package api

import (
	"github.com/BASChain/go-bas/DataSync"
	"github.com/BASChain/go-bas/Bas_Ethereum"
)

func QueryBasByDomainName(q string) *DataSync.DomainRecord {
	hash := Bas_Ethereum.GetHash(q)

	if n,ok := DataSync.Records[hash];!ok{
		return nil
	}else{
		return &n
	}
}

func QueryBasByBCAddr(addr []byte) *DataSync.DomainRecord {
	return nil
}