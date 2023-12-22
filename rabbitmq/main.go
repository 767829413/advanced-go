package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

var (
	user = "user"
	pwd  = "user"
	addr = "mw.internal.cn"
	port = "5672"

	url = fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		user,
		pwd,
		addr,
		port,
	)

	config = amqp.Config{
		Vhost:      "/", //设置服务的Vhost
		Properties: amqp.Table{},
	}
	config1 = amqp.Config{
		Vhost:      "vhost1",
		Properties: amqp.Table{},
	}
	config2 = amqp.Config{
		Vhost:      "vhost2",
		Properties: amqp.Table{},
	}
	config3 = amqp.Config{
		Vhost:      "vhost3",
		Properties: amqp.Table{},
	}
	config4 = amqp.Config{
		Vhost:      "vhost4",
		Properties: amqp.Table{},
	}
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	c, cancel := context.WithCancel(context.Background())
	go exec(c, url, config)
	go exec(c, url, config1)
	go exec(c, url, config2)
	go exec(c, url, config3)
	go exec(c, url, config4)

	time.Sleep(time.Second * 40)
	cancel()
}

func exec(c context.Context, url string, config amqp.Config) {
	conn, err := amqp.DialConfig(url, config)
	failOnError(err, config.Vhost+": Failed to connect to RabbitMQ")
	defer conn.Close()
	// 创建一个通道
	ch, err := conn.Channel()
	failOnError(err, config.Vhost+": Failed to open a channel")
	defer ch.Close()

	// 声明一个队列
	q, err := ch.QueueDeclare(
		config.Vhost+": hello", // 队列名称
		false,                  // 是否持久化
		false,                  // 是否自动删除
		false,                  // 是否独占模式
		false,                  // 是否阻塞
		nil,                    // 额外属性
	)
	failOnError(err, config.Vhost+": Failed to declare a queue")
	var ok = make(chan struct{})
	// 发送消息到队列
	go func() {
		var seq = 1
		for {
			body := fmt.Sprintf(config.Vhost+": Hello, RabbitMQ---%d!", seq)
			err = ch.Publish(
				"",     // 交换机名称
				q.Name, // 队列名称
				false,  // 是否强制
				false,  // 是否立即发送
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			failOnError(err, config.Vhost+": Failed to publish a message")
			fmt.Println(config.Vhost + ": Message sent!")
			ok <- struct{}{}
			seq++
			time.Sleep(2 * time.Second)
		}
	}()

	// 接收来自队列的消息
	msgs, err := ch.Consume(
		q.Name, // 队列名称
		"",     // 消费者名称
		true,   // 是否自动应答
		false,  // 是否独占模式
		false,  // 是否阻塞
		false,  // 是否等待
		nil,    // 额外属性
	)
	failOnError(err, config.Vhost+": Failed to register a consumer")

	// 处理接收到的消息
	go func() {
		for msg := range msgs {
			fmt.Println(config.Vhost+": Received message:", string(msg.Body))
			<-ok
		}
	}()

	// 等待程序退出
	<-c.Done()
}
