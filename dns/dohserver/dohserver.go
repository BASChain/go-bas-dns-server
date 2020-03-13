package dohserver

import (
	"context"
	"fmt"
	"github.com/BASChain/go-bas-dns-server/config"
	"github.com/BASChain/go-bas-dns-server/dns/server"
	"github.com/m13253/dns-over-https/json-dns"
	"github.com/miekg/dns"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"github.com/BASChain/go-bas-dns-server/dns/dohserver/api"
	"path"
)

const (
	VERSION    = "1.1.1"
	USER_AGENT = "DNS-over-HTTPS/" + VERSION + " (github.com/BASChain/go-bas-dns-server)"
)

type DohServer struct {
	dohServer    *http.Server
	dohsServer   *http.Server
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

var (
	gdohserver *DohServer

	dohoncelock sync.Mutex
)

func GetDohDaemonServer() *DohServer {
	if gdohserver == nil {
		dohoncelock.Lock()
		defer dohoncelock.Unlock()

		if gdohserver == nil {
			gdohserver = NewDohServers()
		}
	}
	return gdohserver
}

const(
	TotalPath string = "getDomainTotal"
	DomainList string = "getDomainList"
	AutoComplete string = "autocomplete"
	FreeEth   string = "freeEth"
	FreeBas   string = "freeBas"
)

func NewDohServers() *DohServer {
	cfg := config.GetBasDCfg()
	tv := cfg.TimeOut
	if tv == 0 {
		tv = 10
	}

	timeout := time.Duration(tv) * time.Second

	server := &DohServer{}

	addr := ":" + strconv.Itoa(cfg.DohServerPort)

	server.dohServer = &http.Server{Addr: addr}

	saddr := ":" + strconv.Itoa(cfg.DohsServerPort)
	server.dohsServer = &http.Server{Addr: saddr}

	mux := http.NewServeMux()
	mux.Handle(cfg.DnsPath, &DohServer{})

	mux.Handle(path.Join(cfg.BasApi,TotalPath),api.NewDomainTotal())
	mux.Handle(path.Join(cfg.BasApi,DomainList),api.NewDomainList())
	mux.Handle(path.Join(cfg.BasApi,AutoComplete),api.NewAutoComplete())
	mux.Handle(path.Join(cfg.BasApi,FreeEth),api.NewFreeEth())
	mux.Handle(path.Join(cfg.BasApi,FreeBas),api.NewFreeBas())

	smux := http.NewServeMux()
	smux.Handle(cfg.DnsPath, &DohServer{})
	smux.Handle(path.Join(cfg.BasApi,TotalPath),api.NewDomainTotal())
	smux.Handle(path.Join(cfg.BasApi,DomainList),api.NewDomainList())
	smux.Handle(path.Join(cfg.BasApi,AutoComplete),api.NewAutoComplete())
	smux.Handle(path.Join(cfg.BasApi,FreeEth),api.NewFreeEth())
	smux.Handle(path.Join(cfg.BasApi,FreeBas),api.NewFreeBas())

	server.dohServer.Handler = http.Handler(mux)

	server.dohsServer.Handler = http.Handler(smux)

	server.udpClient = &dns.Client{Net: "udp", UDPSize: dns.DefaultMsgSize, Timeout: timeout}
	server.tcpClient = &dns.Client{Net: "tcp", Timeout: timeout}
	server.tcpClientTls = &dns.Client{Net: "tcp-tls", Timeout: timeout}

	return server

}

func (doh *DohServer) StartDaemon() {
	if doh.dohServer == nil {
		log.Fatal("No Server, Please Init first")
		return
	}

	cfg := config.GetBasDCfg()

	if cfg.GetCertFile() != "" && cfg.GetKeyFile() != "" {
		log.Println("DOHS Server Start at :", cfg.DohsServerPort)
		go doh.dohsServer.ListenAndServeTLS(cfg.GetCertFile(), cfg.GetKeyFile())
	}

	log.Println("DOH Server Start at :", cfg.DohServerPort)
	doh.dohServer.ListenAndServe()

}

func (doh *DohServer) ShutDown() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	doh.dohServer.Shutdown(ctx)
	doh.dohsServer.Shutdown(ctx)
}

