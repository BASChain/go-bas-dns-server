package exlib

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"github.com/BASChain/go-bas-dns-server/dns/mem"
	"github.com/BASChain/go-bas/Transactions"
)

func SendFreeBasByContractWrapper(key *keystore.Key,addr common.Address,amount *big.Int){
	mem.Update(addr,mem.BAS,mem.WAITING)
	go func() {
		if err:=Transactions.SendFreeBasByContract(key,addr,amount);err!=nil{
			mem.Update(addr,mem.BAS,mem.FAILURE)
		}else {
			mem.Update(addr,mem.BAS,mem.SUCCESS)
		}
	}()
}



func SendFreeEthWrapper(key *keystore.Key,toAddress common.Address,amount *big.Int)  {
	mem.Update(toAddress,mem.ETH,mem.WAITING)
	go func() {
		if err:=Transactions.SendFreeEth(key,toAddress,amount);err!=nil{
			mem.Update(toAddress,mem.ETH,mem.FAILURE)
		}else{
			mem.Update(toAddress,mem.ETH,mem.SUCCESS)
		}
	}()
}