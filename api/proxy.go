package api

import (
	"fmt"
	"net"

	"gitlab.paradise-soft.com.tw/dwh/legion/service"
)

func ListenTCP(port int) (err error) {
	var l net.Listener
	l, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err == nil {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					// log.Printf("ListenTCP crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
				}
			}()
			for {
				var conn net.Conn
				conn, err = l.Accept()
				if err == nil {
					go func() {
						defer func() {
							if e := recover(); e != nil {
								// log.Printf("connection handler crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
							}
						}()
						service.TCPCallback(conn)
					}()
				} else {
					fmt.Printf("accept error , ERR:%s", err)
					break
				}
			}
		}()
	}
	return
}
