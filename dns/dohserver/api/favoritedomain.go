package api

import "net/http"

type FavoriteDomain struct {

}


type FavoriteDomainReq struct {

}

type FavoriteDomainResp struct {

}

func NewFavoriteDomain() *FavoriteDomain {
	return &FavoriteDomain{}
}

func (fd *FavoriteDomain)ServeHTTP(w http.ResponseWriter, r *http.Request)  {

}