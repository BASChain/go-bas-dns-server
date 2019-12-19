package dohserver

import (
	"github.com/BASChain/go-bas-dns-server/config"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

type DohServer struct {
	dohServer    *http.Server
	udpClient    *dns.Client
	tcpClient    *dns.Client
	tcpClientTls *dns.Client
}

type DNSRequest struct {
	request         *dns.Msg
	response        *dns.Msg
	transactionID   uint16
	currentUpstream string
	isTailored      bool
	errcode         int
	errtext         string
}

func NewDohServer() *DohServer {
	cfg := config.GetBasDCfg()
	tv := cfg.TimeOut
	if tv == 0 {
		tv = 10
	}

	timeout := time.Duration(tv) * time.Second

	server := &DohServer{}

	addr := ":" + strconv.Itoa(cfg.DohServerPort)

	server.dohServer = &http.Server{Addr: addr}

	mux := http.NewServeMux()
	mux.HandleFunc(cfg.DnsPath, server.handlerFunc)

	server.dohServer.Handler = http.Handler(mux)

	server.udpClient = &dns.Client{Net: "udp", UDPSize: dns.DefaultMsgSize, Timeout: timeout}
	server.tcpClient = &dns.Client{Net: "tcp", Timeout: timeout}
	server.tcpClientTls = &dns.Client{Net: "tcp-tls", Timeout: timeout}

	return server

}

func (doh *DohServer) StartDaemon() error {
	if doh.dohServer == nil {
		return errors.New("No Server, Please Init first")
	}

	cfg := config.GetBasDCfg()
	return doh.dohServer.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
}

func (doh *DohServer) handlerFunc(w http.ResponseWriter, r *http.Request) {

}
