package mem

import (
	"bytes"
	"errors"
	"github.com/BASChain/go-bas/Bas_Ethereum"
	"github.com/BASChain/go-bas/DataSync"
	"github.com/BASChain/go-bas/Miner"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kprc/nbsnetwork/common/list"
	"math/big"
	"sync"
)

type ProfitItem struct {
	receiptHash *Bas_Ethereum.Hash
	domainHash  *Bas_Ethereum.Hash
	domainOwner *common.Address
	from        *common.Address
	Amount      *big.Int
	Allocation  *[4]big.Int
	BlockNumber uint64
	TxIndex     uint
	IsDraw      int //0 is not recongnize, 1 draw,2 wait for draw
	srcType     string
}

type WithdrawDetail struct {
	BlockNumber uint64
	TxIndex     uint
	Amount      *big.Int
}

type ProfitBase struct {
	lock              sync.Mutex
	addr              *common.Address
	lWithdrawDetails  list.List
	lProfitItem       list.List
	lProfitItem4Miner list.List
	totalWithdraw     big.Int
	total4Withdraw        big.Int
	totalWithDrawTimes    int
	totalReceipts         int
	totalFromMiner    big.Int
	totalFromOwner    big.Int
	totalWithdrawFromMiner big.Int
	totalWithdrawFromOwner big.Int
}

type ProfitOwner struct {
	ProfitBase
}

type ProfitMiner struct {
	ProfitBase
}

type MinerProfitStore struct {
	lock  sync.Mutex
	store map[common.Address]*ProfitMiner
}

var (
	minerProfitStore *MinerProfitStore
	mpsLock          sync.Mutex

	//ownerProfitStore *OwnerProfitStore
	//opsLock          sync.Mutex
)

type StoreInterface interface {
	InsertReceipt(receipt *Miner.SimplifiedReceipt) error
	InsertWithdraw(withdraw *Miner.SimplifiedWithdraw) error
}

func profitItemCmp(v1, v2 interface{}) int {
	p1, p2 := v1.(*ProfitItem), v2.(*ProfitItem)

	if bytes.Compare(p1.receiptHash[:], p2.receiptHash[:]) == 0 {
		return 0
	}

	return 1

}

func withdrawCmp(v1, v2 interface{}) int {
	w1, w2 := v1.(*WithdrawDetail), v2.(*WithdrawDetail)

	if w1.BlockNumber == w2.BlockNumber && w1.TxIndex == w2.TxIndex {
		return 0
	} else {
		return 1
	}
}

func profitItemSort(v1, v2 interface{}) int {
	p1, p2 := v1.(*ProfitItem), v2.(*ProfitItem)

	if p1.BlockNumber < p2.BlockNumber {
		return 1
	} else if p1.BlockNumber == p2.BlockNumber {
		if p1.TxIndex < p2.TxIndex {
			return 1
		}
	}
	return -1
}

func withdrawSort(v1, v2 interface{}) int {
	w1, w2 := v1.(*WithdrawDetail), v2.(*WithdrawDetail)
	if w1.BlockNumber < w2.BlockNumber {
		return 1
	} else if w1.BlockNumber == w2.BlockNumber {
		if w1.TxIndex < w2.TxIndex {
			return 1
		}
	}
	return -1
}

func NewMinerProfit(addr *common.Address) *ProfitMiner {
	pm := &ProfitMiner{}

	pm.addr = addr
	pm.lWithdrawDetails = list.NewList(withdrawCmp)
	pm.lWithdrawDetails.SetSortFunc(withdrawSort)
	pm.lProfitItem = list.NewList(profitItemCmp)
	pm.lProfitItem.SetSortFunc(profitItemSort)
	pm.lProfitItem4Miner = list.NewList(profitItemCmp)
	pm.lProfitItem4Miner.SetSortFunc(profitItemSort)

	return pm
}

func GetMinerProfitStore() StoreInterface {
	if minerProfitStore != nil {
		return minerProfitStore
	}

	mpsLock.Lock()
	defer mpsLock.Unlock()

	if minerProfitStore != nil {
		return minerProfitStore
	}

	minerProfitStore := &MinerProfitStore{}
	minerProfitStore.store = make(map[common.Address]*ProfitMiner)

	return minerProfitStore
}

func StartProfitService() {
	Miner.RegMinerReceipt(GetMinerProfitStore().InsertReceipt)
	Miner.RegMinerWithdraw(GetMinerProfitStore().InsertWithdraw)
}

//func GetOwnerProfitStore() StoreInterface {
//
//	if ownerProfitStore != nil {
//		return ownerProfitStore
//	}
//
//	opsLock.Lock()
//	defer opsLock.Unlock()
//
//	if ownerProfitStore != nil {
//		return ownerProfitStore
//	}
//
//	ownerProfitStore := &OwnerProfitStore{}
//	ownerProfitStore.store = make(map[common.Address]*ProfitOwner)
//
//	return ownerProfitStore
//}

