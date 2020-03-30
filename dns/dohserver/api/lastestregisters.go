package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/kprc/nbsnetwork/common/list"
	"github.com/BASChain/go-bas/DataSync"
	"bytes"
	"strings"
)

type LatestRegisters struct {

}


func  cmp(v1 interface{},v2 interface{}) int{
	d1,d2:=v1.(*DataSync.DomainRecord),v2.(*DataSync.DomainRecord)

	return bytes.Compare(d1.Name,d2.Name)
}

func sort(v1 interface{},v2 interface{}) int  {
	d1,d2:=v1.(*DataSync.DomainRecord),v2.(*DataSync.DomainRecord)


	if d1.CommitBlock < d2.CommitBlock {
		return 1
	}

	return -1

}


type LatestRegistersReq struct {
	PageNumber int		`json:"pagenumber"`
	PageSize int		`json:"pagesize"`
	Top      int        `json:"top"`    //0 top level domain,
										//1 normal top level domain,
										//2 level 2 domain,
										//3 level 3 domain
										//...
										//257 for all domain expect 0 and 1 top level domain
										//258 for 0 and 1 top level domain
}


type LatestRegistersResp struct {
	State int			`json:"state"`
	TotalPage int 		`json:"totalpage"`
	PageNumber int		`json:"pagenumber"`
	PageSize int		`json:"pagesize"`
	LatestDomains []*DomainDetail `json:"latestdomains"`
}

func NewLatestRegisters() *LatestRegisters  {
	return &LatestRegisters{}
}

func filter(domain string, top int, rare bool) bool {
	if top == 0{
		return rare
	}

	if top == 1{
		if rare == true{
			return false
		}

		ds := strings.Split(domain,".")
		if len(ds) > 1{
			return false
		}else{
			return true
		}
	}

	if top == 257 {
		b1:=filter(domain,0,rare)
		b2:=filter(domain,1,rare)

		return !(b1||b2)
	}

	if top == 258 {
		b1 := filter(domain,0,rare)
		b2 := filter(domain,1,rare)

		return b1||b2
	}

	ds:=strings.Split(domain,".")
	if len(ds) == top{
		return true
	}

	return false
}

func (lr *LatestRegisters)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
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

	req:=&LatestRegistersReq{}
	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	if req.PageNumber < 1 || req.Top <0 || req.Top > 258{
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	var latestRegDomainList list.List
	latestRegDomainList = list.NewList(cmp)
	latestRegDomainList.SetSortFunc(sort)
	for _,v:=range DataSync.Records{
		//log.Println(string(v.Name),req.Top,v.GetIsRare(),filter(string(v.Name),req.Top,v.GetIsRare()))
		if filter(string(v.Name),req.Top,v.GetIsRare()) == true{
			latestRegDomainList.AddValueOrder(v)
		}

	}

	e:=(req.PageNumber)*req.PageSize
	s:=(req.PageNumber-1)*req.PageSize

	cursor:=latestRegDomainList.ListIterator(0)

	resp:=&LatestRegistersResp{}

	var latestDomains []*DomainDetail

	cnt := 0
	if cursor.Count() <= s {
		resp.State = 0
	}else{
		resp.State = 1
		for{
			d:=cursor.Next()
			if d == nil{
				break
			}
			if cnt >=s && cnt <e{
				ld:=Record2DomainDetail(d.(*DataSync.DomainRecord))
				latestDomains = append(latestDomains,ld)
			}
			cnt ++
		}
		resp.LatestDomains = latestDomains
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