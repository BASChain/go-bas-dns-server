package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/BASChain/go-bas/Market"
	"github.com/BASChain/go-bas/Bas_Ethereum"
	"encoding/hex"
	"github.com/kprc/nbsnetwork/common/list"
)

type SellingDomain struct {

}

type SellingDomainReq struct {
	Wallet string `json:"wallet"`
	PageNumber int	`json:"pagenumber"`
	PageSize int	`json:"pagesize"`

}

type SellingDomainResp struct {
	State int `json:"state"`
	TotalPage int `json:"totalpage"`
	PageNumber int			`json:"pagenumber"`
	PageSize int			`json:"pagesize"`
	Domains []*ExpensiveDomain `json:"domains"`
}

func NewSellingDomain() *SellingDomain {
	return &SellingDomain{}
}

func sellingDomainSort(v1,v2 interface{}) int {
	e1,e2:=v1.(*ExpensiveDomain),v2.(*ExpensiveDomain)

	if e1.RegTime  < e2.RegTime {
		return 1
	}

	return -1
}




func (sd *SellingDomain)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
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

	req:=&SellingDomainReq{}
	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	if req.PageNumber < 1{
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	var addr *common.Address
	if req.Wallet != "" {
		_addr := common.HexToAddress(req.Wallet)
		addr = &_addr
	}

	var owners []map[Bas_Ethereum.Hash]*Market.SellOrder



	if addr != nil{
		if m,ok:=Market.SellOrders[*addr];ok{
			owners = append(owners,m)
		}
	}else{
		for _,v:=range Market.SellOrders{
			owners = append(owners,v)
		}
	}


	sortList := list.NewList(expensiveCmp)
	sortList.SetSortFunc(sellingDomainSort)

	for i:=0;i<len(owners);i++{
		for k,v:=range owners[i]{
			d:=GetRecord(k)
			if d == nil {
				continue
			}
			ed:=&ExpensiveDomain{}
			ed.Domain = string(d.Name)
			ed.PriceOmit = v.GetPrice()
			ed.Price = v.GetPriceStr()
			t,_ := Bas_Ethereum.GetTimestamp(v.BlockNumber)
			ed.RegTime = int64(t)
			ed.ExpireTime = d.GetExpire()
			ed.Owner = d.GetOwner()
			ed.Hash = "0x" + hex.EncodeToString(k[:])
			ed.OrderTime = v.GetTime()

			sortList.AddValueOrder(ed)

		}
	}

	cnt:=0
	b:=(req.PageNumber-1)*req.PageSize
	e:=req.PageNumber * req.PageSize

	cursor:=sortList.ListIterator(e)

	resp := &SellingDomainResp{}



	if cursor.Count() <= b {
		resp.State = 0
	}else{
		resp.State = 1
		for{
			d:=cursor.Next()
			if d == nil{
				break
			}
			if cnt >=b && cnt <e{
				resp.Domains = append(resp.Domains,d.(*ExpensiveDomain))
			}
			cnt ++
		}
	}
	
	resp.PageSize = req.PageSize
	resp.PageNumber = req.PageNumber
	resp.TotalPage = cnt

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
