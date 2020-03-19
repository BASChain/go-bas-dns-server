package api

import (
	"encoding/json"
	"fmt"
	"github.com/BASChain/go-bas/Bas_Ethereum"
	"github.com/BASChain/go-bas/DataSync"
	"io/ioutil"
	"net/http"
)

type RegDomain struct {
}

type RegDomainReq struct {
	DomainName string `json:"domainname"`
}

type RegDomainRecord struct {
	Name          string           `json:"name"`
	Expire        int64            `json:"expire"`
	Owner         string           `json:"owner"`
	IsRoot        bool             `json:"isroot"`
	ROpenToPublic bool             `json:"ropentopublic"`
	RisCustomed   bool             `json:"riscustomed"`
	RIsPureA      bool             `json:"rispurea"`
	RcustomePrice string           `json:"rcustomeprice"`
	DomainHash    string           `json:"domainhash,omitempty"`
	ParentDomain  *RegDomainRecord `json:"parentdomain"`
}

type RegDomainResp struct {
	State        int              `json:"state"`
	DomainRecord *RegDomainRecord `json:"domainrecord"`
}

func NewRegDomain() *RegDomain {
	return &RegDomain{}
}

func (rd *RegDomain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	req := &RegDomainReq{}

	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	dhash := Bas_Ethereum.GetHash(req.DomainName)

	resp := &RegDomainResp{}

	d, ok := DataSync.Records[dhash]
	if !ok {
		resp.State = 0
	} else {
		r := &RegDomainRecord{}
		r.RIsPureA = d.GetIsPureA()
		r.IsRoot = d.GetIsRoot()
		r.Owner = d.GetOwner()
		r.Name = d.GetName()
		r.Expire = d.GetExpire()
		r.ROpenToPublic = d.GetOpenStatus()
		r.RisCustomed = d.GetIsCustomed()
		r.RcustomePrice = d.GetCustomedPrice()

		if !r.IsRoot {
			roothash := d.GetParentHash()
			d, ok = DataSync.Records[roothash]
			if ok {
				r1 := &RegDomainRecord{}
				r1.RIsPureA = d.GetIsPureA()
				r1.IsRoot = d.GetIsRoot()
				r1.Owner = d.GetOwner()
				r1.Name = d.GetName()
				r1.Expire = d.GetExpire()
				r1.ROpenToPublic = d.GetOpenStatus()
				r1.RisCustomed = d.GetIsCustomed()
				r1.RcustomePrice = d.GetCustomedPrice()
				r.ParentDomain = r1
			}
		}
		resp.State = 1
		resp.DomainRecord = r

	}

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
