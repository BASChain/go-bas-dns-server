package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/BASChain/go-bas-dns-server/config"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/BASChain/go-bas/Account"
	"github.com/pkg/errors"
	"github.com/BASChain/go-bas/Transactions"
)

type FreeEth struct {

}


type FreeEthReq struct {
	Wallet string	`json:"wallet"`
	Amount int64	`json:"amount"`
}

type FreeEthResp struct {
	Wallet string `json:"wallet"`
	State int     `json:"state"`
	Amount int64  `json:"amount"`
}

func NewFreeEth() *FreeEth {
	return &FreeEth{}
}


var key *keystore.Key = nil

func RestoreKey() error {
	if key == nil{
		keys:=Account.PrivateKeyRecover(config.GetKeyFile(),"secret")
		if len(keys) == 0{
			return errors.New("load key error")
		}
		key = keys[0]
	}

	return nil
}


func (fe *FreeEth)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST" {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	var body []byte
	var err error

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	fer:=&FreeEthReq{}

	err = json.Unmarshal(body,fer)
	if err!=nil{
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	addr:=common.HexToAddress(fer.Wallet)
	amount := fer.Amount
	if amount == 0{
		amount = config.GetBasDCfg().FreeTokenAmount
	}

	err = RestoreKey()
	if err!=nil{
		panic("load key failed")
	}

	go Transactions.SendFreeEth(key,addr,amount)

	feresp := &FreeEthResp{}
	feresp.Amount = amount
	feresp.Wallet = fer.Wallet
	feresp.State = 1

	var bresp []byte

	bresp,err =json.Marshal(*feresp)
	if err != nil{
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	w.WriteHeader(200)
	w.Write(bresp)

}