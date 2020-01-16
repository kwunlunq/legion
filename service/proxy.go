package service

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"runtime/debug"
	"time"

	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/dwh/proxy/proxy"
	"gitlab.paradise-soft.com.tw/glob/tracer"
)

func TCPCallback(inConn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			tracer.Errorf("panic", "http(s) conn handler crashed with err : %s \nstack: %s", err, string(debug.Stack()))
		}
	}()
	req, err := glob.NewHTTPRequest(&inConn, 4096)
	if err != nil {
		if err != io.EOF {
			tracer.Errorf("testrp", "decoder error , form %s, ERR:%s", err, inConn.RemoteAddr())
		}
		glob.CloseConn(&inConn)
		return
	}
	address := req.Host
	err = OutToTCP(address, &inConn, &req)
	if err != nil {
		tracer.Errorf("testrp", "connect to %s fail, ERR:%s", address, err)
		glob.CloseConn(&inConn)
	}
}
func OutToTCP(address string, inConn *net.Conn, req *glob.HTTPRequest) (err error) {
	// inAddr := (*inConn).RemoteAddr().String()
	inLocalAddr := (*inConn).LocalAddr().String()
	// 防止死循环
	if IsDeadLoop(inLocalAddr, req.Host) {
		err = fmt.Errorf("dead loop detected , %s", req.Host)
		return
	}
	var proxies, proxyList []string
	proxies, err = glob.GetProxies(3, nil, proxy.SetPassSites("leisu"))
	if err != nil {
		tracer.Errorf("testrp", "get p , err:%s", err)
		return
	}
	for _, p := range proxies {
		u, err := url.Parse(p)
		if err != nil {
			return err
		}
		proxyList = append(proxyList, u.Host)
	}

	// proxyList = append(proxyList, "46.101.79.148:24045")
	// proxyList = []string{
	// 	"46.101.78.176:24045",
	// }
	var outConns []net.Conn
	var proxyConnsReader []io.ReadWriter
	tracer.Tracef("testrp", "conn %s", proxyList)

	for _, p := range proxyList {
		outConn, connErr := net.DialTimeout("tcp", p, time.Duration(5)*time.Second)

		if connErr != nil {
			tracer.Errorf("testrp", "connect to %s , err:%s", p, connErr)
		} else {
			outConn.Write(req.HeadBuf)
			outConns = append(outConns, outConn)
			proxyConnsReader = append(proxyConnsReader, outConn)
		}
	}

	if len(outConns) == 0 {
		err = errors.New("no proxy can used")
		tracer.Errorf("testrp", err.Error())
		return
	}

	glob.IoBind(*inConn, proxyConnsReader, func(isSrcErr bool, err error) {
		glob.CloseConn(inConn)
		for _, outConn := range outConns {
			glob.CloseConn(&outConn)
		}
		if err != nil && err != io.EOF {
			tracer.Errorf("testrp", "conn error: %s", err)
			return
		}
	}, func(n int, d bool) {})

	return
}

func IsDeadLoop(inLocalAddr string, host string) bool {
	inIP, inPort, err := net.SplitHostPort(inLocalAddr)
	if err != nil {
		return false
	}
	outDomain, outPort, err := net.SplitHostPort(host)
	if err != nil {
		return false
	}
	if inPort == outPort {
		var outIPs []net.IP
		outIPs, err = net.LookupIP(outDomain)
		if err == nil {
			for _, ip := range outIPs {
				if ip.String() == inIP {
					return true
				}
			}
		}
		interfaceIPs, err := glob.GetAllInterfaceAddr()
		if err == nil {
			for _, localIP := range interfaceIPs {
				for _, outIP := range outIPs {
					if localIP.Equal(outIP) {
						return true
					}
				}
			}
		}
	}
	return false
}
