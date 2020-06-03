package mem

import (
	"github.com/BASChain/go-bmail-resolver"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/pkg/errors"
	"log"
	"net"
	"sync"
)

type DomainA struct {
	ip     []net.IP
	cname  []string
	whence int64 //ms
}

var (
	daMem     map[string]*DomainA
	daMemLock sync.Mutex

	reqMem     map[string]struct{}
	reqMemLock sync.Mutex
)

func GetDomainA(domain string) (net.IP, string, error) {
	da, err := getDomainA(domain)
	if err != nil {
		reqDomainA(domain)
		return net.IPv4zero, "", err
	}

	if tools.GetNowMsTime()-da.whence >= 300000 {
		reqDomainA(domain)
	}

	if len(da.ip) > 0 {
		if len(da.cname) > 0 {
			return da.ip[0], da.cname[0], nil
		} else {
			return da.ip[0], "", nil
		}
	} else {
		if len(da.cname) > 0 {
			return net.IPv4zero, da.cname[0], nil
		}
	}

	return net.IPv4zero, "", errors.New("not found")

}

func init() {
	reqMem = make(map[string]struct{})
	daMem = make(map[string]*DomainA)
}

func getDomainA(domain string) (*DomainA, error) {
	daMemLock.Lock()
	defer daMemLock.Unlock()

	if da, ok := daMem[domain]; !ok {
		return nil, errors.New("not found")
	} else {
		return da, nil
	}

	return nil, errors.New("not found")
}

func reqDomainA(domain string) {

	if _, ok := reqMem[domain]; ok {
		return
	}

	reqMemLock.Lock()

	if _, ok := reqMem[domain]; ok {
		reqMemLock.Unlock()
		return
	}

	reqMem[domain] = struct{}{}

	reqMemLock.Unlock()

	go updateDomainA(domain)

}

func updateDomainA(domain string) {

	defer func() {
		reqMemLock.Lock()
		delete(reqMem, domain)
		reqMemLock.Unlock()
	}()

	nr := resolver.NewEthResolver(true)

	ips, cns, err := nr.DomainA3(domain)
	if err != nil {
		log.Println("eth resolver", domain, err)
		return
	}
	if len(ips) > 0 {
		log.Println("eth resolver", domain, ips[0].String())
	}
	if len(cns) > 0 {
		log.Println("eth resolver", domain, cns[0])
	}

	var (
		da *DomainA
		ok bool
	)

	daMemLock.Lock()
	defer daMemLock.Unlock()

	if da, ok = daMem[domain]; !ok {
		da = &DomainA{}
	}

	da.ip = ips
	da.cname = cns
	da.whence = tools.GetNowMsTime()

	daMem[domain] = da

	return

}
