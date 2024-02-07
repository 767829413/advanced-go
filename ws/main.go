package main

import (
	"bufio"
	"fmt"
	"github.com/767829413/advanced-go/util"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	sendChan = make(chan []byte)
	pongChan = make(chan []byte)
	history  = make([]string, 0)
)

func main() {
	c := GetConf()
	run(c)
}

func run(c *Config) {
	wsUrl := getWsUrl(c)
	log.Println("wsUrl is: ", wsUrl)
	conn := connect(wsUrl)
	go send(conn)
	go receiver(conn)
	Login(c.UserName, c.GroupId)
	Ready()
	if c.UserType == "speaker" {
		EnableVideo(c.LoginName)
		TurnOnLive()
	}
	r := regexp.MustCompile("[ 	]+")

	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		command := strings.Trim(string(data), " 	")
		if len(command) == 0 {
			continue
		}
		command = preDeal(command)
		if len(command) == 0 {
			fmt.Println("show history end")
			continue
		}
		cmds := r.Split(command, -1)
		dealInput(cmds)
	}
}

func preDeal(cmd string) string {
	if cmd == "h" {
		for _, hiscmd := range history {
			fmt.Println(hiscmd)
		}
		return ""
	} else if strings.HasPrefix(cmd, "h") {
		numstr := cmd[1:]
		num, _ := strconv.Atoi(numstr)
		if num > 0 || num < 10 {
			if len(history) >= num {
				cmd = history[len(history)-num]
			}
		}
	}
	history = append(history, cmd)
	if len(history) > 10 {
		history = history[len(history)-10:]
	}
	return cmd
}

func receiver(conn *websocket.Conn) {
	conn.SetPingHandler(func(message string) error {
		pongChan <- []byte(message)
		return nil
	})
	conn.SetPongHandler(func(message string) error {
		return nil
	})
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read data failed as: %s\n", err)
			return
		}

		if len(msg) == 0 {
			continue
		}
		log.Printf("[SERVER DATA] >>>> %s\n", msg)
	}
}

func send(conn *websocket.Conn) {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case msg := <-sendChan:
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("write data failed as: %s\n", err)
			}
		case msg := <-pongChan:
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("write pong failed as: %s\n", err)
			}
		case <-ticker.C:
			msg := []byte(fmt.Sprintf("c[1,2,%d]", time.Now().UnixNano()/1e6))
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("write ping failed as: %s\n", err)
			}
		}
	}
}

func connect(wsUrl string) *websocket.Conn {
	dialer := websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment, // From default dialer
		HandshakeTimeout:  45 * time.Second,          // From default dialer
		EnableCompression: true,
	}
	conn, _, err := dialer.Dial(wsUrl, nil)
	if err != nil {
		log.Fatal("dial failed as ", err)
	}

	return conn
}

func getWsUrl(c *Config) string {
	if len(c.Query) > 0 {
		return c.Domain + "/upws?" + c.Query
	} else {
		query := buildQuery(c, c.AppKey)
		return c.Domain + "/upws?" + query.Encode()
	}
}

func buildQuery(c *Config, key string) url.Values {
	param := url.Values{}
	param.Set("appId", c.AppId)
	param.Set("meetingId", c.MeetingId)
	param.Set("loginName", c.LoginName)
	param.Set("userType", c.UserType)
	param.Set("validTime", "3600000")
	param.Set("validBegin", strconv.FormatInt(util.NowMs(), 10))
	param.Set("meetingType", c.MeetingType)
	if c.EndTimeRelativeTime != 0 {
		param.Set("endTime", strconv.FormatInt(util.NowMs()+(c.EndTimeRelativeTime)*1000, 10))
	}
	signature, _ := util.EncryptUrlValue(param, key)
	param.Set("signature", signature)

	return param
}

const serverKey = "d4478d30ec7f4735a0614c419d863938"

type Config struct {
	Domain    string
	AppId     string `json:"appId"`
	AppKey    string `json:"appKey"`
	GroupId   string `json:"groupId"`
	MeetingId string `json:"meetingId"`
	// 1 普通课堂 ,5 说课教研
	MeetingType string `json:"meetingType"`
	LoginName   string `json:"loginName"`
	UserType    string `json:"userType"`
	UserName    string `json:"userName"`
	Query       string `json:"query"`
	// 相对于第一个客户端启动时，课堂的结束时间，单位秒;默认无
	EndTimeRelativeTime int64 `json:"endTimeRelativeTime"`
}

