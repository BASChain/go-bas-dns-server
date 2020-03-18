package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/BASChain/go-bas/DataSync"
	"strings"
)

type SellAutoComplete struct {

}


type SellAutoCompleteReq struct {
	Wallet string `json:"wallet"`
	Text string `json:"text"`
}


type DHPaire struct {
	DomainName string `json:"domainname"`
	WalletAddress  string `json:"walletaddress"`
	Expire int64 `json:"expire"`
}


type SellAutoCompleteResp struct {
	State int `json:"state"`
	DomainHashPair []DHPaire `json:"domainhashpair"`
}

func NewSellAutoComplete() *SellAutoComplete  {
	return &SellAutoComplete{}
}

func (sac *SellAutoComplete)ServeHTTP(w http.ResponseWriter, r *http.Request)   {
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

	req:=&SellAutoCompleteReq{}

	err = json.Unmarshal(body,req)
	if err!=nil{
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	searchText := req.Text

	var addr *common.Address

	if req.Wallet != ""{
		address := common.HexToAddress(req.Wallet)
		addr = &address
	}


	var dhp []DHPaire

	for _,r:=range DataSync.Records{
		if strings.Contains(r.GetName(),searchText){
			if (addr == nil || (*addr) != (*r.GetOwnerOrig())) && r.GetBCAddr() != ""{
				item:=DHPaire{}
				item.DomainName = r.GetName()
				item.WalletAddress = r.GetBCAddr()
				item.Expire = r.GetExpire()
				dhp = append(dhp,item)
			}
		}
	}

	resp := &SellAutoCompleteResp{}

	if len(dhp) == 0{
		resp.State = 0
	}else{
		resp.State = 1
		resp.DomainHashPair = dhp
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





