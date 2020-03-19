package api

import (
	"github.com/BASChain/go-bas/Bas_Ethereum"
	"github.com/BASChain/go-bas/DataSync"
)

func QueryBasByDomainName(q string) *DataSync.DomainRecord {
	hash := Bas_Ethereum.GetHash(q)

	if n, ok := DataSync.Records[hash]; !ok {
		return nil
	} else {
		return &n
	}
}

func QueryBasByBCAddr(addr []byte) *DataSync.DomainRecord {
	return nil
}
