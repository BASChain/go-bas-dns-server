package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
	"github.com/BASChain/go-bas/DataSync"
)

type SubDomainList struct {

}

type SubDomainListReq struct {
	RootDomain string `json:"rootdomain"`
	PageNumber int	`json:"pagenumber"`
	PageSize int	`json:"pagesize"`
}

type SubDomainDesc struct {
	Domain string `json:"domain"`
	RegTime int64 `json:"regtime"`
	ExpireTime int64 `json:"expiretime"`
	Owner string `json:"owner"`
}

type SubDomainResp struct {
	State int `json:"state"`
	TotalPage int `json:"totalpage"`
	PageNumber int			`json:"pagenumber"`
	PageSize int			`json:"pagesize"`
	SubDomains []*SubDomainDesc `json:"subdomains"`
}

func NewSubDomainList() *SubDomainList {
	return &SubDomainList{}
}

func (sbl *SubDomainList)ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	req:=&SubDomainListReq{}
	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	rootDomain:=req.RootDomain
	if rootDomain == "" || req.PageNumber < 1{
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	rds:=strings.Split(rootDomain,".")
	for i:=0;i<len(rds);i++{
		if rds[i] == ""{
			w.WriteHeader(500)
			fmt.Fprintf(w, "{}")
			return
		}
	}

	cnt:=0
	b:= (req.PageNumber -1)*req.PageSize
	e:=req.PageNumber*req.PageSize

	var sdds []*SubDomainDesc

	for _,v:=range DataSync.Records{
		if strings.HasSuffix(v.GetName(),rootDomain){
			if cnt >=b && cnt < e{
				sdd := &SubDomainDesc{}
				sdd.Owner = v.GetOwner()
				sdd.ExpireTime = v.GetExpire()
				sdd.Domain = v.GetName()
				sdd.RegTime = v.GetRegTime()
				sdds = append(sdds,sdd)
			}
			cnt ++
		}
	}

	resp:=&SubDomainResp{}
	if len(sdds) == 0{
		resp.State = 0
	}else{
		resp.State = 1
		resp.SubDomains = sdds
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