package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/BASChain/go-bas/Market"
	"github.com/BASChain/go-bas/Bas_Ethereum"
	"strings"
	"encoding/hex"
)

type MarketSearch struct {

}

type MarketSearchReq struct {
	Text string `json:"text"`
	PageNumber int	`json:"pagenumber"`
	PageSize int	`json:"pagesize"`
}


type MarketSearchResp struct {
	State int `json:"state"`
	TotalPage int `json:"totalpage"`
	PageNumber int			`json:"pagenumber"`
	PageSize int			`json:"pagesize"`

	Domains []*ExpensiveDomain `json:"domains"`
}

func NewMarketSearch() *MarketSearch {
	return &MarketSearch{}
}

func (ms *MarketSearch)ServeHTTP(w http.ResponseWriter, r *http.Request){
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

	req:=&MarketSearchReq{}
	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	cnt:=0
	b:=(req.PageNumber-1)*req.PageSize
	e:=req.PageNumber * req.PageSize

	resp:=&MarketSearchResp{}

	searchText := req.Text
	for _,v:=range Market.SellOrders{
		for kk,vv:=range v{
			d:=GetRecord(kk)
			if d == nil {
				continue
			}

			if !strings.Contains(string(d.Name),searchText){
				continue
			}
			if cnt >=b && cnt <e{
				ed:=&ExpensiveDomain{}
				ed.Domain = string(d.Name)
				ed.PriceOmit = vv.GetPrice()
				ed.Price = vv.GetPriceStr()
				t,_ := Bas_Ethereum.GetTimestamp(vv.BlockNumber)
				ed.RegTime = int64(t)
				ed.ExpireTime = d.GetExpire()
				ed.Owner = d.GetOwner()
				ed.Hash = "0x"+hex.EncodeToString(kk[:])

				resp.Domains = append(resp.Domains,ed)
			}
			cnt ++

		}
	}

	resp.PageSize = req.PageSize
	resp.PageNumber = req.PageNumber
	if len(resp.Domains)  == 0{
		resp.State = 0
	}else {
		resp.State = 1
	}
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