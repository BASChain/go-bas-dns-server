package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/BASChain/go-bas/Bas_Ethereum"
	"github.com/BASChain/go-bas/DataSync"
	"encoding/hex"
)

type DomainInfo struct {

}

type DomainInfoReq struct {
	DomainName string `json:"domainname"`
}

type DomainRecordInfo struct {
	Name string `json:"name"`
	Ipv4 string `json:"ipv4"`
	Ipv6 string `json:"ipv6"`
	BCAddr string `json:"bcaddr"`
	Alias string `json:"alias"`
	ExtraInfo string `json:"extrainfo"`
	DomainHash string `json:"domainhash,omitempty"`
}

type DomainInfoResp struct {
	State int `json:"state"`
	DnsInfo *DomainRecordInfo `json:"dnsinfo"`
	AssetInfo *RegDomainRecord `json:"assetinfo"`
}

func NewDomainInfo() *DomainInfo {
	return &DomainInfo{}
}

func (rd *DomainInfo)ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	req := &DomainInfoReq{}

	err = json.Unmarshal(body,req)
	if err!=nil{
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	dhash := Bas_Ethereum.GetHash(req.DomainName)

	resp:=&DomainInfoResp{}

	d,ok:=DataSync.Records[dhash]
	if !ok{
		resp.State = 0
	}else{
		dnsinfo:=&DomainRecordInfo{}
		dnsinfo.Name = d.GetName()
		dnsinfo.Ipv4 = d.GetIPv4Str()
		dnsinfo.Ipv6 = d.GetIpv6Str()
		dnsinfo.Alias = d.GetAliasName()
		dnsinfo.BCAddr = d.GetBCAddr()
		dnsinfo.ExtraInfo = d.GetExtraInfo()
		dnsinfo.DomainHash = "0x"+hex.EncodeToString(dhash[:])

		assetinfo:=&RegDomainRecord{}

		assetinfo.RIsPureA = d.GetIsPureA()
		assetinfo.IsRoot = d.GetIsRoot()
		assetinfo.Owner = d.GetOwner()
		assetinfo.Name = d.GetName()
		assetinfo.Expire = d.GetExpire()
		assetinfo.ROpenToPublic = d.GetOpenStatus()
		assetinfo.RisCustomed = d.GetIsCustomed()
		assetinfo.RcustomePrice = d.GetCustomedPrice()
		assetinfo.DomainHash = "0x"+hex.EncodeToString(dhash[:])

		if !assetinfo.IsRoot {
			roothash := d.GetParentHash()
			d,ok=DataSync.Records[roothash]
			if ok{
				r1:=&RegDomainRecord{}
				r1.RIsPureA = d.GetIsPureA()
				r1.IsRoot = d.GetIsRoot()
				r1.Owner = d.GetOwner()
				r1.Name = d.GetName()
				r1.Expire = d.GetExpire()
				r1.ROpenToPublic = d.GetOpenStatus()
				r1.RisCustomed = d.GetIsCustomed()
				r1.RcustomePrice = d.GetCustomedPrice()
				r1.DomainHash = "0x"+hex.EncodeToString(roothash[:])
				assetinfo.ParentDomain = r1
			}
		}
		resp.State = 1

		resp.DnsInfo = dnsinfo
		resp.AssetInfo = assetinfo

	}

	var bresp []byte

	bresp,err =json.Marshal(*resp)
	if err != nil{
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	w.WriteHeader(200)
	w.Write(bresp)

}




