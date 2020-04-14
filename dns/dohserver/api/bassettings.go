package api

import (
	"encoding/json"
	"fmt"
	"github.com/BASChain/go-bas/DataSync"
	"net/http"
)

type BasSettings struct {
}

type BasSettingsResp struct {
	ARootGas       string `json:"arootgas"`
	BRootGas       string `json:"brootgas"`
	SubGas         string `json:"subgas"`
	CustomPriceGas string `json:"custompricegas"`
	MaxYear        int64  `json:"maxyear"`
	RareTypeLength int64  `json:"raretypelength"`
}

func NewBasSettings() *BasSettings {
	return &BasSettings{}
}

func (bs *BasSettings) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Query().Get("refresh") == "true" {
		DataSync.Settings()
		w.WriteHeader(200)
		w.Write([]byte("Refresh Success"))

		return
	}

	resp := &BasSettingsResp{}

	DataSync.CheckSettings()
	resp.ARootGas = DataSync.GetARootGas()
	resp.BRootGas = DataSync.GetBRootGas()
	resp.SubGas = DataSync.GetSubGas()
	resp.CustomPriceGas = DataSync.GetCustomPriceGas()
	resp.MaxYear = DataSync.GetMaxYear()
	resp.RareTypeLength = DataSync.GetRareTypeLength()

	var bresp []byte
	var err error

	bresp, err = json.Marshal(*resp)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	w.WriteHeader(200)
	w.Write(bresp)
}
