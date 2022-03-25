package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	list := os.Args
	path := list[1]
	fileInfo := delErr(os.Stat(path))
	fmt.Println("file name:", fileInfo.Name())
	fmt.Println("file size:", fileInfo.Size())

	conn := delErr(net.Dial("tcp", "localhost:8000"))
	defer conn.Close()

	// 发送文件名
	_ = delErr(conn.Write([]byte(fileInfo.Name())))

	// 读取返回消息
	buf := make([]byte, 1024)
	n := delErr(conn.Read(buf))
	if "ok" == string(buf[:n]) {
		fmt.Println("receive server: ok")
		sendFile(conn, path)
	} else {
		fmt.Println("receive server: no")
	}

}

func delErr[T any](res T, err error) T {
	if err != nil {
		panic(err)
	}
	return res

}

func sendFile(conn net.Conn, filePath string) {
	f := delErr(os.Open(filePath))
	defer f.Close()
	buf := make([]byte, 1024*1024)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("send file over")
			} else {
				panic(err)
			}
			return
		}
		_ = delErr(conn.Write(buf[:n]))
	}
}
