package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/BASChain/go-bas/DataSync"
	"github.com/kprc/nbsnetwork/tools"
	"strings"
	"github.com/kprc/nbsnetwork/common/list"
)

type TopLevelDomains struct {

}


type TopLevelDomainsReq struct{
    Text string 		`json:"text"`
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

func tldCmp(v1,v2 interface{}) int  {
	d1,d2:=v1.(string),v2.(string)

	if d1 == d2{
		return 0
	}

	return 1
}

func tldSort(v1,v2 interface{}) int  {
	d1,d2:=v1.(string),v2.(string)
	if strings.Compare(d1,d2) >= 0{
		return 1
	}

	return -1
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

	tldList := list.NewList(tldCmp)
	tldList.SetSortFunc(tldSort)


	curTime:=tools.GetNowMsTime()/1000

	for _,r:=range DataSync.Records{

		if r.GetIsRoot() && r.GetIsRare() && r.GetExpire() > curTime && r.GetOpenStatus(){
			if req.Text == "" || (req.Text != "" && strings.Contains(r.GetName(),req.Text)){
				tldList.AddValueOrder(r.GetName())
			}
		}
	}
	resp:=&TopLevelDomainsResp{}
	cnt := 0
	rb := (req.PageNumber-1) * req.PageSize
	re := (req.PageNumber) * req.PageSize

	cursor:=tldList.ListIterator(0)
	if cursor.Count() <= rb{
		resp.State = 0
	}else{
		resp.State = 1
		for{
			nxt:=cursor.Next()
			if nxt == nil{
				break
			}
			if cnt >=rb && cnt <re{
				resp.Domains = append(resp.Domains,nxt.(string))
			}
			cnt ++
		}
	}
	resp.PageSize = req.PageSize
	resp.PageNumber = req.PageNumber
	resp.TotalCnt = cnt

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