func (doh *DohServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS, POST")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.Header().Set("Server", USER_AGENT)
	w.Header().Set("X-Powered-By", USER_AGENT)

	if r.Method == "OPTIONS" {
		w.Header().Set("Content-Length", "0")
		return
	}

	if r.Form == nil {
		const maxMemory = 32 << 20 // 32 MB
		r.ParseMultipartForm(maxMemory)
	}
	contentType := r.Header.Get("Content-Type")
	if ct := r.FormValue("ct"); ct != "" {
		contentType = ct
	}
	if contentType == "" {
		// Guess request Content-Type based on other parameters
		if r.FormValue("name") != "" {
			contentType = "application/dns-json"
		} else if r.FormValue("dns") != "" {
			contentType = "application/dns-message"
		}
	}

	var responseType string
	for _, responseCandidate := range strings.Split(r.Header.Get("Accept"), ",") {
		responseCandidate = strings.SplitN(responseCandidate, ";", 2)[0]
		if responseCandidate == "application/json" {
			responseType = "application/json"
			break
		} else if responseCandidate == "application/dns-udpwireformat" {
			responseType = "application/dns-message"
			break
		} else if responseCandidate == "application/dns-message" {
			responseType = "application/dns-message"
			break
		}
	}
	if responseType == "" {
		// Guess response Content-Type based on request Content-Type
		if contentType == "application/dns-json" {
			responseType = "application/json"
		} else if contentType == "application/dns-message" {
			responseType = "application/dns-message"
		} else if contentType == "application/dns-udpwireformat" {
			responseType = "application/dns-message"
		}
	}

	var req *DNSRequest
	if contentType == "application/dns-json" {
		req = doh.parseRequestGoogle(ctx, w, r)
	} else if contentType == "application/dns-message" {
		req = doh.parseRequestIETF(ctx, w, r)
	} else if contentType == "application/dns-udpwireformat" {
		req = doh.parseRequestIETF(ctx, w, r)
	} else {
		jsonDNS.FormatError(w, fmt.Sprintf("Invalid argument value: \"ct\" = %q", contentType), 415)
		return
	}
	if req.errcode == 444 {
		return
	}
	if req.errcode != 0 {
		jsonDNS.FormatError(w, req.errtext, req.errcode)
		return
	}

	req = doh.patchRootRD(req)

	var err error
	req, err = doh.doDNSQuery(ctx, req)
	if err != nil {
		jsonDNS.FormatError(w, fmt.Sprintf("DNS query failure (%s)", err.Error()), 503)
		return
	}

	if responseType == "application/json" {
		doh.generateResponseGoogle(ctx, w, r, req)
	} else if responseType == "application/dns-message" {
		doh.generateResponseIETF(ctx, w, r, req)
	} else {
		panic("Unknown response Content-Type")
	}

}

func (doh *DohServer) findClientIP(r *http.Request) net.IP {
	noEcs := r.URL.Query().Get("no_ecs")
	if strings.ToLower(noEcs) == "true" {
		return nil
	}

	XForwardedFor := r.Header.Get("X-Forwarded-For")
	if XForwardedFor != "" {
		for _, addr := range strings.Split(XForwardedFor, ",") {
			addr = strings.TrimSpace(addr)
			ip := net.ParseIP(addr)
			if jsonDNS.IsGlobalIP(ip) {
				return ip
			}
		}
	}
	XRealIP := r.Header.Get("X-Real-IP")
	if XRealIP != "" {
		addr := strings.TrimSpace(XRealIP)
		ip := net.ParseIP(addr)
		if jsonDNS.IsGlobalIP(ip) {
			return ip
		}
	}
	remoteAddr, err := net.ResolveTCPAddr("tcp", r.RemoteAddr)
	if err != nil {
		return nil
	}
	if ip := remoteAddr.IP; jsonDNS.IsGlobalIP(ip) {
		return ip
	}
	return nil
}

func (doh *DohServer) patchRootRD(req *DNSRequest) *DNSRequest {
	for _, question := range req.request.Question {
		if question.Name == "." {
			req.request.RecursionDesired = true
		}
	}
	return req
}

func (doh *DohServer) indexQuestionType(msg *dns.Msg, qtype uint16) int {
	for i, question := range msg.Question {
		if question.Qtype == qtype {
			return i
		}
	}
	return -1
}

func (doh *DohServer) doDNSQuery(ctx context.Context, req *DNSRequest) (resp *DNSRequest, err error) {

	msg := req.request
	if msg == nil || len(msg.Question) == 0 {
		req.response = server.DeriveMsg(req.request, dns.RcodeFormatError)
		return req, nil
	}

	q := msg.Question[0]
	if q.Qclass != dns.ClassINET {
		req.response = server.DeriveMsg(req.request, dns.RcodeNotImplemented)
		return req, nil
	}

	switch q.Qtype {
	case dns.TypeA:
		req.response, err = server.BCReplayTypeA(req.request, q)

		if err != nil {
			req.response, err = server.BCReplyTraditionTypeA(msg)
			if req.response != nil && req.request != nil && len(req.request.Question)>0{
				req.response.Answer = append(req.response.Answer,server.BuildNullAnswer(req.request.Question[0]))
			}
		}
	case server.TypeBCAddr:
		req.response, err = server.BCReplyTypeBCA(req.request, q)

	default:
		req.response, err = server.BCReplyTraditionTypeA(msg)
		if req.response != nil && req.request != nil && len(req.request.Question)>0{
			req.response.Answer = append(req.response.Answer,server.BuildNullAnswer(req.request.Question[0]))
		}
	}

	return req, nil
}
