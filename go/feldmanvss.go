package main

import (
	"fmt"
	"github.com/Nik-U/pbc"
	"math/big"
)

// FeldmanVSS
type FeldmanVSS struct {
	pairing *pbc.Pairing
	g       *pbc.Element
	p       *big.Int
	n0      *big.Int
	n1      *big.Int
	n2      *big.Int
	l       int
}

func NewFeldmanVSS() (bool, *FeldmanVSS) {
	feldmanvss := &FeldmanVSS{}

	// params
	pStr := "3602291881362578269408900972923883981249023743695260275790375337088899103553606296737776176077021631283499575659377869566382616215275016597346338059"
	p, _ := new(big.Int).SetString(pStr, 10)
	n0, _ := new(big.Int).SetString("1363895147340162124487750544377566700025348452567", 10)
	n1, _ := new(big.Int).SetString("1257354545315887944833595666025792933231792977521", 10)
	n2, _ := new(big.Int).SetString("1296657106138026641358592699056954007605324218609", 10)
	n := new(big.Int)
	n.Mul(n0, n1)
	n.Mul(n, n2)
	l := 1620

	feldmanvss.p = p
	feldmanvss.n0 = n0
	feldmanvss.n1 = n1
	feldmanvss.n2 = n2
	feldmanvss.l = l

	params := fmt.Sprintf("type a1\np %s\nn %s\nn0 %s\nn1 %s\nn2 %s\nl %d", pStr, n.String(), n0.String(), n1.String(), n2.String(), l)

	// elliptic curve
	paramsObj, err := pbc.NewParamsFromString(params)
	if err != nil {
		fmt.Println("[Error]Something wrong with method pbc.NewParamsFromString()")
		return false, feldmanvss
	}
	feldmanvss.pairing = pbc.NewPairing(paramsObj)

	// g1
	xstr := "952915521556523589053324002114806389324617573949095646588529396283787676834376782374438266727784956841412620903849225915513670943755291617331017855"
	gX, _ := new(big.Int).SetString(xstr, 10)
	ystr := "1501612210867762095431526541616590680978425303013763238315111339334102400746738556118006988352474767350695416424976512889273006328883996609061373589"
	gY, _ := new(big.Int).SetString(ystr, 10)
	xBytes := gX.Bytes()
	yBytes := gY.Bytes()
	serializedData := append(xBytes, yBytes...)
	feldmanvss.g = feldmanvss.pairing.NewG1().SetBytes(serializedData)
	return true, feldmanvss
}