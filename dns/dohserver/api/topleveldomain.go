package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/BASChain/go-bas/DataSync"
	"github.com/kprc/nbsnetwork/tools"
)

type TopLevelDomains struct {

}


type TopLevelDomainsReq struct{
	PageNumber int		`json:"pagenumber"`
	PageSize int		`json:"pagesize"`
}



type TopLevelDomainsResp struct {
	State int			`json:"state"`
	TotalCnt int		`json:"totalcnt"`
	PageNumber int		`json:"pagenumber"`
	PageSize int		`json:"pagesize"`
	Domains []string	`json:"domains"`
}

func NewTopLevelDomains() *TopLevelDomains {
	return &TopLevelDomains{}
}

func (tld *TopLevelDomains)ServeHTTP(w http.ResponseWriter,r *http.Request) {
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

	req:=&TopLevelDomainsReq{}
	err = json.Unmarshal(body,req)
	if err!=nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	if req.PageNumber < 1{
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	var domains []string

	cnt := 0
	rb := (req.PageNumber-1) * req.PageSize
	re := (req.PageNumber) * req.PageSize

	curTime:=tools.GetNowMsTime()/1000

	for _,r:=range DataSync.Records{

		if r.GetIsRoot() && r.GetIsRare() && r.GetExpire() > curTime && r.GetOpenStatus(){
			if cnt >=rb && cnt < re{
				domains = append(domains,r.GetName())
			}
			cnt ++
		}
	}

	resp:=&TopLevelDomainsResp{}
	if len(domains) == 0{
		resp.State = 0
	}else{
		resp.State = 1

	}

	resp.TotalCnt = cnt
	resp.PageSize = req.PageSize
	resp.PageNumber = req.PageNumber
	resp.Domains = domains


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



