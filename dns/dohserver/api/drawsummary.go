package api

import "math/big"

type DrawSummary struct {
}

type DrawSummaryReq struct {
	Wallet string `json:"wallet"`
}

type DrawSummaryResp struct {
	Wallet              string  `json:"wallet"`
	TotalWDrawTimes     int     `json:"totalwdrawtimes"`
	TotalWait2WDraw     big.Int `json:"totalwait2wdraw"`
	TotalWDrawed        big.Int `json:"totalwdrawed"`
	TotalMinerEarned    big.Int `json:"totalminerearned"`
	TotalOwnerEarned    big.Int `json:"totaoownerearned"`
	Wait2WDrawFromMiner big.Int `json:"wait2wdrawfromminer"`
	Wait2WDrawFromOwner big.Int `json:"wait2wdrawfromowner"`
	TotalReceipts       int	    `json:"totalreceipts"`
}




