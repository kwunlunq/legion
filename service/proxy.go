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

	tracer.Infof("testrp", "use proxy : %s", address)
	//os.Exit(0)
	err = OutToTCP(address, &inConn, &req)
	if err != nil {
		tracer.Errorf("testrp", "connect to %s fail, ERR:%s", address, err)

		glob.CloseConn(&inConn)
	}
}
func OutToTCP(address string, inConn *net.Conn, req *glob.HTTPRequest) (err error) {
	// inAddr := (*inConn).RemoteAddr().String()
	inLocalAddr := (*inConn).LocalAddr().String()
	//防止死循环
	if IsDeadLoop(inLocalAddr, req.Host) {
		glob.CloseConn(inConn)
		err = fmt.Errorf("dead loop detected , %s", req.Host)
		return
	}
	proxys, err := glob.GetProxies(3, nil, proxy.SetPassSites("leisu"))
	if err != nil {
		tracer.Errorf("testrp", "get proxy , err:%s", err)
		glob.CloseConn(inConn)
		return
	}
	var proxyList []string
	for _, proxy := range proxys {
		u, err := url.Parse(proxy)
		if err != nil {
			return err
		}
		proxyList = append(proxyList, u.Host)
	}

	proxyList = append(proxyList, "46.101.79.148:24045")
	// proxyList = []string{
	// 	"46.101.78.176:24045",
	// 	// "168.149.142.170:8080",
	// 	// "198.199.119.119:3128",
	// }
	// , "180.168.13.26:8000"
	// 67.205.149.230:8080
	var outConns []net.Conn
	var outConnsReader []io.ReadWriter
	tracer.Infof("testrp", "conn %s", proxyList)

	for _, proxy := range proxyList {
		var outConn net.Conn

		outConn, err = net.DialTimeout("tcp", proxy, time.Duration(5)*time.Second)

		if err != nil {
			tracer.Errorf("testrp", "connect to %s , err:%s", proxy, err)
		}
		outConn.Write(req.HeadBuf)
		outConns = append(outConns, outConn)
		outConnsReader = append(outConnsReader, outConn)
	}

	if len(outConns) == 0 {
		tracer.Errorf("testrp", "no proxy can used")
		err = errors.New("no proxy can used")
		return
	}

	glob.IoBind(*inConn, outConnsReader, func(isSrcErr bool, err error) {
		if err != nil {
			// log.Println(err)
		}
		// tracer.Infof("testrp", "conn %s - %s - %s -%s released [%s]", inAddr, inLocalAddr, outLocalAddr, outAddr,
		// 	req.Host)
		tracer.Infof("testrp", "close")
		glob.CloseConn(inConn)
		for _, outConn := range outConns {
			glob.CloseConn(&outConn)
		}
	}, func(n int, d bool) {}, 0)

	// for _, outConn := range outConns {
	// 	// var outConn net.Conn
	// 	// // var _outConn interface{}
	// 	// outConn, err = net.DialTimeout("tcp", proxy, time.Duration(5)*time.Second)
	// 	//
	// 	// // _outConn, err = s.outPool.Pool.Get()
	// 	// // if err == nil {
	// 	// // 	outConn = _outConn.(net.Conn)
	// 	// // }
	// 	// if err != nil {
	// 	// 	tracer.Errorf("testrp", "connect to %s , err:%s", proxy, err)
	// 	// 	glob.CloseConn(inConn)
	// 	// 	return
	// 	// }
	//
	// 	// outAddr := outConn.RemoteAddr().String()
	// 	// outLocalAddr := outConn.LocalAddr().String()
	//
	// 	// log.Printf("conn %s - %s - %s - %s connected [%s]", inAddr, inLocalAddr, outLocalAddr, outAddr, req.Host)
	// }

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
