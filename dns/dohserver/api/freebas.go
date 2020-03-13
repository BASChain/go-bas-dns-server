package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/BASChain/go-bas/Transactions"
)

type FreeBas struct {

}


type FreeBasReq struct {
	Wallet string `json:"wallet"`
	Amount int64  `json:"amount"`
}


type FreeBasResp struct {
	State int `json:"state"`
	Wallet string `json:"wallet"`
	Amount int64  `json:"amount"`
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

	var b bool
	b,err = Transactions.CheckIfApplied(addr)
	if b{
		resp.State = 0
	}else{
		resp.State = 1

		RestoreKey()

		go Transactions.SendFreeBas(key,addr)

		resp.Wallet = fbr.Wallet
		amount,_ := Transactions.GetFreeBasAmount()
		if amount == nil{
			resp.Amount = 0
		}else{
			resp.Amount = amount.Int64()
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