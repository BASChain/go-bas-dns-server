package api

import (
	"encoding/json"
	"fmt"
	"github.com/BASChain/go-bas-dns-server/dns/mem"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"math/big"
	"net/http"
)

type DrawSummary struct {
}

type DrawSummaryReq struct {
	Wallet string `json:"wallet"`
}

type DrawSummaryResp struct {
	State               int    `json:"state"`
	Wallet              string `json:"wallet"`
	TotalWDrawTimes     int    `json:"totalwdrawtimes"`
	TotalWait2WDraw     string `json:"totalwait2wdraw"`
	TotalWDrawed        string `json:"totalwdrawed"`
	TotalMinerEarned    string `json:"totalminerearned"`
	TotalOwnerEarned    string `json:"totaoownerearned"`
	Wait2WDrawFromMiner string `json:"wait2wdrawfromminer"`
	Wait2WDrawFromOwner string `json:"wait2wdrawfromowner"`
	TotalReceipts       int    `json:"totalreceipts"`
}

func NewDrawSummary() *DrawSummary {
	return &DrawSummary{}
}

func subbigint(x, y big.Int) *big.Int {
	z := &big.Int{}

	z.Sub(&x, &y)

	return z
}

func (dl *DrawSummary) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	req := &DrawSummaryReq{}

	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	addr := common.HexToAddress(req.Wallet)

	store := mem.GetMinerProfitStore()

	resp := &DrawSummaryResp{}
	resp.Wallet = req.Wallet

	store.Lock()

	m := store.GetProfitMiner(&addr)

	if m == nil {
		resp.State = 0
	} else {
		resp.State = 1
		resp.TotalWDrawTimes = m.GetTotalWithdrawTimes()
		resp.TotalReceipts = m.GetTotalReceipts()

		resp.TotalWait2WDraw = subbigint(m.GetTotal4Withdraw(), m.GetTotalWithdraw()).String()
		twd := m.GetTotalWithdraw()
		resp.TotalWDrawed = (&twd).String()
		tme := m.GetTotalWithdrawFromMiner()
		resp.TotalMinerEarned = (&tme).String()
		twfo := m.GetTotalWithdrawFromOwner()
		resp.TotalOwnerEarned = (&twfo).String()
		resp.Wait2WDrawFromMiner = subbigint(m.GetTotalFromMiner(), m.GetTotalWithdrawFromMiner()).String()
		resp.Wait2WDrawFromOwner = subbigint(m.GetTotalFromOwner(), m.GetTotalWithdrawFromOwner()).String()
	}
	store.UnLock()

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
