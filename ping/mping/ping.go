//go:build linux
// +build linux

package mping

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"go.uber.org/ratelimit"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/sys/unix"
)

var (
	// ErrStampNotFund is returned when timestamp not found
	ErrStampNotFund       = errors.New("timestamp not found")
	id                    uint16
	validTargets          = make(map[string]bool)
	supportTxTimestamping = true
	supportRxTimestamping = true

	// 发送报文前缀,便于区分
	msgPrefix   = []byte("fuvking")
	targetAddrs []string
	stat        *buckets
)

func init() {
	id = uint16(os.Getpid() & 0xffff)
}

var connOnce sync.Once

func start() error {
	for _, target := range targetAddrs {
		validTargets[target] = true
	}

	if len(targetAddrs) == 0 {
		return errors.New("no target")
	}

	conn, err := openConn()
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	if *tos > 0 {
		pconn := ipv4.NewConn(conn)
		err = pconn.SetTOS(*tos)
		if err != nil {
			return fmt.Errorf("failed to set tos: %w", err)
		}
	}

	done := make(chan error, 3)
	go func() {
		err := send(conn)
		done <- err
	}()

	go func() {
		err := printStat()
		done <- err
	}()
	go func() {
		read(conn)
		done <- err
	}()

	return <-done
}

func openConn() (*net.IPConn, error) {
	// 使用 net.ListenPacket 函数创建一个 IP 连接，指定使用 IPv4 和 ICMP 协议，监听所有可用的网络接口
	conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, err
	}

	ipconn := conn.(*net.IPConn)
	// 使用 File 方法获取 IP 连接的文件描述符
	f, err := ipconn.File()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 将文件描述符转换为整数，以便后续使用系统调用
	fd := int(f.Fd())

	// https://patchwork.ozlabs.org/project/netdev/patch/1226415412.31699.2.camel@ecld0pohly/
	// https://www.kernel.org/doc/Documentation/networking/timestamping.txt
	// 定义时间戳标志，表示要启用的时间戳选项
	flags := unix.SOF_TIMESTAMPING_SYS_HARDWARE | unix.SOF_TIMESTAMPING_RAW_HARDWARE | unix.SOF_TIMESTAMPING_SOFTWARE | unix.SOF_TIMESTAMPING_RX_HARDWARE | unix.SOF_TIMESTAMPING_RX_SOFTWARE |
		unix.SOF_TIMESTAMPING_TX_HARDWARE | unix.SOF_TIMESTAMPING_TX_SOFTWARE |
		unix.SOF_TIMESTAMPING_OPT_CMSG | unix.SOF_TIMESTAMPING_OPT_TSONLY

	// 使用 syscall.SetsockoptInt 设置 socket 选项 SO_TIMESTAMPING，启用所定义的时间戳标志
	if err := syscall.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_TIMESTAMPING, flags); err != nil {
		// 处理时间戳选项设置失败的情况
		// 如果设置 SO_TIMESTAMPING 失败，将标志 supportTxTimestamping 和 supportRxTimestamping 设置为 false.
		supportTxTimestamping = false
		supportRxTimestamping = false

		// 尝试单独设置 SO_TIMESTAMPNS，并根据是否成功设置 supportRxTimestamping
		if err := syscall.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_TIMESTAMPNS, 1); err == nil {
			supportRxTimestamping = true
		}

		return ipconn, nil
	}
	// 设置接收和发送超时时间,接收和发送的超时时间均为 1 秒
	timeout := syscall.Timeval{Sec: 1, Usec: 0}
	if err := syscall.SetsockoptTimeval(fd, unix.SOL_SOCKET, unix.SO_RCVTIMEO, &timeout); err != nil {
		return nil, err
	}
	if err := syscall.SetsockoptTimeval(fd, unix.SOL_SOCKET, unix.SO_SNDTIMEO, &timeout); err != nil {
		return nil, err
	}

	return ipconn, nil
}

// 用于构建发送报文
var payload []byte

