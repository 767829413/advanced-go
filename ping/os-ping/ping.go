package osping

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os"
	"os/exec"
)

func OsPing() {
	host := os.Args[1]
	output, err := exec.Command("ping", host).CombinedOutput()
	if err != nil {
		panic(err.Error())
	}
	// 处理命令行中文转码的问题
	newByte, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(output)
	fmt.Println(string(newByte))
}
