package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/kprc/nbsnetwork/common/list"
	"github.com/BASChain/go-bas/DataSync"
	"encoding/hex"
)

type FavoriteDomain struct {

}

type FavoriteItem struct {
	Owner string `json:"owner"`
	Name string `json:"name"`
	ExpireTime int64 `json:"expiretime"`
	RegTime int64 `json:"regtime"`
	Price string `json:"price"`
	Hash string `json:"hash"`
	SubDomainCount int `json:"subdomaincount"`
}

type FavoriteDomainReq struct {
	PageNumber int	`json:"pagenumber"`
	PageSize int	`json:"pagesize"`
}

type FavoriteDomainResp struct {
	State int `json:"state"`
	TotalPage int `json:"totalpage"`
	PageNumber int			`json:"pagenumber"`
	PageSize int			`json:"pagesize"`

	Domains []*FavoriteItem `json:"domains"`
}

func favoriteCmp(v1 ,v2 interface{}) int  {
	f1,f2:=v1.(*FavoriteItem),v2.(*FavoriteItem)

	if f1.Hash == f2.Hash{
		return  0
	}

	return 1

}

func favoriteSort(v1,v2 interface{}) int  {
	f1,f2:=v1.(*FavoriteItem),v2.(*FavoriteItem)

	if f1.SubDomainCount < f2.SubDomainCount{
		return 1
	}

	return -1

}

func favoriteFDo(arg interface{},v interface{}) (ret interface{},err error) {
	f2:=v.(*FavoriteItem)
	f2.SubDomainCount ++

	return arg,nil
}



func NewFavoriteDomain() *FavoriteDomain {
	return &FavoriteDomain{}
}

func (fd *FavoriteDomain)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
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

	req:=&FavoriteDomainReq{}
	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	favoriteList := list.NewList(favoriteCmp)
	favoriteList.SetSortFunc(favoriteSort)

	for k,v:=range DataSync.Records{
		if v.GetIsRoot(){
			item:=&FavoriteItem{}
			item.Hash = "0x"+hex.EncodeToString(k[:])
			item.Name = v.GetName()
			item.RegTime = v.GetRegTime()
			item.ExpireTime = v.GetExpire()
			item.Owner = v.GetOwner()
			item.Price = v.GetCustomedPrice()

			if l:=favoriteList.Find(item);l==nil{
				favoriteList.AddValue(item)
			}
			continue
		}

		root:=v.SRootHash
		if r,ok:=DataSync.Records[root];!ok{
			continue
		}else{
			item:=&FavoriteItem{}
			item.Hash = "0x"+hex.EncodeToString(root[:])
			item.Name = r.GetName()
			item.RegTime = r.GetRegTime()
			item.ExpireTime = r.GetExpire()
			item.Owner = r.GetOwner()
			item.Price = r.GetCustomedPrice()
			if l:=favoriteList.Find(item);l==nil{
				favoriteList.AddValue(item)
			}
			favoriteList.FindDo(item,favoriteFDo)
		}
	}
	favoriteList.Sort()

	cnt:=0
	b:=(req.PageNumber-1)*req.PageSize
	e:=req.PageNumber*req.PageSize

	resp:=&FavoriteDomainResp{}

	var ds []*FavoriteItem

	cursor:=favoriteList.ListIterator(e)
	if cursor.Count() <= b{
		resp.State = 0
	}else{
		resp.State = 1
		for{
			nxt:=cursor.Next()
			if nxt == nil{
				break
			}
			if cnt >=b && cnt <e{
				ds = append(ds,nxt.(*FavoriteItem))
			}
			cnt ++
		}
		resp.Domains = ds
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