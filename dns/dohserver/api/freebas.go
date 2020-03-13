package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/BASChain/go-bas/Transactions"
	"github.com/BASChain/go-bas-dns-server/config"
	"math/big"
)

type FreeBas struct {

}


type FreeBasReq struct {
	Wallet string `json:"wallet"`
	Amount string  `json:"amount"`
}


type FreeBasResp struct {
	State int `json:"state"`
	Wallet string `json:"wallet"`
	Amount string `json:"amount"`
	ErrMsg string `json:"errmsg"`
}

func NewFreeBas() *FreeBas {
	return &FreeBas{}
}

func (fb *FreeBas)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
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

	fbr:=&FreeBasReq{}

	err = json.Unmarshal(body,fbr)
	if err!=nil{
		w.WriteHeader(500)
		fmt.Fprintf(w,"{}")
		return
	}
	addr:=common.HexToAddress(fbr.Wallet)

	resp:=&FreeBasResp{}
	resp.Wallet = fbr.Wallet

	var b bool
	b,err = Transactions.CheckIfApplied(addr)
	if b{
		resp.State = 0
		resp.ErrMsg = "You have Applied"
	}else{

		amount := fbr.Amount
		if amount == ""{
			amount = config.GetBasDCfg().FreeBasAmount
		}

		z:=&big.Int{}
		sndamount,b:=z.SetString(amount,10)
		if !b{
			resp.State = 0
			resp.ErrMsg = "Amount error"
		}else{
			resp.State = 1
			RestoreKey()
			go Transactions.SendFreeBasByContract(key,addr,sndamount)
			resp.Amount = amount
		}

	}


	var bresp []byte

	bresp,err =json.Marshal(*resp)
	if err != nil{
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	w.WriteHeader(200)
	w.Write(bresp)

}