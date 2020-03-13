package server

import (
	"github.com/BASChain/go-bas-dns-server/config"

	"github.com/btcsuite/btcutil/base58"
	"github.com/miekg/dns"

	"log"
	"net"
	"strconv"
	"github.com/BASChain/go-bas-dns-server/dns/dohserver/api"

	"errors"
)

const (
	TypeBCAddr = 65
)

var (
	dnshandle dns.HandlerFunc
)

var (
	udpServer *dns.Server
	tcpServer *dns.Server
)



func sendErrMsg(w dns.ResponseWriter, msg *dns.Msg, errCode int) {
	m := DeriveMsg(msg, errCode)

	w.WriteMsg(m)
}

func DeriveMsg(msg *dns.Msg, errCode int) *dns.Msg {
	m := msg.Copy()
	m.Compress = true
	m.Response = true

	m.Rcode = errCode

	return m
}

func buildAnswer(ipv4 [4]byte, q dns.Question) []dns.RR {
	A := &dns.A{}

	A.Hdr.Name = q.Name
	A.Hdr.Rrtype = dns.TypeA
	A.Hdr.Class = dns.ClassINET
	A.Hdr.Ttl = 10
	A.Hdr.Rdlength = 4

	A.A = net.IPv4(ipv4[0], ipv4[1], ipv4[2], ipv4[3])

	log.Println("Request Name: ", q.Name, A.A.String())

	var rr []dns.RR

	rr = append(rr, A)

	return rr
}

func BuildNullAnswer(q dns.Question) dns.RR  {
	NULL:=&dns.NULL{}

	NULL.Hdr.Name = q.Name

	NULL.Hdr.Rrtype = dns.TypeNULL
	NULL.Hdr.Class = dns.ClassINET
	NULL.Hdr.Ttl = 10
	NULL.Hdr.Rdlength = uint16(len("TraditionSystemName"))

	NULL.Data="TraditionSystemName"

	return NULL
}


func replyTypA(w dns.ResponseWriter, msg *dns.Msg, q dns.Question) error {
	if m, err := BCReplayTypeA(msg, q); err != nil {
		return err
	} else {
		w.WriteMsg(m)
		return nil
	}

}

func BCReplayTypeA(msg *dns.Msg, q dns.Question) (resp *dns.Msg, err error) {
	qn := q.Name
	if qn[len(qn)-1] == '.' {
		qn = qn[:len(qn)-1]
	}

	//if bdr, err := BAS_Ethereum.QueryByString(qn); err != nil {
	//	return nil, errors.New("Not Found")
	//} else {
	//	dr := &DR{bdr}
	//	if dr.IntIPv4() == 0 {
	//		return nil, errors.New("Not Found")
	//	}
	//	m := msg.Copy()
	//	m.Compress = true
	//	m.Response = true
	//	m.Answer = buildAnswer(dr.IPv4, q)
	//
	//	return m, nil
	//}
	 dr:=api.QueryBasByDomainName(qn)
	 if dr == nil{
		return nil, errors.New("Not Found")
	 }

	 if dr.GetIPv4() == 0{
		 return nil, errors.New("Not Found")
	 }
	 m := msg.Copy()
	 m.Compress = true
	 m.Response = true
	 m.Answer = buildAnswer(dr.GetIPv4Addr(), q)

	 return m, nil

}

func replyTraditionTypA(w dns.ResponseWriter, msg *dns.Msg) {
	m, _ := BCReplyTraditionTypeA(msg)
	w.WriteMsg(m)
}

func BCReplyTraditionTypeA(msg *dns.Msg) (resp *dns.Msg, err error) {
	cnt := 0

	for {
		cnt++
		s := GetDns()

		if s == "" {
			m := DeriveMsg(msg, dns.RcodeServerFailure)
			return m, nil
		}

		if m, err := dns.Exchange(msg, s+":53"); err != nil {
			FailDns(s)
			log.Println("failed "+s+":53", msg.Question[0].Name)
			if cnt >= MaxTimes() {
				m := DeriveMsg(msg, dns.RcodeBadKey)
				return m, nil
			}
		} else {
			log.Println("success "+s+":53", msg.Question[0].Name)
			return m, nil
		}

	}

}

