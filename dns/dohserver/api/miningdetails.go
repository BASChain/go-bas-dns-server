package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/BASChain/go-bas-dns-server/dns/mem"
)

type MiningDetail struct {

}


type MiningDetailReq struct {
	Wallet string `json:"wallet"`
	PageNumber int `json:"pagenumber"`
	PageSize int `json:"pagesize"`
}



type MiningDetailItem struct {
	ReceiptHash string `json:"receipthash"`
	RootDomainName string `json:"rootdomainname"`
	OpName string `json:"opname"`
	FromDomainName string `json:"fromdomainname"`
	ToMiner string `json:"tominer"`
	ToOwner string `json:"toowner"`
}

type MiningDetailResp struct {
	State      int                `json:"state"`
	TotalPage  int                `json:"totalpage"`
	PageNumber int                `json:"pagenumber"`
	PageSize   int                `json:"pagesize"`
	Mdis       []*MiningDetailItem `json:"mdis"`
}

func (md *MiningDetail)  ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	req := &MiningDetailReq{}

	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	addr := common.HexToAddress(req.Wallet)

	store:=mem.GetMinerProfitStore()
	store.Lock()
	m:=store.GetProfitMiner(&addr)
	store.UnLock()

	resp:=&MiningDetailResp{}
	resp.PageSize = req.PageSize
	resp.PageNumber = req.PageNumber


	m.Lock()
	defer m.UnLock()

	//lMiner:=m.GetProfitItem4MinerList()
	lOwner:=m.GetProfitItemList()

	resp.TotalPage = int(lOwner.Count())

	start:=(req.PageNumber-1)*req.PageSize
	end:=req.PageNumber*req.PageSize

	cursor:=lOwner.ListIteratorB(start,end)
	if cursor.Count() <= 0{
		resp.State = 0
	}else{
		for{
			n:=cursor.Next()
			if n == nil{
				break
			}
			ni:=n.(*mem.ProfitItem)
			mdi:=&MiningDetailItem{}
			mdi.ReceiptHash = ni.GetReceiptHash().String()
			mdi.FromDomainName = ""
			mdi.RootDomainName = ""
			mdi.OpName = ni.GetSrcTyp()

			resp.Mdis = append(resp.Mdis,mdi)
		}

	}



}













