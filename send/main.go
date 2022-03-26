package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	var (
		path string
		host string
	)
	fmt.Print("请输入需要传入文件地址，按回车键结束 \n")
	fmt.Scanf("%s", &path)
	fmt.Print("请输入需要发送的ip地址，按回车键结束 \n")
	fmt.Scanf("%s", &host)
	fmt.Println("path:", path, "host:", host)

	fileInfo := delErr(os.Stat(path))
	conn := delErr(net.Dial("tcp", host))
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
