package config

import (
	"encoding/json"
	"github.com/BASChain/go-bas/service"
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"os"
	"path"
	"sync"
)

const (
	BASD_HomeDir      = ".basd"
	BASD_CFG_FileName = "basd.json"
)

var EthNet string

type BasContactConfig struct {
	TokenAddr     string   `json:"tokenaddr"`
	OwnershipAddr string   `json:"ownershipaddr"`
	AssetAddr     string   `json:"assetaddr"`
	DNSAddr       string   `json:"dnsaddr"`
	OANNAddr      string   `json:"oannaddr"`
	MinerAddr     string   `json:"mineraddr"`
	MarketAddr    string   `json:"marketaddr"`
	RemoteServers []string `json:"remoteservers"`
}

type BASDConfig struct {
	UpdPort        int              `json:"updport"`
	TcpPort        int              `json:"tcpport"`
	RopstenNAP     string           `json:"ropstennap"`
	TokenAddr      string           `json:"tokenaddr"`
	MgrAddr        string           `json:"mgraddr"`
	CmdListenPort  string           `json:"cmdlistenport"`
	ResolvDns      []string         `json:"resolvdns"`
	DohServerPort  int              `json:"dohserverport"`
	DohsServerPort int              `json:"dohsserverport"`
	CertFile       string           `json:"certfile"`
	KeyFile        string           `json:"keyfile"`
	DnsPath        string           `json:"dnspath"`
	TimeOut        int              `json:"timeout"`
	TryTimes       int              `json:"trytimes"`
	BasApi         string           `json:"basapi"`
	DnsBasApi      string           `json:"dnsbasapi"`
	ContactApi     string           `json:"contactapi"`
	MyWalletApi    string           `json:"mywalletapi"`
	MarketApi      string           `json:"marketapi"`
	MinerApi       string           `json:"minerapi"`
	FreeEthAmount  string           `json:"freetokenamount"`
	FreeBasAmount  string           `json:"freebasamount"`
	TestContactCfg BasContactConfig `json:"testcontactcfg"`
	MainContactCfg BasContactConfig `json:"maincontactcfg"`
}

var (
	bascfgInst     *BASDConfig
	bascfgInstLock sync.Mutex
)

func (bc *BASDConfig) InitCfg() *BASDConfig {
	bc.UpdPort = 53
	bc.TcpPort = 53
	bc.CmdListenPort = "127.0.0.1:59538"
	bc.ResolvDns = []string{"202.106.0.20", "8.8.8.8", "202.106.46.151", "8.8.4.4"}
	bc.DohServerPort = 8053
	bc.DohsServerPort = 8043
	bc.DnsPath = "/dns-query"
	bc.TimeOut = 10
	bc.TryTimes = 3
	bc.BasApi = "/api"
	bc.DnsBasApi = "/api/domain"
	bc.ContactApi = "/api/contact"
	bc.MarketApi = "/api/market"
	bc.MyWalletApi = "/api/mywallet"
	bc.MinerApi = "/api/miner"
	bc.FreeEthAmount = "10000000000000000"      //0.01eth
	bc.FreeBasAmount = "4000000000000000000000" //100bas
	bc.MainContactCfg = BasContactConfig{
		TokenAddr:     "0x105B1413461394148023FEB5bE3b4307448872d5",
		OwnershipAddr: "0x35D5FE9dfbED34e0d404A8073D8Ee9618E8dbC16",
		AssetAddr:     "0x36631a815bbecfb8947e814196DbF1768397d75b",
		DNSAddr:       "0xEc784426d352fF80E6c4192a10B009dc45e92DBD",
		OANNAddr:      "0x6a76585B037988281Aa2c80E6E42d689bA940Cef",
		MinerAddr:     "0xb685C02bF992c61c68393aF7fcD8F46833Fb6937",
		MarketAddr:    "0xa26fDE795d1f15768B588Fb6A9342129AC38C648",
		RemoteServers: []string{
			"wss://mainnet.infura.io/ws/v3/831ab04fa4964991b5fba5c52106d7b0",
			"wss://mainnet.infura.io/ws/v3/8b8db3cca50a4fcf97173b7619b1c4c3",
			//"ws://75.135.96.248:3334",
		},
	}
	bc.TestContactCfg = BasContactConfig{
		TokenAddr:     "0x9d0314f9Bacd569DCB22276867AAEeE1C8A87614",
		OwnershipAddr: "0x4b91b82bed39B1d946C9E3BC12ba09C2F22fd3ee",
		AssetAddr:     "0x2B1110a13183A7045C7BCE3ba0092Ff0de4FD241",
		DNSAddr:       "0x8951f6B80b880E8A47d0d18000A4c90F288F61a3",
		OANNAddr:      "0x5e6B639843da8A9883aF8055C71D21d7dd4c30C3",
		MinerAddr:     "0xCAB59645aE535A7b5a4f81d8D17E2fe0d2Cf4687",
		MarketAddr:    "0xA32ccce4B7aB28d3Ce40BBa03A2748bCbe4544dB",
		RemoteServers: []string{
			"wss://ropsten.infura.io/ws/v3/831ab04fa4964991b5fba5c52106d7b0",
			"wss://ropsten.infura.io/ws/v3/8b8db3cca50a4fcf97173b7619b1c4c3",
			"ws://75.135.96.248:3334",
		},
	}
	return bc
}

