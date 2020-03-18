package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/BASChain/go-bas-dns-server/config"
	"github.com/BASChain/go-bas/Transactions"
	"math/big"
	"github.com/BASChain/go-bas-dns-server/dns/mem"
)

type FreeEth struct {

}


type FreeEthReq struct {
	Wallet string	`json:"wallet"`
	Amount string	`json:"amount"`
}

type FreeEthResp struct {
	Wallet string `json:"wallet"`
	State int     `json:"state"`
	ErrMsg    string  `json:"errmsg"`
	Amount string  `json:"amount"`
}

func NewFreeEth() *FreeEth {
	return &FreeEth{}
}


//var key *keystore.Key = nil
//
//func RestoreKey() error {
//	if key == nil{
//		keys:=Account.PrivateKeyRecover(config.GetKeyFile(),"secret")
//		if len(keys) == 0{
//			return errors.New("load key error")
//		}
//		key = keys[0]
//	}
//
//	return nil
//}


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
	if amount == ""{
		amount = config.GetBasDCfg().FreeEthAmount
	}


	feresp := &FreeEthResp{}
	feresp.Wallet = fer.Wallet

	var flag bool

	z:=big.Int{}
	sndamount,b:=z.SetString(amount,10)

	if !b{
		feresp.ErrMsg = "Amount Error"
		feresp.State = 0
		flag = true
	}
	var state int
	state,err=mem.GetState(addr,mem.ETH)
	if err == nil{
		if state == mem.SUCCESS{
			feresp.State = 0
			feresp.ErrMsg = "You have Applied"
			flag = true
		}

		if state == mem.WAITING{
			feresp.State = 0
			feresp.ErrMsg = "Your Applying is running"
			flag = true
		}
	}

	if !flag{
		feresp.Amount = amount
		feresp.State = 1
		feresp.ErrMsg = "success"
		Transactions.SendFreeEthWrapper(config.GetLoanKey(),addr,sndamount)
	}

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