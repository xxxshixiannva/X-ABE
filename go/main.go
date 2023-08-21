package main

import (
	"fmt"
)

func main() {
	// init ip list
	fmt.Println("[init]======================== Init IP of all nodes list ======================")
	status, ipinfo := NewIPinfo()
	if status == false {
		return
	}
	fmt.Println("[init]Success!")

	// init PedersenVSS
	fmt.Println("[init]===================== Init Parameters of PedersenVSS ====================")
	status, pedersenvss := NewPedersenVSS(ipinfo)
	fmt.Println("[init]Success!")

	// init FeldmanVSS
	fmt.Println("[init]===================== Init Parameters of FeldmanVSS =====================")
	status, feldmanvss := NewFeldmanVSS()
	if status == false {
		return
	}
	fmt.Println("[init]Success!")

	// init sender
	fmt.Println("[init]========================= Init Message Sender ===========================")
	status, sender := NewMessageSender(ipinfo, pedersenvss, feldmanvss)
	fmt.Println("[init]Success!")
	// check all nodes online
	go func() {
		if sender.onlinenum == len(sender.IPinfo.IPlist)-1 {
			for _, value := range sender.Conns {
				value.Write([]byte("all nodes online|||||"))
			}
			fmt.Println("[log]All Nodes Online!")
			sender.step1_generate()
		}
		return
	}() 

	// init listener
	fmt.Println("[init]======================== Init Message Receiver ==========================")
	listener := NewMessageListener(ipinfo, sender)
	fmt.Println("[init]Success!")
	listener.Start()
}