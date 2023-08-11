package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
)

func sendFile(fileName string, w http.ResponseWriter) {
	path := "/root/go/src/github.com/767829413/advanced-go/demo/net/http/server/static/" + fileName
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("os.Open error: ", err)
		w.Write([]byte("No such file or directory"))
		return
	}
	defer f.Close()

	buf := make([]byte, 4096)
	for {
		n, _ := f.Read(buf)
		if n == 0 {
			return
		}
		w.Write(buf[:n])
	}

}

func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request url: ", r.URL.String())
	fileName := path.Base(r.URL.String())
	fmt.Println("request fileName: ", fileName)
	sendFile(fileName, w)
}

func main() {
	http.HandleFunc("/", myHandler)

	http.ListenAndServe("wsl:9898", nil)
}
