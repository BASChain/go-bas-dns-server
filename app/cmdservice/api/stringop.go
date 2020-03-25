package api

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/BASChain/go-bas-dns-server/app/cmdcommon"
	"github.com/BASChain/go-bas-dns-server/app/cmdpb"
	"github.com/BASChain/go-bas/Bas_Ethereum"
	"github.com/BASChain/go-bas/DataSync"
	"github.com/ethereum/go-ethereum/common"
	"net"
	"strconv"
)

type CmdStringOPSrv struct {
}

func (cso *CmdStringOPSrv) StringOpDo(cxt context.Context, so *cmdpb.StringOP) (*cmdpb.DefaultResp, error) {
	msg := ""
	switch so.Op {
	case cmdcommon.CMD_ASSET:
		msg = listAssets(so.Param)
	case cmdcommon.CMD_DOMAIN:
		msg = GetRecords(so.Param)
	default:
		return encapResp("Command Not Found"), nil
	}

	return encapResp(msg), nil
}

func listAssets(wallet string) string {
	msg := ""
	if wallet == "" {
		for k, ass := range DataSync.Assets {
			msg += "Wallet: " + k.String() + "\r\n"
			msg += getAssetInfo(ass)
			msg += "\r\n"
		}

		if msg == "" {
			msg = "No assets"
		}

		return msg
	}
	msg += "Wallet: " + wallet + "\r\n"
	addr := common.HexToAddress(wallet)
	if a, ok := DataSync.Assets[addr]; !ok {
		msg = "NotFound"
	} else {
		msg = getAssetInfo(a)
	}

	return msg
}

func getAssetInfo(domains []Bas_Ethereum.Hash) string {
	msg := ""

	for i := 0; i < len(domains); i++ {
		if dr, ok := DataSync.Records[domains[i]]; ok {
			msg += "DHash: " + hex.EncodeToString(domains[i][:])
			msg += getDomain(dr)
			msg += "\r\n"
		}
	}

	return msg
}

func GetRecords(r string) string {
	msg := ""
	if r == "" {
		for k, d := range DataSync.Records {
			msg += "DHash: " + hex.EncodeToString(k[:])
			msg += getDomain(d)
			msg += "\r\n"
		}

		return msg
	}

	hash := Bas_Ethereum.GetHash(r)
	msg += "DHash: " + r
	if n, ok := DataSync.Records[hash]; !ok {
		return "Domain Not Found"
	} else {
		msg = getDomain(n)
	}

	return msg

}

func getDomain(domain *DataSync.DomainRecord) string {
	msg := ""

	msg += fmt.Sprintf("   %-20s ", domain.GetName())
	ip := domain.GetIPv4Addr()
	msg += fmt.Sprintf("%-16s ", net.IPv4(ip[0], ip[1], ip[2], ip[3]).String())
	msg += fmt.Sprintf("%-12s ", strconv.FormatInt(domain.GetExpire(), 10))

	return msg
}
