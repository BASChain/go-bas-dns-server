package server

import (
	"github.com/BASChain/go-bas-dns-server/config"

	"github.com/BASChain/go-bas-dns-server/lib/dns"
	"github.com/btcsuite/btcutil/base58"

	"github.com/BASChain/go-bas-dns-server/dns/dohserver/api"
	"log"
	"net"
	"strconv"

	"encoding/binary"
	"errors"
	"github.com/BASChain/go-bas-dns-server/dns/mem"
)

const (
	TypeBCAddr  = 65 //only for Question
	TypeBasAddr = 66 //for answer
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
	A := BuildAAnswer(ipv4, q)

	var rr []dns.RR

	rr = append(rr, A)

	return rr
}

func BuildAAnswer(ipv4 [4]byte, q dns.Question) dns.RR {
	A := &dns.A{}

	A.Hdr.Name = q.Name
	A.Hdr.Rrtype = dns.TypeA
	A.Hdr.Class = dns.ClassINET
	A.Hdr.Ttl = 10
	A.Hdr.Rdlength = 4

	A.A = net.IPv4(ipv4[0], ipv4[1], ipv4[2], ipv4[3])

	log.Println("Request Name: ", q.Name, A.A.String())

	return A
}

func BuildNullAnswer(q dns.Question, data string) dns.RR {
	NULL := &dns.NULL{}

	NULL.Hdr.Name = q.Name

	NULL.Hdr.Rrtype = dns.TypeNULL
	NULL.Hdr.Class = dns.ClassINET
	NULL.Hdr.Ttl = 10
	NULL.Hdr.Rdlength = uint16(len("TraditionSystemName"))

	if data == "" {
		NULL.Data = "TraditionSystemName"
	} else {
		NULL.Data = data
	}

	return NULL
}

func BuildCnameAnswer(cname string, q dns.Question) dns.RR {
	CNAME := &dns.CNAME{}

	CNAME.Hdr.Name = q.Name
	CNAME.Hdr.Rrtype = dns.TypeCNAME
	CNAME.Hdr.Class = dns.ClassINET
	CNAME.Hdr.Ttl = 10
	CNAME.Hdr.Rdlength = uint16(len(cname))

	CNAME.Target = cname

	return CNAME
}

func replyTypA(w dns.ResponseWriter, msg *dns.Msg, q dns.Question) error {
	if m, err := BCReplayTypeA2(msg, q); err != nil {
		return err
	} else {
		w.WriteMsg(m)
		return nil
	}

}

func BCReplayTypeA2(msg *dns.Msg, q dns.Question) (resp *dns.Msg, err error) {
	qn := q.Name
	if qn[len(qn)-1] == '.' {
		qn = qn[:len(qn)-1]
	}
	//log.Println("query :",qn)

	ip, cn, err := mem.GetDomainA(qn)
	if binary.BigEndian.Uint32(ip.To4()) == 0 {
		if cn != "" {
			m := msg.Copy()
			//m.Question[0].Name=dr.GetAliasName()
			m.Compress = true
			m.Response = true
			m.Answer = append(m.Answer, BuildCnameAnswer(cn, q))
			m.Answer = append(m.Answer, BuildNullAnswer(q, "AliasName"))
			//m.Answer = append(m.Answer,BuildAAnswer(dr.GetIPv4Addr(), q))
			//log.Println("response to client, type alias",qn)
			return m, nil
		}
	} else {
		m := msg.Copy()
		m.Compress = true
		m.Response = true
		var ipparam [4]byte
		copy(ipparam[:], ip)
		m.Answer = buildAnswer(ipparam, q)
		//log.Println("response to client, type a",qn)
		return m, nil
	}

	return nil, errors.New("No settings")

}

func BCReplayTypeA(msg *dns.Msg, q dns.Question) (resp *dns.Msg, err error) {
	qn := q.Name
	if qn[len(qn)-1] == '.' {
		qn = qn[:len(qn)-1]
	}
	//log.Println("query :",qn)

	dr := api.QueryBasByDomainName(qn)
	if dr == nil {
		//log.Println("query ",qn,"failed")
		return nil, errors.New("Not Found")
	}

	if dr.GetIPv4() == 0 {

		if dr.GetAliasName() != "" {
			m := msg.Copy()
			//m.Question[0].Name=dr.GetAliasName()
			m.Compress = true
			m.Response = true
			m.Answer = append(m.Answer, BuildCnameAnswer(dr.GetAliasName(), q))
			m.Answer = append(m.Answer, BuildNullAnswer(q, "AliasName"))
			//m.Answer = append(m.Answer,BuildAAnswer(dr.GetIPv4Addr(), q))
			//log.Println("response to client, type alias",qn)
			return m, nil
		}
	} else {
		m := msg.Copy()
		m.Compress = true
		m.Response = true
		m.Answer = buildAnswer(dr.GetIPv4Addr(), q)
		//log.Println("response to client, type a",qn)
		return m, nil
	}
	return nil, errors.New("No settings")

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
