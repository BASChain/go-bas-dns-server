package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/BASChain/go-bas/Market"
	"github.com/BASChain/go-bas/DataSync"
	"github.com/BASChain/go-bas/Bas_Ethereum"
	"github.com/kprc/nbsnetwork/common/list"
	"math/big"
)

type ExpensiveDomains struct {

}


type ExpensiveDomainsReq struct {
	PageNumber int			`json:"pagenumber"`
	PageSize int			`json:"pagesize"`
}

type ExpensiveDomain struct {
	Domain string `json:"domain"`
	Price string `json:"price"`
	PriceOmit *big.Int `json:"-"`
	RegTime int64 `json:"regtime"`
	ExpireTime int64 `json:"expiretime"`
	Owner string `json:"owner"`
}

type ExpensiveDomainsResp struct {
	State int `json:"state"`
	TotalPage int `json:"totalpage"`
	PageNumber int	`json:"pagenumber"`
	PageSize int `json:"pagesize"`
	Domains []*ExpensiveDomain `json:"domains"`
}

func NewExpensiveDomains() *ExpensiveDomains {
	return &ExpensiveDomains{}
}

func expensiveCmp(v1 interface{},v2 interface{}) int  {
	e1,e2:=v1.(ExpensiveDomain),v2.(ExpensiveDomain)

	if e1.Domain == e2.Domain{
		return 0
	}
	return 1
}

func expensiveSort(v1 interface{},v2 interface{}) int  {
	e1,e2:=v1.(ExpensiveDomain),v2.(ExpensiveDomain)

	if e1.PriceOmit.Cmp(e2.PriceOmit) < 0{
		return 1
	}

	return -1

}


func GetRecord(hash Bas_Ethereum.Hash) *DataSync.DomainRecord  {
	if d,ok:=DataSync.Records[hash];!ok{
		return nil
	}else{
		return d
	}
}

func (ed *ExpensiveDomains) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	req:=&ExpensiveDomainsReq{}
	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	expensiveList := list.NewList(expensiveCmp)
	expensiveList.SetSortFunc(expensiveSort)

	for _,v:=range Market.SellOrders{
		for kk,vv:=range v{
			d:=GetRecord(kk)
			if d == nil {
				continue
			}
			e:=&ExpensiveDomain{}
			e.Domain = string(d.Name)
			e.PriceOmit = vv.GetPrice()
			e.Price = vv.GetPriceStr()
			t,_ := Bas_Ethereum.GetTimestamp(vv.BlockNumber)
			e.RegTime = int64(t)
			e.ExpireTime = d.GetExpire()
			e.Owner = d.GetOwner()

			expensiveList.AddValueOrder(e)
		}
	}


	cnt:=0
	b:=(req.PageNumber-1)*req.PageSize
	e:=req.PageNumber*req.PageSize

	resp:=&ExpensiveDomainsResp{}

	var ds []*ExpensiveDomain

	cursor:=expensiveList.ListIterator(e)
	if cursor.Count() <= b{
		resp.State = 0
	}else{
		resp.State = 1
		for{
			nxt:=cursor.Next()
			if nxt == nil{
				break
			}
			if cnt >=b && cnt <e{
				ds = append(ds,nxt.(*ExpensiveDomain))
			}
			cnt ++
		}
		resp.Domains = ds
	}
	resp.PageSize = req.PageSize
	resp.PageNumber = req.PageNumber
	resp.TotalPage = cursor.Count()

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










