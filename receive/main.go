package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	var host string
	fmt.Print("请输入需要发送的ip地址，按回车键结束 \n")
	fmt.Scanf("%s", &host)
	fmt.Println("localhost:", host)
	listener := delErr(net.Listen("tcp", "localhost:"+host))
	defer listener.Close()

	conn := delErr(listener.Accept())
	defer conn.Close()

	buf := make([]byte, 4096)
	n := delErr(conn.Read(buf))
	fmt.Println("receive file name:", string(buf[:n]))
	conn.Write([]byte("ok"))
	receiveFile(conn, string(buf[:n]))
}

func delErr[T any](res T, err error) T {
	if err != nil {
		panic(err)
	}
	return res
}

func receiveFile(conn net.Conn, fileName string) {
	f := delErr(os.Create(fileName))
	defer f.Close()
	buf := make([]byte, 1024*1024)

	for {
		n, _ := conn.Read(buf)
		if n == 0 {
			fmt.Println("receive file over")
			return
		}
		_ = delErr(f.Write(buf[:n]))
	}
}
