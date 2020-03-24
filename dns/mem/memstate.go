package mem

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type MemState struct {
	FreeState [2]int //1 Waiting result, 2 success, 3 failure
	WriteTime int64
}

var (
	lock      sync.Mutex
	m         map[common.Address]*MemState
	quit      chan int
	wg        sync.WaitGroup
	roundTime int64
)

const (
	BAS int = 0
	ETH int = 1

	WAITING int = 1
	SUCCESS int = 2
	FAILURE int = 3
)

func Update(addr common.Address, typ int, state int) error {

	if typ != BAS && typ != ETH {
		return errors.New("type error")
	}
	if state != WAITING && state != SUCCESS && state != FAILURE {
		return errors.New("state error")
	}

	lock.Lock()
	defer lock.Unlock()

	if m == nil {
		m = make(map[common.Address]*MemState)
	}

	if s, ok := m[addr]; !ok {
		s = &MemState{}
		s.FreeState[typ] = state
		s.WriteTime = tools.GetNowMsTime()
		m[addr] = s
	} else {
		s.FreeState[typ] = state
		s.WriteTime = tools.GetNowMsTime()
	}

	return nil

}

func GetState(addr common.Address, typ int) (int, error) {

	lock.Lock()
	defer lock.Unlock()

	if typ != BAS && typ != ETH {
		return 0, errors.New("type error")
	}

	if s, ok := m[addr]; !ok {
		return 0, errors.New("No Address")
	} else {
		return s.FreeState[typ], nil
	}
}

func MemStateStop() {
	quit <- 1
	wg.Wait()
}

func MemStateStart() {

	quit = make(chan int, 1)

	wg.Add(1)

	go memStateTimeOut()
}

func memStateTimeOut() {

	defer wg.Done()

	for {
		select {
		case <-quit:
			return
		default:

		}

		curTime := tools.GetNowMsTime()

		if curTime-roundTime < 300000 {
			time.Sleep(time.Second)
			continue
		}

		roundTime = curTime

		ks := make([]common.Address, 0)

		lock.Lock()

		if m != nil {
			for k, v := range m {
				if curTime-v.WriteTime > 300000 {
					ks = append(ks, k)
				}
			}
		}

		for i := 0; i < len(ks); i++ {
			delete(m, ks[i])
		}

		lock.Unlock()
	}
}
