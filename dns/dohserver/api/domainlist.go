package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/BASChain/go-bas/DataSync"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"net/http"
	"github.com/BASChain/go-bas/Bas_Ethereum"
	"github.com/BASChain/go-bas/Market"
)

type DomainList struct {
}

type DomainListReq struct {
	Wallet     string `json:"wallet"`
	PageNumber int    `json:"pageNumber"`
	PageSize   int    `json:"pageSize"`
}

type DomainListItem struct {
	IsOrder     bool   `json:"isorder"`
	Name        string `json:"name"`
	Expire      int64  `json:"expire"`
	OpenApplied bool   `json:"openApplied"`
	Hash        string `json:"hash"`
}

type DomainListResp struct {
	State int               `json:"state"`
	Owner string            `json:"owner"`
	Data  []*DomainListItem `json:"data"`
}

func NewDomainList() *DomainList {
	return &DomainList{}
}

func IsOrder(addr common.Address,hash Bas_Ethereum.Hash) bool  {
	if m,ok:=Market.SellOrders[addr];!ok{
		return false
	}else{
		if _,ok:=m[hash];ok{
			return true
		}
	}
	return false
}

func (dl *DomainList) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

	dtl := &DomainListReq{}

	err = json.Unmarshal(body, dtl)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	addr := common.HexToAddress(dtl.Wallet)

	dtlresp := &DomainListResp{Owner: dtl.Wallet}

	DataSync.MemLock()
	defer DataSync.MemUnlock()

	hasharr, ok := DataSync.Assets[addr]
	if !ok {
		dtlresp.State = 0
	} else {
		dtlresp.State = 1
	}

	for i := (dtl.PageNumber - 1) * dtl.PageSize; i < len(hasharr) && i < (dtl.PageNumber)*dtl.PageSize; i++ {
		dtli := &DomainListItem{}
		dm, ok := DataSync.Records[hasharr[i]]
		if !ok {
			continue
		}
		dtli.IsOrder = IsOrder(*dm.GetOwnerOrig(),hasharr[i])
		dtli.Name = dm.GetName()
		dtli.Expire = dm.GetExpire()
		dtli.OpenApplied = dm.GetOpenStatus()
		dtli.Hash = "0x" + hex.EncodeToString(hasharr[i][:])
		dtlresp.Data = append(dtlresp.Data, dtli)
	}

	var bresp []byte

	bresp, err = json.Marshal(*dtlresp)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	w.WriteHeader(200)
	w.Write(bresp)

	return

}
