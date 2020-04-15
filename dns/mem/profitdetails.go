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
	Amount      *big.Int
	Allocation  *[][]big.Int
	BlockNumber uint64
	TxIndex     uint
	IsDraw      bool
	srcType     string
}

type WithdrawDetail struct {
	BlockNumber uint64
	TxIndex     uint
	Amount      *big.Int
}

type ProfitBase struct {
	lock               sync.Mutex
	addr               *common.Address
	lWithdrawDetails   list.List
	lProfitItem        list.List
	totalWithdraw      big.Int
	totalWait2Withdraw big.Int
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

type OwnerProfitStore struct {
	lock  sync.Mutex
	store map[common.Address]*ProfitOwner
}

var (
	minerProfitStore *MinerProfitStore
	mpsLock          sync.Mutex

	ownerProfitStore *OwnerProfitStore
	opsLock          sync.Mutex
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

func prifitItemSort(v1, v2 interface{}) int {
	p1, p2 := v1.(*ProfitItem), v2.(*ProfitItem)

	if p1.BlockNumber > p2.BlockNumber {
		return 1
	} else if p1.BlockNumber == p2.BlockNumber {
		if p1.TxIndex > p2.TxIndex {
			return 1
		}
	}
	return -1
}

func withdrawSort(v1, v2 interface{}) int {
	w1, w2 := v1.(*WithdrawDetail), v2.(*WithdrawDetail)
	if w1.BlockNumber > w2.BlockNumber {
		return 1
	} else if w1.BlockNumber == w2.BlockNumber {
		if w1.TxIndex > w2.TxIndex {
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
	pm.lProfitItem.SetSortFunc(prifitItemSort)

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

func GetOwnerProfitStore() StoreInterface {

	if ownerProfitStore != nil {
		return ownerProfitStore
	}

	opsLock.Lock()
	defer opsLock.Unlock()

	if ownerProfitStore != nil {
		return ownerProfitStore
	}

	ownerProfitStore := &OwnerProfitStore{}
	ownerProfitStore.store = make(map[common.Address]*ProfitOwner)

	return ownerProfitStore
}

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

	cursor:=trs.GetList().ListIterator(0)
	for{
		n:=cursor.Next()
		if n == nil{
			break
		}
		tr:=n.(*DataSync.TransferRecord)
		if tr.BlockNumber > blockNum{
			continue
		}else if tr.BlockNumber == blockNum{
			if tr.TxIndex > txidx{
				continue
			}
		}
		if nearest == nil{
			nearest = tr
			break
		}

	}

	if nearest == nil{
		return nil
	}else{
		return &nearest.To
	}

}

func (mps *MinerProfitStore) InsertReceipt(receipt *Miner.SimplifiedReceipt) error {
	if receipt == nil {
		return errors.New("Parameter Error")
	}
	dr, pay, err := GetDomainRecord(receipt.ReceiptNumber)
	if err != nil {
		return err
	}

	mps.lock.Lock()

	m, ok := mps.store[pay.Payer]
	if !ok {
		ownerAddr := GetOwner(receipt.BlockNumber, receipt.TxIndex, string(pay.Name))
		if ownerAddr == nil {
			ownerAddr = &dr.Owner
		}
		m = NewMinerProfit(ownerAddr)
		mps.store[*ownerAddr] = m
	}

	mps.lock.Unlock()
	m.lock.Lock()
	defer m.lock.Unlock()

	return nil
}

func (mps *MinerProfitStore) InsertWithdraw(withdraw *Miner.SimplifiedWithdraw) error {
	return nil
}

func (ops *OwnerProfitStore) InsertReceipt(receipt *Miner.SimplifiedReceipt) error {
	return nil
}

func (ops *OwnerProfitStore) InsertWithdraw(withdraw *Miner.SimplifiedWithdraw) error {
	return nil
}