func replyTypPTR(w dns.ResponseWriter, msg *dns.Msg, q dns.Question) error {
	return nil
}

func replyTraditionTypPTR(w dns.ResponseWriter, msg *dns.Msg, q dns.Question) {
	return
}

func replyTypBCA(w dns.ResponseWriter, msg *dns.Msg, q dns.Question) {
	m, _ := BCReplyTypeBCA(msg, q)

	w.WriteMsg(m)
}

func BCReplyTypeBCA(msg *dns.Msg, q dns.Question) (resp *dns.Msg, err error) {
	qn := q.Name
	if qn[len(qn)-1] == '.' {
		qn = qn[:len(qn)-1]

	}

	var b []byte
	b = base58.Decode(qn)
	var barr [32]byte

	for i := 0; i < len(b); i++ {
		barr[i] = b[i]
	}

	//if bdr, err := BAS_Ethereum.QueryByBCAddress(barr); err != nil {
	//	m := DeriveMsg(msg, dns.RcodeBadKey)
	//	return m, nil
	//} else {
	//	dr := &DR{bdr}
	//	if dr.IntIPv4() == 0 {
	//		m := DeriveMsg(msg, dns.RcodeBadKey)
	//		return m, nil
	//	}
	//	m := msg.Copy()
	//	m.Compress = true
	//	m.Response = true
	//
	//	m.Answer = buildAnswer(dr.IPv4, q)
	//
	//	return m, nil
	//}

	m := DeriveMsg(msg, dns.RcodeBadKey)
	return m, nil


}

func DnsHandleTradition(w dns.ResponseWriter, msg *dns.Msg) {
	if len(msg.Question) == 0 {
		sendErrMsg(w, msg, dns.RcodeFormatError)
		return
	}
	q := msg.Question[0]

	if q.Qclass != dns.ClassINET {
		sendErrMsg(w, msg, dns.RcodeNotImplemented)
		return
	}

	switch q.Qtype {
	case dns.TypeA:
		//log.Println("type a reply", q.Name)
		if err := replyTypA(w, msg, q); err != nil {
			//log.Println("replyTypA error,begin reply tradition type a",q.Name)
			replyTraditionTypA(w, msg)
		}
	case dns.TypePTR:
		//if err:=replyTypPTR(w,msg,q);err!=nil{
		//	replyTraditionTypPTR(w,msg,q)
		//}
		replyTraditionTypA(w, msg)
	case TypeBCAddr:
		replyTypBCA(w, msg, q)
	default:
		//sendErrMsg(w,msg,dns.RcodeNotImplemented)
		//log.Println("Default reply",q.Name)
		replyTraditionTypA(w, msg)
		return
	}

}

func DNSServerDaemon() {
	cfg := config.GetBasDCfg()

	uport := cfg.UpdPort
	uaddr := ":" + strconv.Itoa(uport)

	dnshandle = DnsHandleTradition

	log.Println("DNS Server Start at udp", uaddr)

	udpServer = &dns.Server{}
	udpServer.Addr = uaddr
	udpServer.Handler = dnshandle
	udpServer.Net = "udp4"

	go udpServer.ListenAndServe()

	tport := cfg.TcpPort

	taddr := ":" + strconv.Itoa(tport)

	log.Println("DNS Server Start at tcp", taddr)

	tcpServer = &dns.Server{Addr: taddr, Net: "tcp4", Handler: dnshandle}

	tcpServer.ListenAndServe()
}

func DNSServerStop() {

	if udpServer != nil {
		udpServer.Shutdown()
		udpServer = nil
	}

	if tcpServer != nil {
		tcpServer.Shutdown()
		tcpServer = nil
	}

}