func (bc *BASDConfig) SettingNet() {
	var c *BasContactConfig

	if EthNet == "main" {
		c = &bc.MainContactCfg
	}

	if EthNet == "test" {
		c = &bc.TestContactCfg
	}

	log.Println("Current Eth Net is", EthNet)

	if c != nil {

		service.ChangeAccessPoint(c.RemoteServers)

		service.ChangeContractAddresses(c.TokenAddr, c.OwnershipAddr, c.AssetAddr, c.DNSAddr, c.OANNAddr, c.MinerAddr, c.MarketAddr)

	} else {
		panic("Setting eth net failed")
	}
}

func (bc *BASDConfig) Load() *BASDConfig {
	if !tools.FileExists(GetBASDCFGFile()) {
		return nil
	}

	jbytes, err := tools.OpenAndReadAll(GetBASDCFGFile())
	if err != nil {
		log.Println("load file failed", err)
		return nil
	}

	//bc1:=&BASDConfig{}

	err = json.Unmarshal(jbytes, bc)
	if err != nil {
		log.Println("load configuration unmarshal failed", err)
		return nil
	}

	return bc

}

func newBasDCfg() *BASDConfig {

	bc := &BASDConfig{}

	bc.InitCfg()

	return bc
}

func GetBasDCfg() *BASDConfig {
	if bascfgInst == nil {
		bascfgInstLock.Lock()
		defer bascfgInstLock.Unlock()
		if bascfgInst == nil {
			bascfgInst = newBasDCfg()
		}
	}

	return bascfgInst
}

func PreLoad() *BASDConfig {
	bc := &BASDConfig{}

	return bc.Load()
}

func LoadFromCfgFile(file string) *BASDConfig {
	bc := &BASDConfig{}

	bc.InitCfg()

	bcontent, err := tools.OpenAndReadAll(file)
	if err != nil {
		log.Fatal("Load Config file failed")
		return nil
	}

	err = json.Unmarshal(bcontent, bc)
	if err != nil {
		log.Fatal("Load Config From json failed")
		return nil
	}

	bascfgInstLock.Lock()
	defer bascfgInstLock.Unlock()
	bascfgInst = bc

	return bc

}

func LoadFromCmd(initfromcmd func(cmdbc *BASDConfig) *BASDConfig) *BASDConfig {
	bascfgInstLock.Lock()
	defer bascfgInstLock.Unlock()

	lbc := newBasDCfg().Load()

	if lbc != nil {
		bascfgInst = lbc
	} else {
		lbc = newBasDCfg()
	}

	bascfgInst = initfromcmd(lbc)

	return bascfgInst
}

func GetBASDHomeDir() string {
	curHome, err := tools.Home()
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(curHome, BASD_HomeDir)
}

func GetBASDCFGFile() string {
	return path.Join(GetBASDHomeDir(), BASD_CFG_FileName)
}

func (bc *BASDConfig) GetCertFile() string {

	if bc.CertFile == "" {
		return ""
	}

	cf := path.Join(GetBASDHomeDir(), bc.CertFile)

	if tools.FileExists(cf) {
		return cf
	} else {
		return ""
	}
}

func (bc *BASDConfig) GetKeyFile() string {
	if bc.KeyFile == "" {
		return ""
	}
	kf := path.Join(GetBASDHomeDir(), bc.KeyFile)
	if tools.FileExists(kf) {
		return kf
	} else {
		return ""
	}
}

func (bc *BASDConfig) Save() {
	jbytes, err := json.MarshalIndent(*bc, " ", "\t")

	if err != nil {
		log.Println("Save BASD Configuration json marshal failed", err)
	}

	if !tools.FileExists(GetBASDHomeDir()) {
		os.MkdirAll(GetBASDHomeDir(), 0755)
	}

	err = tools.Save2File(jbytes, GetBASDCFGFile())
	if err != nil {
		log.Println("Save BASD Configuration to file failed", err)
	}

}

func IsInitialized() bool {
	if tools.FileExists(GetBASDCFGFile()) {
		return true
	}

	return false
}
