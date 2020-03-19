package api

import (
	"encoding/json"
	"fmt"
	"github.com/BASChain/go-bas/DataSync"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"net/http"
)

type DomainTotal struct {
}

type DomainTotalReq struct {
	Wallet string `json:"wallet"`
}

type DomainTotalResp struct {
	State int `json:"state"`
	Data  int `json:"data"`
}

const (
	PathNotFound    int = 1
	AddressNotFound int = 2
)

type PathError struct {
	ErrorCode int
	ErrorMsg  string
}

func EncapError(code int, msg string) []byte {
	pe := PathError{ErrorCode: code, ErrorMsg: msg}

	jpe, err := json.Marshal(pe)
	if err != nil {
		return []byte("{}")
	}

	return jpe
}

func NewDomainTotal() *DomainTotal {
	return &DomainTotal{}
}

func (dt *DomainTotal) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	dtr := &DomainTotalReq{}

	err = json.Unmarshal(body, dtr)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	addr := common.HexToAddress(dtr.Wallet)

	dtresp := &DomainTotalResp{}

	DataSync.MemLock()
	defer DataSync.MemUnlock()

	hasharr, ok := DataSync.Assets[addr]
	if !ok {
		dtresp.State = 0
	} else {
		dtresp.State = 1
		dtresp.Data = len(hasharr)
	}

	var bresp []byte

	bresp, err = json.Marshal(*dtresp)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	w.WriteHeader(200)
	w.Write(bresp)

	return
}
