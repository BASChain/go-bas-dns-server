package mem

import (
	"sync"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/kprc/nbsnetwork/tools"
)

type MemState struct {
	FreeState [2]int //1 Waiting result, 2 success, 3 failure
	WriteTime int64
}

var (
	lock sync.Mutex
	m map[common.Address]*MemState
)

const(
	BAS int = 0
	ETH int = 1

	WAITING int = 1
	SUCCESS int = 2
	FAILURE int = 3
)

func Update(addr common.Address,typ int,state int) error {
	if typ != BAS && typ != ETH{
		return errors.New("type error")
	}
	if state != WAITING && state != SUCCESS && state != FAILURE{
		return errors.New("state error")
	}
	lock.Lock()
	defer lock.Unlock()

	if s,ok:=m[addr];!ok{
		s = &MemState{}
		s.FreeState[typ] = state
		s.WriteTime = tools.GetNowMsTime()
		m[addr] = s
	}else{
		s.FreeState[typ] = state
		s.WriteTime = tools.GetNowMsTime()
	}

	return nil

}

func GetState(addr common.Address,typ int) (int,error) {
	lock.Lock()
	defer lock.Unlock()

	if typ != BAS && typ != ETH{
		return 0,errors.New("type error")
	}

	if s,ok:=m[addr];!ok{
		return 0,errors.New("No Address")
	}else{
		return s.FreeState[typ],nil
	}
}