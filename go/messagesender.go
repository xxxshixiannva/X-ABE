package main

import (
	"fmt"
	"math/big"
	"net"
	"sync"
	"time"
	"strconv"
	"strings"
)

func findindex(allips []string, myip string) int {
	for i, ip := range allips {
		if ip == myip {
			return i
		}
	}
	return -1
}

// sender
type MessageSender struct {
	IPinfo     *IPinfo
	Conns          map[string]net.Conn
	onlinenum      int
	pedersenvss    *PedersenVSS
	feldmanvss     *FeldmanVSS
	cik_map        map[string]string
	sij_map        map[string]string
	aik_map        map[string]string
	qualified      map[string]bool
	qualified2    map[string]bool
	maplock        sync.RWMutex
	xj             *big.Int
	xj2            *big.Int
}

func NewMessageSender(ipinfo *IPinfo, pedersenvss *PedersenVSS, feldmanvss *FeldmanVSS) (bool, *MessageSender) {
	client := &MessageSender{
		IPinfo:     ipinfo,
		Conns:      make(map[string]net.Conn),
		pedersenvss:    pedersenvss,
	    feldmanvss:     feldmanvss,
		cik_map:    make(map[string]string),
		sij_map:    make(map[string]string),
		aik_map:    make(map[string]string),
		qualified:  make(map[string]bool),
		qualified2: make(map[string]bool),
		xj:         new(big.Int),
		xj2:        new(big.Int),
	}
	for _, ip := range client.IPinfo.IPlist {
		if client.IPinfo.Myip != ip {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, client.IPinfo.Port))
			if err != nil {
				continue
			}
			client.Conns[ip] = conn
		}
	}
	client.onlinenum = len(client.Conns)
	return true, client
}