func send(conn *net.IPConn) error {
	defer connOnce.Do(func() { conn.Close() })

	// 获取文件描述符
	// 通过 conn.File() 获取连接的文件描述符
	f, err := conn.File()
	if err != nil {
		return err
	}
	defer f.Close()
	// 将文件描述符转换为整数,以便进行系统调用
	fd := int(f.Fd())

	// 限流器，按照需要的速率发送,使用 ratelimit.New 创建一个速率限制器，以控制发送速率
	limiter := ratelimit.New(*rate, ratelimit.Per(time.Second))

	// 初始化 ICMP 报文序列号
	var seq uint16

	// 构建 ICMP 报文数据
	data := make([]byte, *packetSize)
	copy(data, msgPrefix)

	// 生成随机数据填充报文
	_, err = rand.Read(data[len(msgPrefix)+8:])
	if err != nil {
		return err
	}

	// 保存报文中的有效负载
	payload = data[len(msgPrefix)+8:]

	// 构建目标 IP 地址列表
	targets := make([]*net.IPAddr, 0, len(targetAddrs))
	for _, taget := range targetAddrs {
		targets = append(targets, &net.IPAddr{IP: net.ParseIP(taget)})
	}

	// 循环发送 ICMP 报文
	sentPackets := 0
	for {
		// 用来将发送的包和回来的包匹配
		seq++

		// 使用速率限制器控制发送速率
		limiter.Take()
		for _, target := range targets {
			// 把发送时的时间戳放入payload, 以便计算时延
			ts := time.Now().UnixNano()
			binary.LittleEndian.PutUint64(data[len(msgPrefix):], uint64(ts))

			// 构建发送的icmp包, ICMP Echo 请求报文
			req := &icmp.Message{
				Type: ipv4.ICMPTypeEcho,
				Body: &icmp.Echo{
					ID:   int(id),
					Seq:  int(seq),
					Data: data,
				},
			}

			key := ts / int64(time.Second)

			// 将 ICMP 报文序列化为二进制数据
			data, err := req.Marshal(nil)
			if err != nil {
				continue
			}

			// 发送 ICMP 报文到目标地址
			_, err = conn.WriteTo(data, target)
			if err != nil {
				return err
			}

			// 构建并保存结果信息
			rs := &Result{
				txts:   ts,
				target: target.IP.String(),
				seq:    seq,
			}

			// 如果支持 TX 时间戳，则获取 TX 时间戳并保存
			if supportTxTimestamping {
				if txts, err := getTxTs(fd); err != nil {
					if strings.HasPrefix(err.Error(), "resource temporarily unavailable") {
						continue
					}
					fmt.Printf("failed to get TX timestamp: %s", err)
					rs.txts = txts
				}
			}

			// 将结果信息添加到统计中
			stat.Add(key, rs)
		}

		// 更新发送的计数器
		sentPackets++

		// 如果指定了发送的总数，达到总数后休眠并返回
		if *count > 0 && sentPackets >= *count {
			time.Sleep(time.Second * time.Duration((*delay + 1)))
			fmt.Printf("sent packets: %d\n", sentPackets)
			return nil
		}
	}
}

func read(conn *net.IPConn) error {
	defer connOnce.Do(func() { conn.Close() })

	// 初始化读取缓冲区和带外数据缓冲区
	pktBuf := make([]byte, 1500)
	oob := make([]byte, 1500)

	for {
		// 设置读取超时为 10 毫秒
		_ = conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
		// 读取ICMP返回的包
		n, oobn, _, ra, err := conn.ReadMsgIP(pktBuf, oob)

		if err != nil {
			return err
		}

		var rxts int64
		// 如果支持 RX 时间戳，从带外数据中获取 RX 时间戳
		if supportRxTimestamping {
			if rxts, err = getTsFromOOB(oob, oobn); err != nil {
				return fmt.Errorf("failed to get RX timestamp: %s", err)
			}
		} else {
			// 不支持 RX 时间戳就使用当前时间作为 RX 时间戳
			rxts = time.Now().UnixNano()
		}

		// 如果读取的字节数小于 IPv4 报文头长度，报文格式错误
		if n < ipv4.HeaderLen {
			return errors.New("malformed IPv4 packet")
		}

		// 获取响应地址
		target := ra.String()

		// 过滤不是发包过程中设置的响应地址
		if !validTargets[target] {
			continue
		}

		// 解析 ICMP 消息
		msg, err := icmp.ParseMessage(1, pktBuf[ipv4.HeaderLen:n])
		if err != nil {
			continue
		}

		// 如果消息类型不是 ICMP Echo Reply，则继续下一次循环
		if msg.Type != ipv4.ICMPTypeEchoReply {
			continue
		}

		switch pkt := msg.Body.(type) {
		case *icmp.Echo:
			// 检查id,是不是当前进程id,不匹配则继续下一次循环
			if uint16(pkt.ID) != id {
				continue
			}

			// 如果 Echo 请求的数据长度小于要求的最小长度，则继续下一次循环
			if len(pkt.Data) < len(msgPrefix)+8 {
				continue
			}
			// 如果 Echo 请求的数据不以预定的前缀开头，则继续下一次循环
			if !bytes.HasPrefix(pkt.Data, msgPrefix) {
				continue
			}

			// 从 Echo 请求的数据中提取 TX 时间戳
			txts := int64(binary.LittleEndian.Uint64(pkt.Data[len(msgPrefix):]))
			key := txts / int64(time.Second)

			// 是否检查数据是否被修改（bitflip）
			bitflip := false
			if *bitflipCheck {
				bitflip = !bytes.Equal(pkt.Data[len(msgPrefix)+8:], payload)
			}

			// 将结果添加到统计中
			stat.AddReply(key, &Result{
				txts:     txts,
				rxts:     rxts,
				target:   target,
				latency:  time.Now().UnixNano() - txts,
				received: true,
				seq:      uint16(pkt.Seq),
				bitflip:  bitflip,
			})
		}
	}
}