func GetConf() *Config {
	return &Config{
		Domain:              "ws://localhost:80",
		AppId:               "plaso",
		AppKey:              "123456",
		MeetingId:           "hello",
		LoginName:           "user1",
		UserType:            "listener",
		UserName:            "username1",
		GroupId:             "1",
		EndTimeRelativeTime: 0,
	}
}

const (
	UKeyframe = iota
	UEdit
	_
	UImage
	ULine
	URube
	UStraightLine
	URect
	UTriangle
	UCircle
	UEllipse              // 椭圆
	UDeskShareOnOff  = 99 // 桌面共享开关
	UScribbleLayer   = 100
	UMediaLayer      = 101
	URoomLayerStatus = 102
)

const (
	AttrDelete = iota
	AttrSize
	AttrPosition
	AttrUndo
	AttrUndoCount
	AttrAppendPoints
	AttrMakeEnd
	AttrZIndex
	AttrActive
	AttrBgAddr
	AttrChangeLayer
	AttrRotate
)

const (
	MIC int64 = 1 << iota
	CAM
	PEN
	ONLINE    = 1 << 8  // 在线状态
	RaiseHand = 1 << 10 // 举手状态
	OnWall    = 1 << 11 // 上墙状态
	StagePerm = MIC | CAM | PEN
)

const (
	FeatureMessage     = 1 // 消息状态
	FeatureLive        = 2 // 直播状态
	FeaturePrivateChat = 4 // 私聊状态
)

const (
	UNDO = "1"
	REDO = "2"
)

func sendData(data string) {
	sendChan <- []byte(data)
}

var elementId = 0

func drawElement(elementType int, attrs string) int {
	elementId += 1
	a := fmt.Sprintf("[[2],%d,%d,%s]", elementType, elementId, attrs)
	sendData(a)
	return elementId
}

func editElement(elementType int, elementId int, attrId int, attrValue string) {
	var attr string
	if attrValue == "" {
		attr = fmt.Sprintf("[%d]", attrId)
	} else {
		attr = fmt.Sprintf("[%d,%s]", attrId, attrValue)
	}
	a := fmt.Sprintf("[[2],%d,%d,%d,[%s]]", UEdit, elementType, elementId, attr)
	sendData(a)
}

func Login(userName, groupId string) {
	sendData(fmt.Sprintf("c[100,%q,\"a.jpg\",[%q,%q]]", userName, groupId, ""))
}
func Ready() {
	sendData("c[21]")
}
func Logout() {
	a := "c[102, 9]"
	sendData(a)
}
func EndLesson() {
	sendData("c[105]")
}
func DisablePen(userId string) {
	sendData(fmt.Sprintf("c[3,0,%q,%d,0]", userId, PEN))
}
func EnablePen(userId string) {
	sendData(fmt.Sprintf("c[3,1,%q,%d,0]", userId, PEN))
}
func DisableVideo(userId string) {
	sendData(fmt.Sprintf("c[3,0,%q,%d,0]", userId, CAM))
}
func EnableVideo(userId string) {
	sendData(fmt.Sprintf("c[3,1,%q,%d,0]", userId, CAM))
}
func TurnOnLive() {
	sendData(fmt.Sprintf("c[24,1,%d]", FeatureLive))
}
func TurnOffLive() {
	sendData(fmt.Sprintf("c[24,0,%d]", FeatureLive))
}
func ReportVideo() {
	sendData(fmt.Sprintf("c[9,1,%d]", CAM))
}

func CreateScribbleLayer() {
	var addr = "scribble 1"
	a := fmt.Sprintf("[[0],%d,0,%q,[0,0],[0,0,0,1]]", UScribbleLayer, addr)
	sendData(a)
}