func GetDomainRecord(receipt Bas_Ethereum.Hash) (domain *DataSync.DomainRecord, pay *DataSync.Receipt, err error) {
	DataSync.PayLock()
	v, ok := DataSync.PayRecords[receipt]
	if !ok {
		DataSync.PayUnLock()
		return nil, nil, errors.New("No Pay Receipt")
	}
	pay = v.Clone()

	DataSync.PayUnLock()

	dhash := Bas_Ethereum.GetHash(string(pay.Name))

	DataSync.MemLock()
	defer DataSync.MemUnlock()

	var dr *DataSync.DomainRecord

	dr, ok = DataSync.Records[dhash]
	if !ok {
		return nil, pay, errors.New("No Domain Record")
	}

	domain = dr.Clone()

	return
}

func GetOwner(blockNum uint64, txidx uint, domainName string) *common.Address {
	dhash := Bas_Ethereum.GetHash(domainName)

	DataSync.TLock()
	trs, ok := DataSync.TransferRecords[dhash]
	if !ok {
		DataSync.TUnLock()
		return nil
	}
	DataSync.TUnLock()

	var nearest *DataSync.TransferRecord

	trs.Lock()
	defer trs.UnLock()

	cursor := trs.GetList().ListIterator(0)
	for {
		n := cursor.Next()
		if n == nil {
			break
		}
		tr := n.(*DataSync.TransferRecord)
		if tr.BlockNumber > blockNum {
			continue
		} else if tr.BlockNumber == blockNum {
			if tr.TxIndex > txidx {
				continue
			}
		}
		if nearest == nil {
			nearest = tr
			break
		}

	}

	if nearest == nil {
		return nil
	} else {
		return &nearest.To
	}

}

func ownerIsMiner(own *common.Address, sett *Miner.Setting) bool {
	if own == nil || sett == nil {
		return false
	}

	for _, m := range sett.Miners {
		if bytes.Compare(m.Bytes(), own.Bytes()) == 0 {
			return true
		}
	}

	return false
}

func priceCalc(amount *big.Int,promo uint8)  {
	
}


func (mps *MinerProfitStore) InsertReceipt(receipt *Miner.SimplifiedReceipt) error {
	if receipt == nil {
		return errors.New("Parameter Error")
	}

	if receipt.Allocation >= Miner.AllocationMax {
		return errors.New("Parameter Error")
	}

	dr, pay, err := GetDomainRecord(receipt.ReceiptNumber)
	if err != nil {
		return err
	}

	ownerAddr := GetOwner(receipt.BlockNumber, receipt.TxIndex, string(pay.Name))
	if ownerAddr == nil {
		ownerAddr = &dr.Owner
	}

	mps.lock.Lock()

	m, ok := mps.store[*ownerAddr]
	if !ok {
		m = NewMinerProfit(ownerAddr)
		mps.store[*ownerAddr] = m
	}

	mps.lock.Unlock()

	var alc *[4]big.Int

	err, sett := Miner.SettingRecords.GetClosest(receipt.BlockNumber, receipt.TxIndex)
	if err == nil {
		rcd := sett.(Miner.Setting).Allocation[receipt.Allocation]
		alc = &rcd
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	dhash := Bas_Ethereum.GetHash(string(dr.Name))

	pi := &ProfitItem{
		receiptHash: &receipt.ReceiptNumber,
		domainHash:  &dhash,
		domainOwner: m.addr,
		Amount:      &receipt.Amount,
		Allocation:  alc,
		BlockNumber: receipt.BlockNumber,
		TxIndex:     receipt.TxIndex,
		IsDraw:      0,
		srcType:     pay.Option,
		from:        &pay.Payer,
	}
	m.totalReceipts ++

	m.lProfitItem.AddValueOrder(pi)

	if alc != nil {
		s := sett.(Miner.Setting)
		if ownerIsMiner(m.addr, &s) {
			m.lProfitItem4Miner.AddValueOrder(pi)
		}
	}

	return nil
}

func (mps *MinerProfitStore) InsertWithdraw(withdraw *Miner.SimplifiedWithdraw) error {
	if withdraw == nil {
		return errors.New("Parameter error")
	}
	mps.lock.Lock()
	m, ok := mps.store[withdraw.Drawer]
	if !ok {
		m = NewMinerProfit(&withdraw.Drawer)
		mps.store[withdraw.Drawer] = m
	}

	mps.lock.Unlock()

	m.lock.Lock()
	defer m.lock.Unlock()

	wd := &WithdrawDetail{}
	wd.TxIndex = withdraw.TxIndex
	wd.BlockNumber = withdraw.BlockNumber
	wd.Amount = &withdraw.Amount

	m.lWithdrawDetails.AddValueOrder(wd)

	m.totalWithdraw.Add(&m.totalWithdraw, wd.Amount)
	m.totalWithDrawTimes ++


	return nil
}
