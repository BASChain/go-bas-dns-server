package api

import (
	"net/http"
	"github.com/BASChain/go-bas/Bas_Ethereum"
	"encoding/json"
	"fmt"
)

type BasSettings struct {

}


type BasSettingsResp struct {
	ARootGas string `json:"arootgas"`
	BRootGas string `json:"brootgas"`
	SubGas   string `json:"subgas"`
	CustomPriceGas string `json:"custompricegas"`
	MaxYear int64	`json:"maxyear"`
	RareTypeLength int64	`json:"raretypelength"`
}

func NewBasSettings()  *BasSettings {
	return &BasSettings{}
}

func (bs *BasSettings)ServeHTTP(w http.ResponseWriter, r *http.Request)   {

	if r.URL.Query().Get("refresh") == "true"{
		Bas_Ethereum.Settings()
		w.WriteHeader(200)
		w.Write([]byte("Refresh Success"))

		return
	}

	resp:=&BasSettingsResp{}

	Bas_Ethereum.CheckSettings()
	resp.ARootGas = Bas_Ethereum.GetARootGas()
	resp.BRootGas = Bas_Ethereum.GetBRootGas()
	resp.SubGas   = Bas_Ethereum.GetSubGas()
	resp.CustomPriceGas = Bas_Ethereum.GetCustomPriceGas()
	resp.MaxYear = Bas_Ethereum.GetMaxYear()
	resp.RareTypeLength = Bas_Ethereum.GetRareTypeLength()


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