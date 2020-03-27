package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/BASChain/go-bas/DataSync"
	"github.com/BASChain/go-bas/Bas_Ethereum"
	"github.com/kprc/nbsnetwork/common/list"
	"math/big"
	"github.com/BASChain/go-bas/Market"
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
	From string  `json:"from,omitempty"`
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
	e1,e2:=v1.(*ExpensiveDomain),v2.(*ExpensiveDomain)

	if e1.Domain == e2.Domain{
		return 0
	}
	return 1
}

func expensiveSort(v1 interface{},v2 interface{}) int  {
	e1,e2:=v1.(*ExpensiveDomain),v2.(*ExpensiveDomain)

	if e1.PriceOmit.Cmp(e2.PriceOmit) < 0{
		return 1
	}

	return -1

}

func fDo(arg interface{}, v interface{}) (ret interface{},err error) {
	e1,e2:=arg.(*ExpensiveDomain),v.(*ExpensiveDomain)

	if e1.PriceOmit.Cmp(e2.PriceOmit) > 0{
		e2.PriceOmit = e1.PriceOmit
		e2.Price = e1.Price
		e2.RegTime = e1.RegTime
		e2.ExpireTime = e1.ExpireTime
		e2.Owner = e1.Owner
		e2.From = e1.From
	}

	return e1,nil
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

	for i:=0;i<len(Market.Sold);i++{
		deal:=Market.Sold[i]
		d:=GetRecord(deal.GetHash())
		if d == nil || deal.GetAGreedPrice() == nil{
			continue
		}
		ed:=&ExpensiveDomain{}
		ed.From = deal.GetFromOwner()
		ed.Owner = deal.GetOwner()
		ed.PriceOmit = deal.GetAGreedPrice()
		ed.Price = ed.PriceOmit.String()
		ed.ExpireTime = d.GetExpire()
		ed.RegTime = d.GetRegTime()
		ed.Domain = d.GetName()

		if _,err:=expensiveList.FindDo(ed,fDo);err!=nil{
			expensiveList.AddValue(ed)
		}
	}

	expensiveList.Sort()

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