func printStat() error {
	// 获取延迟时间（以秒为单位）,这里默认 3s
	delayInSeconds := int64(*delay)
	// 创建定时器，每秒触发一次
	ticker := time.NewTicker(time.Second)

	// 记录上一个统计数据的时间戳
	var lastKey int64

	for range ticker.C {
	recheck:
		// 获取最新的统计数据
		bucket := stat.Last()
		if bucket == nil {
			continue
		}

		// 如果统计数据的时间戳小于等于上一个时间戳，弹出并继续检查
		if bucket.Key <= lastKey {
			stat.Pop()
			goto recheck
		}

		// 如果统计数据的时间戳小于当前时间减去延迟时间，弹出并继续检查
		if bucket.Key <= time.Now().UnixNano()/int64(time.Second)-delayInSeconds {
			pop := stat.Pop().(*Bucket)
			if pop.Key < bucket.Key {
				goto recheck
			}

			// 更新上一个时间戳
			lastKey = pop.Key

			// 初始化存储每个目标的结果信息的关联
			targetResult := make(map[string]*TargetResult)

			// 加上一个读锁避免并发冲突
			pop.Mu.RLock()

			// 遍历每个结果，将信息汇总到 targetResult 映射中
			for _, r := range pop.Value {
				target := r.target

				tr := targetResult[target]
				if tr == nil {
					tr = &TargetResult{}
					targetResult[target] = tr
				}

				tr.latency += r.latency

				if r.received {
					tr.received++
				} else {
					tr.loss++
				}

				if *bitflipCheck && r.bitflip {
					tr.bitflipCount++
				}

			}
			pop.Mu.RUnlock()

			// 打印每个目标的统计信息
			for target, tr := range targetResult {
				total := tr.received + tr.loss
				var lossRate float64
				if total == 0 {
					lossRate = 0
				} else {
					lossRate = float64(tr.loss) / float64(total)
				}

				logLevel := "INFO"
				if tr.loss > 0 {
					logLevel = "WARN"
				}

				// 根据是否进行 bitflip 检查，打印不同的信息
				if *bitflipCheck {
					if tr.received == 0 {
						log.Printf(
							"[%s] %s: sent:%d, recv:%d, loss rate: %.2f%%, latency: %v, bitflip: %d\n",
							logLevel,
							target,
							total,
							tr.received,
							lossRate*100,
							0,
							tr.bitflipCount,
						)
					} else {
						log.Printf("[%s] %s: sent:%d, recv:%d,  loss rate: %.2f%%, latency: %v, bitflip: %d\n", logLevel, target, total, tr.received, lossRate*100, time.Duration(tr.latency/int64(tr.received)), tr.bitflipCount)
					}
				} else {
					if tr.received == 0 {
						log.Printf("[%s] %s: sent:%d, recv:%d, loss rate: %.2f%%, latency: %v\n", logLevel, target, total, tr.received, lossRate*100, 0)
					} else {
						log.Printf("[%s] %s: sent:%d, recv:%d,  loss rate: %.2f%%, latency: %v\n", logLevel, target, total, tr.received, lossRate*100, time.Duration(tr.latency/int64(tr.received)))
					}
				}
			}

		}
	}
	return nil
}

func getTsFromOOB(oob []byte, oobn int) (int64, error) {
	cms, err := syscall.ParseSocketControlMessage(oob[:oobn])
	if err != nil {
		return 0, err
	}
	for _, cm := range cms {
		if cm.Header.Level == syscall.SOL_SOCKET && cm.Header.Type == syscall.SO_TIMESTAMPING {
			var t unix.ScmTimestamping
			if err := binary.Read(bytes.NewBuffer(cm.Data), binary.LittleEndian, &t); err != nil {
				return 0, err
			}

			for i := 0; i < len(t.Ts); i++ {
				if t.Ts[i].Nano() > 0 {
					return t.Ts[i].Nano(), nil
				}
			}

			return 0, ErrStampNotFund
		}

		if cm.Header.Level == syscall.SOL_SOCKET && cm.Header.Type == syscall.SCM_TIMESTAMPNS {
			var t unix.Timespec
			if err := binary.Read(bytes.NewBuffer(cm.Data), binary.LittleEndian, &t); err != nil {
				return 0, err
			}
			return t.Nano(), nil
		}
	}
	return 0, ErrStampNotFund
}

func getTxTs(socketFd int) (int64, error) {
	pktBuf := make([]byte, 1024)
	oob := make([]byte, 1024)
	_, oobn, _, _, err := syscall.Recvmsg(
		socketFd,
		pktBuf,
		oob,
		syscall.MSG_ERRQUEUE|syscall.MSG_DONTWAIT,
	)
	if err != nil {
		return 0, err
	}
	return getTsFromOOB(oob, oobn)
}
