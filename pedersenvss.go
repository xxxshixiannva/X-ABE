package main

import (
	"crypto/rand"
	"math/big"
)

// PedersenVSS
type PedersenVSS struct {
	p1          *big.Int
	p2          *big.Int
	p3          *big.Int
	t           int
	N           *big.Int
	Polynomial1 []*big.Int
	Polynomial2 []*big.Int
	g           *big.Int
	h           *big.Int
	ipinfo      *IPinfo
}

func NewPedersenVSS(ipinfo *IPinfo) (bool, *PedersenVSS) {
	a, _ := new(big.Int).SetString("1363895147340162124487750544377566700025348452567", 10)
	b, _ := new(big.Int).SetString("1257354545315887944833595666025792933231792977521", 10)
	c, _ := new(big.Int).SetString("1296657106138026641358592699056954007605324218609", 10)

	pedersenvss := &PedersenVSS{
		p1:          a,
		p2:          b,
		p3:          c,
		t:           len(ipinfo.IPlist) / 2,
		Polynomial1: make([]*big.Int, 0),
		Polynomial2: make([]*big.Int, 0),
		g:           big.NewInt(3),
		h:           big.NewInt(65537),
		ipinfo:      ipinfo,
	}

	pedersenvss.N = new(big.Int)
	pedersenvss.N.Mul(pedersenvss.p1, pedersenvss.p2)
	pedersenvss.N.Mul(pedersenvss.N, pedersenvss.p3)

	for i := 0; i < pedersenvss.t+1; i++ {
		randomBigInt, _ := rand.Int(rand.Reader, pedersenvss.N)
		pedersenvss.Polynomial1 = append(pedersenvss.Polynomial1, randomBigInt)
	}
	for i := 0; i < pedersenvss.t+1; i++ {
		randomBigInt, _ := rand.Int(rand.Reader, pedersenvss.N)
		pedersenvss.Polynomial2 = append(pedersenvss.Polynomial2, randomBigInt)
	}

	return true, pedersenvss
}