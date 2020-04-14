package api

import "net/http"

type LatestDealDomain struct {
}

func NewLatestDealDomain() *LatestDealDomain {
	return &LatestDealDomain{}
}

func (ldd *LatestDealDomain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	SH(w, r, 1)
}
