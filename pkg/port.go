package pkg

import (
	"fmt"
	"net"
)

func ScanPort() {
	ch := make(chan struct{}, 1)
	for i := 1; i <= 65535; i++ {

		wg.Add(1)
		ch <- struct{}{}
		go func() {
			defer wg.Done()
			TcpPort("127.0.0.1")
			<-ch
		}()
	}

	wg.Wait()
}

func TcpPort(address string) {
	conn, err := net.DialTimeout("tcp", address, 10)
	if err != nil {
		fmt.Println(err)
	}
	conn.Close()
	fmt.Println("open:", address)
}