func DrawImage() int {
	var addr = "abc"
	attrs := fmt.Sprintf("%q,[0,0],[300,400],0", addr)
	return drawElement(UImage, attrs)
}
func DrawLine() int {
	attrs := fmt.Sprintf("[0,0],1,[0,0,0,1]")
	return drawElement(ULine, attrs)
}
func DrawRube() int {
	attrs := fmt.Sprintf("[0,0],1,false")
	return drawElement(URube, attrs)
}
func DrawSolidLine() int {
	attrs := fmt.Sprintf("[0,0,1,1],1,[0,0,0]")
	return drawElement(UStraightLine, attrs)
}
func DrawRect() int {
	attrs := fmt.Sprintf("[0,0],[300,400],1,[0,0,0]")
	return drawElement(URect, attrs)
}
func DrawTriangle() int {
	attrs := fmt.Sprintf("[0,0,100,0,0,100],1,[0,0,0]")
	return drawElement(UTriangle, attrs)
}
func DrawCircle() int {
	attrs := fmt.Sprintf("[0,0],100,1,[0,0,0]")
	return drawElement(UCircle, attrs)
}
func DrawEllipse() int {
	attrs := fmt.Sprintf("[0,0],[300,400],1,[0,0,0]")
	return drawElement(UEllipse, attrs)
}

func AppendPoint(lineId int) {
	val := "[100,100]"
	editElement(ULine, lineId, AttrAppendPoints, val)
}
func MakeEnd(lineId int) {
	editElement(ULine, lineId, AttrMakeEnd, "")
}

func ChangeLayer(layerId int) {
	a := fmt.Sprintf("[[0],%d,%d,0,[[%d,%d]]]", UEdit, URoomLayerStatus, AttrChangeLayer, layerId)
	sendData(a)
}

func DeleteLine(lineId int) {
	editElement(ULine, lineId, AttrDelete, "")
}

func DeleteImage(id int) {
	editElement(UImage, id, AttrDelete, "")
}

func DeleteRect(id int) {
	editElement(URect, id, AttrDelete, "")
}

func Undo() {
	editElement(100, 2, AttrUndo, UNDO)
}
func Redo() {
	editElement(100, 2, AttrUndo, REDO)
}

func ChangeAddr(id int) {
	editElement(UImage, id, AttrBgAddr, "\"www.baidu.com\"")
}

func ShareDesk() {
	a := fmt.Sprintf("[[1],%d,10]", UDeskShareOnOff)
	sendData(a)
}

func dealInput(cmds []string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic while deal input as: %s\n", err)
		}
	}()
	if len(cmds) == 0 {
		return
	}

	switch cmds[0] {
	case "lo":
		Logout()
	case "ready":
		Ready()
	case "end":
		EndLesson()
	case "layer":
		CreateScribbleLayer()
	case "image":
		DrawImage()
	case "line":
		DrawLine()
	case "rube":
		DrawRube()
	case "solidLine":
		DrawSolidLine()
	case "rect":
		DrawRect()
	case "triangle":
		DrawTriangle()
	case "circle":
		DrawCircle()
	case "ellipse":
		DrawEllipse()
	case "point", "points":
		lineId, _ := strconv.Atoi(cmds[1])
		AppendPoint(lineId)
	case "endLine":
		lineId, _ := strconv.Atoi(cmds[1])
		MakeEnd(lineId)
	case "disablePen":
		DisablePen(cmds[1])
	case "enablePen":
		EnablePen(cmds[1])
	case "disableVideo":
		DisableVideo(cmds[1])
	case "enableVideo":
		EnableVideo(cmds[1])
	case "changeLayer":
		layerId, _ := strconv.Atoi(cmds[1])
		ChangeLayer(layerId)
	case "desk":
		ShareDesk()
	case "deleteLine":
		id, _ := strconv.Atoi(cmds[1])
		DeleteLine(id)
	case "deleteImage":
		id, _ := strconv.Atoi(cmds[1])
		DeleteImage(id)
	case "deleteRect":
		id, _ := strconv.Atoi(cmds[1])
		DeleteRect(id)
	case "undo":
		Undo()
	case "redo":
		Redo()
	case "changeAddr", "addr":
		id, _ := strconv.Atoi(cmds[1])
		ChangeAddr(id)
	case "turnOnLive":
		TurnOnLive()
	case "turnOffLive":
		TurnOffLive()
	case "reportVideo":
		ReportVideo()
	default:
		sendData(cmds[0])
	}
}
