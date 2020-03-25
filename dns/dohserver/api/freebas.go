package api

import (
	"encoding/json"
	"fmt"
	"github.com/BASChain/go-bas-dns-server/config"
	"github.com/BASChain/go-bas-dns-server/dns/mem"
	"github.com/BASChain/go-bas/Transactions"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"math/big"
	"net/http"
	"github.com/BASChain/go-bas-dns-server/dns/dohserver/exlib"
)

type FreeBas struct {
}

type FreeBasReq struct {
	Wallet string `json:"wallet"`
	Amount string `json:"amount"`
}

type FreeBasResp struct {
	State  int    `json:"state"`
	Wallet string `json:"wallet"`
	Amount string `json:"amount"`
	ErrMsg string `json:"errmsg"`
}

func NewFreeBas() *FreeBas {
	return &FreeBas{}
}

func (fb *FreeBas) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	fbr := &FreeBasReq{}

	err = json.Unmarshal(body, fbr)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	addr := common.HexToAddress(fbr.Wallet)

	resp := &FreeBasResp{}
	resp.Wallet = fbr.Wallet

	var flag bool

	var b bool
	b, err = Transactions.CheckIfApplied(addr)
	if b {
		resp.State = 0
		resp.ErrMsg = "You have Applied"
		flag = true
	}

	if !flag {
		var stat int
		stat, err = mem.GetState(addr, mem.BAS)
		if err == nil {
			if stat == mem.SUCCESS {
				resp.State = 0
				resp.ErrMsg = "You have Applied"
				flag = true
			}
			if stat == mem.WAITING {
				resp.State = 0
				resp.ErrMsg = "Your Applying is running"
				flag = true
			}
		}
	}

	var sndamount *big.Int
	var amount string

	if !flag {
		amount = fbr.Amount
		if amount == "" {
			amount = config.GetBasDCfg().FreeBasAmount
		}

		z := &big.Int{}
		sndamount, b = z.SetString(amount, 10)
		if !b {
			resp.State = 0
			resp.ErrMsg = "Amount error"
			flag = true
		}
	}

	if !flag {
		resp.State = 1

		exlib.SendFreeBasByContractWrapper(config.GetLoanKey(), addr, sndamount)
		resp.Amount = amount
		resp.ErrMsg = "success"

	}

	var bresp []byte

	bresp, err = json.Marshal(*resp)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	w.WriteHeader(200)
	w.Write(bresp)

}
