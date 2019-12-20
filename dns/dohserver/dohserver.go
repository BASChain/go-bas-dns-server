package dohserver

import (
	"context"
	"fmt"
	"github.com/BASChain/go-bas-dns-server/config"
	"github.com/m13253/dns-over-https/json-dns"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"strconv"
	"strings"
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

const (
	VERSION    = "1.1.0"
	USER_AGENT = "DNS-over-HTTPS/" + VERSION + " (github.com/BASChain/go-bas-dns-server)"
)

func (doh *DohServer) StartDaemon() error {
	if doh.dohServer == nil {
		return errors.New("No Server, Please Init first")
	}

	cfg := config.GetBasDCfg()
	return doh.dohServer.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
}

func (doh *DohServer) handlerFunc(w http.ResponseWriter, r *http.Request) {
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

	// TODO(m13253): Make ctx work. Waiting for a patch for ExchangeContext from miekg/dns.
	/*
		numServers := len(s.conf.Upstream)
		for i := uint(0); i < s.conf.Tries; i++ {
			req.currentUpstream = s.conf.Upstream[rand.Intn(numServers)]

			upstream, t := addressAndType(req.currentUpstream)

			switch t {
			default:
				log.Printf("invalid DNS type %q in upstream %q", t, upstream)
				return nil, &configError{"invalid DNS type"}
				// Use DNS-over-TLS (DoT) if configured to do so
			case "tcp-tls":
				req.response, _, err = s.tcpClientTLS.Exchange(req.request, upstream)
			case "tcp", "udp":
				// Use TCP if always configured to or if the Query type dictates it (AXFR)
				if t == "tcp" || (s.indexQuestionType(req.request, dns.TypeAXFR) > -1) {
					req.response, _, err = s.tcpClient.Exchange(req.request, upstream)
				} else {
					req.response, _, err = s.udpClient.Exchange(req.request, upstream)
					if err == nil && req.response != nil && req.response.Truncated {
						log.Println(err)
						req.response, _, err = s.tcpClient.Exchange(req.request, upstream)
					}

					// Retry with TCP if this was an IXFR request and we only received an SOA
					if (s.indexQuestionType(req.request, dns.TypeIXFR) > -1) &&
						(len(req.response.Answer) == 1) &&
						(req.response.Answer[0].Header().Rrtype == dns.TypeSOA) {
						req.response, _, err = s.tcpClient.Exchange(req.request, upstream)
					}
				}
			}

			if err == nil {
				return req, nil
			}
			log.Printf("DNS error from upstream %s: %s\n", req.currentUpstream, err.Error())
		}
		return req, err

	*/

	return
}
