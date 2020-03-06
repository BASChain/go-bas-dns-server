package api

import "net/http"

type AutoComplete struct {

}

func NewAutoComplete() *AutoComplete {
	return &AutoComplete{}
}

func (ac *AutoComplete)ServeHTTP(w http.ResponseWriter, r *http.Request)  {


}