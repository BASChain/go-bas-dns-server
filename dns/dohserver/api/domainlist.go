package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/BASChain/go-bas/Bas_Ethereum"
	"github.com/BASChain/go-bas/DataSync"
	"github.com/BASChain/go-bas/Market"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kprc/nbsnetwork/common/list"
	"io/ioutil"
	"net/http"
	"strings"
	"math/big"
)

type DomainList struct {
}

type DomainListReq struct {
	Wallet     string `json:"wallet"`
	PageNumber int    `json:"pageNumber"`
	PageSize   int    `json:"pageSize"`
	DomainType int    `json:"domaintype"`
}

type DomainListItem struct {
	IsOrder     		bool   `json:"isorder"`
	Name        		string `json:"name"`
	Expire      		int64  `json:"expire"`
	OpenApplied 		bool   `json:"openApplied"`
	Hash        		string `json:"hash"`
	RegSubDomainPrice  	string `json:"regsubdomainprice"`
}

type DomainListResp struct {
	State      int               `json:"state"`
	PageNumber int               `json:"pageNumber"`
	PageSize   int               `json:"pageSize"`
	TotalCnt   int               `json:"totalcnt"`
	Owner      string            `json:"owner"`
	Data       []*DomainListItem `json:"data"`
}

func NewDomainList() *DomainList {
	return &DomainList{}
}

func IsOrder(addr common.Address, hash Bas_Ethereum.Hash) bool {
	if m, ok := Market.SellOrders[addr]; !ok {
		return false
	} else {
		if _, ok := m[hash]; ok {
			return true
		}
	}
	return false
}

func domainListSort(v1, v2 interface{}) int {
	d1, d2 := v1.(*DomainListItem), v2.(*DomainListItem)

	if strings.Compare(d1.Name, d2.Name) >= 0 {
		return 1
	}

	return -1

}

func domainListCmp(v1, v2 interface{}) int {
	d1, d2 := v1.(*DomainListItem), v2.(*DomainListItem)

	if d1.Name == d2.Name {
		return 0
	}

	return 1

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

	DataSync.MemLock()
	defer DataSync.MemUnlock()

	dtlresp := &DomainListResp{Owner: dtl.Wallet}
	hasharr, ok := DataSync.Assets[addr]
	if !ok {
		dtlresp.State = 0
	} else {
		dtlresp.State = 1

		dlist := list.NewList(domainListCmp)
		dlist.SetSortFunc(domainListSort)

		for i := 0; i < len(hasharr); i++ {
			dm, ok := DataSync.Records[hasharr[i]]
			if !ok {
				continue
			}
			if dm.IsRoot && dtl.DomainType == 2 {
				continue
			}
			if !dm.IsRoot && dtl.DomainType == 1 {
				continue
			}

			dtli := &DomainListItem{}

			dtli.IsOrder = IsOrder(*dm.GetOwnerOrig(), hasharr[i])
			dtli.Name = dm.GetName()
			dtli.Expire = dm.GetExpire()
			dtli.OpenApplied = dm.GetOpenStatus()
			dtli.Hash = "0x" + hex.EncodeToString(hasharr[i][:])
			if !dm.IsRoot{
				roothash := dm.GetParentHash()
				if droot,ok1:=DataSync.Records[roothash];ok1{
					if droot.RCustomPrice.Cmp(big.NewInt(0)) != 0{
						dtli.RegSubDomainPrice = droot.RCustomPrice.String()
					}
				}
				dtli.RegSubDomainPrice = DataSync.SUBGAS.String()
			}

			dlist.AddValueOrder(dtli)

		}

		cnt := 0
		b := (dtl.PageNumber - 1) * dtl.PageSize
		e := dtl.PageNumber * dtl.PageSize

		cursor := dlist.ListIterator(0)

		if cursor.Count() > b {
			for {
				d := cursor.Next()
				if d == nil {
					break
				}
				if cnt >= b && cnt < e {
					dtlresp.Data = append(dtlresp.Data, d.(*DomainListItem))
				}
				cnt++
			}
		}
		dtlresp.TotalCnt = cursor.Count()

	}

	dtlresp.PageNumber = dtl.PageNumber
	dtlresp.PageSize = dtl.PageSize

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
