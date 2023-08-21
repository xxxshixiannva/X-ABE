package main

import (
	"sync"
	"net"
	"strings"
	"fmt"
	"io"
)

type MessageListener struct {
	IPinfo    *IPinfo
	mapLock   sync.RWMutex
	sender    *MessageSender
}

func NewMessageListener(ipinfo *IPinfo, ms *MessageSender) *MessageListener {
	messagelistener := &MessageListener{
		IPinfo:    ipinfo,
		sender:    ms,
	}
	return messagelistener
}

func (this *MessageListener) Start() {
	// openlistener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.IPinfo.Myip, this.IPinfo.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	fmt.Println("[log]local listener started")
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}
		go this.Handler(conn)
	}
}

func (this *MessageListener) Handler(conn net.Conn) {
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err", err)
				return
			}
			msgstring := string(buf[:n-5])
			// process data from other node
			msgs := strings.Split(msgstring, "|||||")
			for _, msg := range msgs {
				if msg == "all nodes online" {                                    // init 
					fmt.Println("[log]All Nodes Online!")
					for _, ip := range this.sender.IPinfo.IPlist {
						_, ok := this.sender.Conns[ip]
						if this.sender.IPinfo.Myip != ip && !ok {
							conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, this.sender.IPinfo.Port))
							if err != nil {
								continue
							}
							this.sender.Conns[ip] = conn
						}
					}
					this.sender.step1_generate()
				} else if strings.Contains(msg, "[Cik Broadcast from") {          // step1 generate Cik
					fmt.Println(msg)
					ind := strings.Index(msg, "]")
					ip := msg[20:ind]
					values := msg[ind+1:]
					this.sender.maplock.Lock()
					this.sender.cik_map[ip] = values
					this.sender.maplock.Unlock()
					_, ok := this.sender.sij_map[ip]
					if ok && len(this.sender.cik_map) == len(this.IPinfo.IPlist)-1 && len(this.sender.sij_map) == len(this.IPinfo.IPlist)-1{
						this.sender.step1_verify()
					}
				} else if strings.Contains(msg, "[sij,sij' Private from") {       // step1 generate sij
					fmt.Println(msg)
					ind := strings.Index(msg, "]")
					ip := msg[23:ind]
					values := msg[ind+1:]
					this.sender.maplock.Lock()
					this.sender.sij_map[ip] = values
					this.sender.maplock.Unlock()
					_, ok := this.sender.cik_map[ip]
					if ok && len(this.sender.cik_map) == len(this.IPinfo.IPlist)-1 && len(this.sender.sij_map) == len(this.IPinfo.IPlist)-1{
						this.sender.step1_verify()
					}
				} else if strings.Contains(msg, "[Aik Broadcast from") {       // step1 generate sij
					ind := strings.Index(msg, "]")
					ip := msg[20:ind]
					values := msg[ind+1:]
					this.sender.maplock.Lock()
					this.sender.aik_map[ip] = values
					this.sender.maplock.Unlock()
					if len(this.sender.aik_map) == len(this.IPinfo.IPlist)-1 {
						this.sender.step4_verify()
					}
				} else {
					fmt.Println(msg)
				}
			}
			
		}
	}()
}
