package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/BASChain/go-bas-dns-server/dns/mem"
	"math/big"
	"github.com/BASChain/go-bas/Miner"
	"github.com/kprc/nbsnetwork/common/list"
	"github.com/pkg/errors"
)

type MiningDetail struct {

}


type MiningDetailReq struct {
	Wallet string `json:"wallet"`
	PageNumber int `json:"pagenumber"`
	PageSize int `json:"pagesize"`
}



type MiningDetailItem struct {
	ReceiptHash string `json:"receipthash"`
	RootDomainName string `json:"rootdomainname"`
	OpName string `json:"opname"`
	FromDomainName string `json:"fromdomainname"`
	ToMiner string `json:"tominer"`
	ToOwner string `json:"toowner"`
}

type MiningDetailResp struct {
	State      int                `json:"state"`
	TotalPage  int                `json:"totalpage"`
	PageNumber int                `json:"pagenumber"`
	PageSize   int                `json:"pagesize"`
	Mdis       []*MiningDetailItem `json:"mdis"`
}

func NewMiningDetail() *MiningDetail {
	return &MiningDetail{}
}

func calc4Miner(mdi []*MiningDetailItem,lMiner list.List )  {
	for i:=0;i<len(mdi);i++{
		lMiner.FindDo(mdi[i], func(arg interface{}, v interface{}) (ret interface{}, err error) {
			m:=arg.(*MiningDetailItem)
			p:=v.(*mem.ProfitItem)
			if m.ReceiptHash == p.GetReceiptHash().String(){
				var (
					alc *[4]big.Int
					sett Miner.Setting
				)
				if p.Allocation == nil{
					bn,txid:=p.GetTractId()
					alc,sett=mem.GetSetting(bn,txid,p.GetAllocTyp())
					p.Allocation = alc
					if alc != nil{
						p.MinerCnt = len(sett.Miners)
					}
				}

				if p.Allocation != nil{
					toMiner:=mem.PriceCalc4Miner(p.MinerCnt,p.Amount,&p.Allocation[p.AllocTyp])

					m.ToMiner = toMiner.String()

					return toMiner,nil
				}else{
					m.ToMiner = ""
					return nil,errors.New("calc error")
				}

			}
			return nil,errors.New("no find")

		})
	}


}

func (md *MiningDetail)  ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	req := &MiningDetailReq{}

	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	addr := common.HexToAddress(req.Wallet)

	store:=mem.GetMinerProfitStore()
	store.Lock()
	m:=store.GetProfitMiner(&addr)
	store.UnLock()

	resp:=&MiningDetailResp{}
	resp.PageSize = req.PageSize
	resp.PageNumber = req.PageNumber


	if m == nil{
		resp.State = 0

	}else{

		m.Lock()
		defer m.UnLock()

		lMiner:=m.GetProfitItem4MinerList()
		lOwner:=m.GetProfitItemList()

		resp.TotalPage = int(lOwner.Count())

		start:=(req.PageNumber-1)*req.PageSize
		end:=req.PageNumber*req.PageSize

		cursor:=lOwner.ListIteratorB(start,end)
		if cursor.Count() <= 0{
			resp.State = 0
		}else{
			for{
				n:=cursor.Next()
				if n == nil{
					break
				}
				ni:=n.(*mem.ProfitItem)
				mdi:=&MiningDetailItem{}
				mdi.ReceiptHash = ni.GetReceiptHash().String()
				r:=GetRecord(*ni.GetDomainHash())
				if r==nil{
					mdi.FromDomainName = ""
					mdi.RootDomainName = ""
				}else{
					mdi.FromDomainName = r.GetName()

					root:=GetRecord(r.GetParentHash())
					if root == nil{
						mdi.RootDomainName = ""
					}else {
						mdi.RootDomainName = root.GetName()
					}
				}

				var alc *[4]big.Int
				if ni.Allocation == nil{
					bn,txid:=ni.GetTractId()
					alc,_=mem.GetSetting(bn,txid,ni.GetAllocTyp())
				}else{
					alc = ni.Allocation
				}
				if alc != nil{
					mdi.ToOwner = mem.PriceCalc4Owner(ni.GetAmount(),&alc[Miner.ToRoot]).String()
				}


				mdi.OpName = ni.GetSrcTyp()

				resp.Mdis = append(resp.Mdis,mdi)
			}

			calc4Miner(resp.Mdis,lMiner)

		}
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













