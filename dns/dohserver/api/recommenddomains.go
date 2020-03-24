package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
	"github.com/BASChain/go-bas/DataSync"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/BASChain/go-bas/Bas_Ethereum"
	"encoding/hex"
)

type RecommendDomains struct {

}

type RecommendDomainsReq struct {
	SearchDomains string    `json:"searchdomains"`
	PageNumber int			`json:"pagenumber"`
	PageSize int			`json:"pagesize"`
}

type DomainDetail struct {
	DnsInfo   *DomainRecordInfo `json:"dnsinfo"`
	AssetInfo *RegDomainRecord  `json:"assetinfo"`
}

type RecommandDomain struct {
	RecommendName string `json:"recommendname"`
	RootDomain *DomainDetail	`json:"rootdomain"`
}

type RecommendDomainsResp struct {
	State int `json:"state"`
	TotalCnt int `json:"totalcnt"`
	PageNumber int `json:"pagenumber"`
	PageSize int `json:"pagesize"`
	Registered []*DomainDetail		`json:"registered"`
	Recommend []*RecommandDomain	`json:"recommend"`
}

func NewRecommendDomains() *RecommendDomains {
	return &RecommendDomains{}
}

func FindAllTopLevelDomain() []DataSync.DomainRecord {
	curTime := tools.GetNowMsTime()/1000
	var ds []DataSync.DomainRecord
	for _,r:=range DataSync.Records{
		if r.GetIsRoot() && r.GetIsPureA() && r.GetExpire() > curTime && r.GetOpenStatus(){
			ds = append(ds,r)
		}
	}

	return ds
}

func Record2DomainDetail(d *DataSync.DomainRecord) *DomainDetail {
	dd:=&DomainDetail{}

	if d == nil{
		return nil
	}

	dhash := Bas_Ethereum.GetHash(d.GetName())

	dnsinfo := &DomainRecordInfo{}

	dnsinfo.Name = d.GetName()
	dnsinfo.Ipv4 = d.GetIPv4Str()
	dnsinfo.Ipv6 = d.GetIpv6Str()
	dnsinfo.Alias = d.GetAliasName()
	dnsinfo.BCAddr = d.GetBCAddr()
	dnsinfo.ExtraInfo = d.GetExtraInfo()
	dnsinfo.DomainHash = "0x" + hex.EncodeToString(dhash[:])

	assetinfo := &RegDomainRecord{}

	assetinfo.RIsPureA = d.GetIsPureA()
	assetinfo.IsRoot = d.GetIsRoot()
	assetinfo.Owner = d.GetOwner()
	assetinfo.Name = d.GetName()
	assetinfo.Expire = d.GetExpire()
	assetinfo.ROpenToPublic = d.GetOpenStatus()
	assetinfo.RisCustomed = d.GetIsCustomed()
	assetinfo.RcustomePrice = d.GetCustomedPrice()
	assetinfo.DomainHash = "0x" + hex.EncodeToString(dhash[:])

	if !assetinfo.IsRoot {
		roothash := d.GetParentHash()
		d1, ok := DataSync.Records[roothash]
		if ok {
			r1 := &RegDomainRecord{}
			r1.RIsPureA = d1.GetIsPureA()
			r1.IsRoot = d1.GetIsRoot()
			r1.Owner = d1.GetOwner()
			r1.Name = d1.GetName()
			r1.Expire = d1.GetExpire()
			r1.ROpenToPublic = d1.GetOpenStatus()
			r1.RisCustomed = d1.GetIsCustomed()
			r1.RcustomePrice = d1.GetCustomedPrice()
			r1.DomainHash = "0x" + hex.EncodeToString(roothash[:])
			assetinfo.ParentDomain = r1
		}
	}

	dd.AssetInfo = assetinfo
	dd.DnsInfo = dnsinfo

	return dd

}

func GetRootDomainRecord(domain string) *DataSync.DomainRecord  {
	if domain == ""{
		return nil
	}

	dhash := Bas_Ethereum.GetHash(domain)

	d,ok:=DataSync.Records[dhash]
	if !ok{
		return nil
	}

	return &d
}


func GetSubDomainRecord(domain string) *DataSync.DomainRecord {

	if domain == ""{
		return nil
	}
	dsegs := strings.Split(domain,".")
	if len(dsegs) == 1{
		return nil
	}
	for _,s:=range dsegs{
		if s == ""{
			return nil
		}
	}

	dhash := Bas_Ethereum.GetHash(domain)

	d,ok:=DataSync.Records[dhash]
	if !ok{
		return nil
	}

	return &d

}

func (rd *RecommendDomains)  ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	req := &RecommendDomainsReq{}

	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	searchDomain := req.SearchDomains

	if searchDomain == "" {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	domainSegs := strings.Split(searchDomain,".")
	if len(domainSegs) < 2 {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	if req.PageNumber < 1{
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}


	for _,d:=range domainSegs{
		if d == ""{
			w.WriteHeader(500)
			fmt.Fprintf(w, "{}")
			return
		}
	}

	rootDomain := domainSegs[len(domainSegs)-1]

	prefixDomain := ""
	for i:=0;i<len(domainSegs)-1;i++{
		if prefixDomain != ""{
			prefixDomain += "."+domainSegs[i]
		}else{
			prefixDomain += domainSegs[i]
		}
	}

	cnt := 0
	cb:=(req.PageNumber-1)*req.PageSize
	ce:=req.PageNumber*req.PageSize

	var Registered []*DomainDetail
	var Recommend []*RecommandDomain

	d:=GetSubDomainRecord(searchDomain)
	if d!=nil{
		if cnt >=cb && cnt <ce{
			dd:=Record2DomainDetail(d)
			Registered = append(Registered,dd)
		}
		cnt ++
	}else{
		r:=GetRootDomainRecord(rootDomain)
		if r!=nil{
			if cnt >=cb && cnt <ce{
				dr:=Record2DomainDetail(r)
				rec:=&RecommandDomain{}
				rec.RecommendName = searchDomain
				rec.RootDomain = dr
				Recommend = append(Recommend,rec)
			}
			cnt ++
		}
	}

	roots:=FindAllTopLevelDomain()
	curTime:=tools.GetNowMsTime()/1000
	for i:=0;i<len(roots);i++{
		r:=roots[i]
		rname:=r.GetName()
		if rname == rootDomain{
			continue
		}
		name:=prefixDomain+"."+rname
		d:=GetSubDomainRecord(name)
		if d == nil{
			if cnt >=cb && cnt <ce{
				dr:=Record2DomainDetail(&(roots[i]))
				rec:=&RecommandDomain{}
				rec.RecommendName = name
				rec.RootDomain = dr
				Recommend = append(Recommend,rec)
			}
			cnt ++
		}else{
			if curTime > d.GetExpire(){
				continue
			}

			if cnt >=cb && cnt <ce{
				dd:=Record2DomainDetail(d)
				Registered = append(Registered,dd)
			}

			cnt ++
		}
	}

	resp:=&RecommendDomainsResp{}

	if len(Registered)==0 && len(Recommend)==0{
		resp.State = 0
	}else{
		resp.State = 1
	}
	resp.TotalCnt = cnt
	resp.PageNumber =req.PageNumber
	resp.PageSize = req.PageSize
	resp.Recommend = Recommend
	resp.Registered = Registered

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







