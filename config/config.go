package config

import (
	"encoding/json"
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

type BASDConfig struct {
	UpdPort        int      `json:"updport"`
	TcpPort        int      `json:"tcpport"`
	RopstenNAP     string   `json:"ropstennap"`
	TokenAddr      string   `json:"tokenaddr"`
	MgrAddr        string   `json:"mgraddr"`
	CmdListenPort  string   `json:"cmdlistenport"`
	ResolvDns      []string `json:"resolvdns"`
	DohServerPort  int      `json:"dohserverport"`
	DohsServerPort int      `json:"dohsserverport"`
	CertFile       string   `json:"certfile"`
	KeyFile        string   `json:"keyfile"`
	DnsPath        string   `json:"dnspath"`
	TimeOut        int      `json:"timeout"`
	TryTimes       int      `json:"trytimes"`
	BasApi         string   `json:"basapi"`
	DnsBasApi      string   `json:"dnsbasapi"`
	ContactApi     string   `json:"contactapi"`
	MyWalletApi    string   `json:"mywalletapi"`
	MarketApi      string   `json:"marketapi"`
	FreeEthAmount  string   `json:"freetokenamount"`
	FreeBasAmount  string   `json:"freebasamount"`
}

var (
	bascfgInst     *BASDConfig
	bascfgInstLock sync.Mutex
)

func (bc *BASDConfig) InitCfg() *BASDConfig {
	bc.UpdPort = 53
	bc.TcpPort = 53
	bc.CmdListenPort = "127.0.0.1:59527"
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
	bc.FreeEthAmount = "100000000000000000"     //0.1eth
	bc.FreeBasAmount = "1000000000000000000000" //100bas

	return bc
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

func GetKeyFile() string {
	return path.Join(GetBASDHomeDir(), "UTC--2020-03-11T06-56-52.423772000Z--33324a5ee0b35f17536ceda27274e88e76640f24")
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
