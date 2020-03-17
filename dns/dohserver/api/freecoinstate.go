package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/BASChain/go-bas-dns-server/dns/mem"
)

type FreeCoinState struct {

}

type FreeCoinStateReq struct {
	Wallet string `json:"wallet"`
	Type   int `json:"type"`
}

type FreeCoinStateResp struct {
	Wallet string `json:"wallet"`
	Type int `json:"type"`
	State int `json:"state"`
	Msg string `json:"msg"`
}

func NewFreeCoinState() *FreeCoinState {
	return &FreeCoinState{}
}

func (fcs *FreeCoinState) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
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

	req:= &FreeCoinStateReq{}

	err = json.Unmarshal(body,req)
	if err!=nil{
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	addr:=common.HexToAddress(req.Wallet)
	
	resp:=&FreeCoinStateResp{}
	resp.Wallet = req.Wallet
	resp.Type = req.Type

	if req.Type != mem.ETH && req.Type != mem.BAS{
		resp.Msg = "Type Error"
	}else{
		var state int
		state,err = mem.GetState(addr,req.Type)
		if err == nil{
			resp.State = state
			if resp.State == mem.WAITING{
				resp.Msg = "Waiting a result"
			}
			if resp.State == mem.SUCCESS{
				resp.Msg = "Success"
			}
			if resp.State == mem.FAILURE{
				resp.Msg = "Failure"
			}
		}else{
			resp.Msg = "No State"
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