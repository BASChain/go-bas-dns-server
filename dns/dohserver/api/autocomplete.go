package api

import (
	"encoding/json"
	"fmt"
	"github.com/BASChain/go-bas/DataSync"
	"io/ioutil"
	"net/http"
	"strings"
)

type AutoComplete struct {
}

type AutoCompleteReq struct {
	Text string `json:"text"`
}

type AutoCompleteResp struct {
	State int      `json:"state"`
	Data  []string `json:"data"`
}

func NewAutoComplete() *AutoComplete {
	return &AutoComplete{}
}

func (ac *AutoComplete) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	acr := &AutoCompleteReq{}

	err = json.Unmarshal(body, acr)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	searchText := acr.Text
	var data []string

	for _, r := range DataSync.Records {
		if strings.Contains(r.GetName(), searchText) {
			data = append(data, r.GetName())
		}
	}

	resp := &AutoCompleteResp{}

	if len(data) > 0 {
		resp.State = 1
	} else {
		resp.State = 0
	}

	resp.Data = data

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