func (this *MessageSender) step1_generate() {
	time.Sleep(3 * time.Second)
	fmt.Println("[log]================= Step1 Pedersen-VSS Generate and Verify =================")
	time.Sleep(2 * time.Second)

	// broadcast modexp elements
	send := "[Cik Broadcast from " + this.IPinfo.Myip + "]"
	for i:=0; i<this.pedersenvss.t+1; i++ {
		C_ik := big.NewInt(1)
		C_ik.Mul(C_ik, new(big.Int).Exp(this.pedersenvss.g, this.pedersenvss.Polynomial1[i] , this.pedersenvss.N))
		C_ik.Mul(C_ik, new(big.Int).Exp(this.pedersenvss.h, this.pedersenvss.Polynomial2[i] , this.pedersenvss.N))
		C_ik.Mod(C_ik, this.pedersenvss.N)
		send += (C_ik.String()+" ")
	}
	for _, value := range this.Conns {
		_, err := value.Write([]byte(send[:len(send)-1]+"|||||"))
		if err != nil {
			fmt.Println(err)
		}
	}

	// secretly send two polynomials
	time.Sleep(2 * time.Second)
	f1 := make([]*big.Int, 0)
	f2 := make([]*big.Int, 0)
	bigj := new(big.Int)
	for j:=1; j<=len(this.IPinfo.IPlist); j++ {
		bigj, _ = bigj.SetString(strconv.Itoa(j), 10)
		result := big.NewInt(0)
		power, _ := new(big.Int).SetString("1", 10)
		for _, ele := range this.pedersenvss.Polynomial1 {
			tmp := new(big.Int).Mul(ele, power)
			result.Add(result, tmp)
			power.Mul(power, bigj)
		}
		f1 = append(f1, result)
		result = big.NewInt(0)
		power, _ = new(big.Int).SetString("1", 10)
		for _, ele := range this.pedersenvss.Polynomial2 {
			tmp := new(big.Int).Mul(ele, power)
			result.Add(result, tmp)
			power.Mul(power, bigj)
		}
		f2 = append(f2, result)	
	}

	for ip, value := range this.Conns {
		if ip != this.IPinfo.Myip {
			index := findindex(this.IPinfo.IPlist, ip)
			_, err := value.Write([]byte("[sij,sij' Private from " + this.IPinfo.Myip + "]" + f1[index].String() + " " + f2[index].String() + "|||||"))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (this *MessageSender) step1_verify() {
	for ip, _ := range this.cik_map {
		// all Cik
		Ciks_string := strings.Split(this.cik_map[ip], " ")
		//fmt.Println(Ciks_string)
		// sij
		sij1_string := strings.Split(this.sij_map[ip], " ")[0]
		sij2_string := strings.Split(this.sij_map[ip], " ")[1]
		//fmt.Println(sij1_string, sij2_string)
		// calcaulate g^sij1*h^sij2
		left := big.NewInt(1)
		sij1, _ := new(big.Int).SetString(sij1_string, 10)
		sij2, _ := new(big.Int).SetString(sij2_string, 10)
		left.Mul(left, new(big.Int).Exp(this.pedersenvss.g, sij1, this.pedersenvss.N))
		left.Mul(left, new(big.Int).Exp(this.pedersenvss.h, sij2, this.pedersenvss.N))
		left.Mod(left, this.pedersenvss.N)

		// calcaulate g^sij1*h^sij2
		right := big.NewInt(1)
		index := findindex(this.IPinfo.IPlist, this.IPinfo.Myip) + 1
		ind, _ := new(big.Int).SetString(strconv.Itoa(index), 10)
		bigi := new(big.Int)
		for i:=0; i<=this.pedersenvss.t; i++ {
			bigi, _ = bigi.SetString(strconv.Itoa(i), 10)
			cik, _ := new(big.Int).SetString(Ciks_string[i], 10)
			right.Mul(right, new(big.Int).Exp(cik, new(big.Int).Exp(ind, bigi, nil), this.pedersenvss.N))
			right.Mod(right, this.pedersenvss.N)
		}
		if right.String() == left.String() {
			this.maplock.Lock()
			this.qualified[ip] = true
			this.maplock.Unlock()
			fmt.Println("[Verifying Pedersenvss from " + ip + "] Success!")
		} else {
			this.maplock.Lock()
			this.qualified[ip] = false
			this.maplock.Unlock()
			fmt.Println("[Verifying Pedersenvss from " + ip + "] Fail!")
		}
	}
	
	if len(this.qualified) == len(this.IPinfo.IPlist) - 1 {
		fmt.Println("[log]================== Step2 all nodes qualification built ===================")
		this.step3_GenerateSecretKey()
		this.step4_generate()
	}
} 

func (this *MessageSender) step3_GenerateSecretKey() {
	time.Sleep(3 * time.Second)
	fmt.Println("[log]================= Step3 Generate distributed secret key ==================")
	time.Sleep(2 * time.Second)

	this.xj = big.NewInt(0)
	this.xj2 = big.NewInt(0)
	
	// init yourself
	ind := findindex(this.IPinfo.IPlist, this.IPinfo.Myip) + 1
	bigj, _ := new(big.Int).SetString(strconv.Itoa(ind), 10)
	power, _ := new(big.Int).SetString("1", 10)
	for _, ele := range this.pedersenvss.Polynomial1 {
		tmp := new(big.Int).Mul(ele, power)
		this.xj.Add(this.xj, tmp)
		power.Mul(power, bigj)
	}

	power, _ = new(big.Int).SetString("1", 10)
	for _, ele := range this.pedersenvss.Polynomial2 {
		tmp := new(big.Int).Mul(ele, power)
		this.xj2.Add(this.xj2, tmp)
		power.Mul(power, bigj)
	}

	for ip, flag := range this.qualified { 
		if flag == true {
			sij1_string := strings.Split(this.sij_map[ip], " ")[0]
			sij2_string := strings.Split(this.sij_map[ip], " ")[1]
			sij1, _ := new(big.Int).SetString(sij1_string, 10)
			sij2, _ := new(big.Int).SetString(sij2_string, 10)
			this.xj.Add(this.xj, sij1)
			this.xj2.Add(this.xj2, sij2)
			this.xj.Mod(this.xj, this.pedersenvss.N)
			this.xj2.Mod(this.xj2, this.pedersenvss.N)
		}
	}

	fmt.Println("[log]distributed secret key generated")
	fmt.Println(this.xj.String())
}

func (this *MessageSender) step4_generate() {
	time.Sleep(3 * time.Second)
	fmt.Println("[log]==================== Step4 Generate in elliptic curve ====================")
	time.Sleep(2 * time.Second)
	send := []byte("[Aik Broadcast from " + this.IPinfo.Myip + "]")
	seperator := []byte("|||")
	for _, value := range this.pedersenvss.Polynomial1 {
		expElement := this.feldmanvss.pairing.NewZr().SetBig(value.Mod(value, this.pedersenvss.p1))
		result := this.feldmanvss.pairing.NewG1().PowZn(this.feldmanvss.g, expElement)
		rb := result.CompressedBytes()
		send = append(send, rb...)
		send = append(send, seperator...)
	}
	send = append(send, []byte("||")...)
	for _, value := range this.Conns {
		_, err := value.Write(send)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (this *MessageSender) step4_verify() {
	for ip, _ := range this.aik_map {
		Aiks_string := strings.Split(this.aik_map[ip], "|||")
		sij1_string := strings.Split(this.sij_map[ip], " ")[0]
		sij1, _ := new(big.Int).SetString(sij1_string, 10)
		expElement := this.feldmanvss.pairing.NewZr().SetBig(sij1)
		left := this.feldmanvss.pairing.NewG1().PowZn(this.feldmanvss.g, expElement)

		expElement = expElement.SetBig(this.feldmanvss.n0)
		right := this.feldmanvss.pairing.NewG1().Set0()
		index := findindex(this.IPinfo.IPlist, this.IPinfo.Myip) + 1
		ind, _ := new(big.Int).SetString(strconv.Itoa(index), 10)
		bigi := new(big.Int)
		for i:=0; i<=this.pedersenvss.t; i++ {
			compressedBytes := []byte(Aiks_string[i])
			Aik := this.feldmanvss.pairing.NewG1().SetCompressedBytes(compressedBytes)
			bigi, _ = bigi.SetString(strconv.Itoa(i), 10)
			expElement = expElement.SetBig(new(big.Int).Exp(ind, bigi, nil))
			tmp := this.feldmanvss.pairing.NewG1().PowZn(Aik, expElement)
			right = this.feldmanvss.pairing.NewG1().Mul(right, tmp)
		}
		if right.Equals(left) {
			this.maplock.Lock()
			this.qualified2[ip] = true
			this.maplock.Unlock()
			fmt.Println("[Verifying Feldmanvss from " + ip + "] Success!")
		} else {
			this.maplock.Lock()
			this.qualified2[ip] = false
			this.maplock.Unlock()
			fmt.Println("[Verifying Feldmanvss from " + ip + "] Fail!")
		}
	}
	
	if len(this.qualified2) == len(this.IPinfo.IPlist) - 1 {
		this.step5_GeneratePublicKey()
	}
} 

func (this *MessageSender) step5_GeneratePublicKey() {
	time.Sleep(3 * time.Second)
	fmt.Println("[log]================== Step5 Generate distributed public key =================")
	time.Sleep(2 * time.Second)
	
	expElement := this.feldmanvss.pairing.NewZr().SetBig(this.pedersenvss.Polynomial1[0])
	result := this.feldmanvss.pairing.NewG1().PowZn(this.feldmanvss.g, expElement)
	fmt.Println("[Part of Public Key]:",result.String())
	fmt.Print("[CompressedBytes(Part of Public Key)]")
	fmt.Println(result.CompressedBytes())
}
