package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	listener := delErr(net.Listen("tcp", "localhost:8000"))
	fmt.Println("loacl host", listener.Addr())
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
