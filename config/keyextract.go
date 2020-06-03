package config

import (
	//"github.com/BASChain/go-bas/Account"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"sync"
)

var key *keystore.Key
var lock sync.Mutex

func RestoreKey() {

	//if key != nil {
	//	return
	//}
	//
	//lock.Lock()
	//defer lock.Unlock()
	//
	//if key != nil {
	//	return
	//}
	//
	//data, err := Asset("../go-bas/key/UTC--2020-03-11T06-56-52.423772000Z--33324a5ee0b35f17536ceda27274e88e76640f24")
	//if err != nil {
	//	panic("load key failure")
	//}
	////keys := Account.PrivateKeyRecoverByBytes(data, "secret")
	//
	//if len(keys) > 0 {
	//	key = keys[0]
	//} else {
	//	panic("load key error")
	//}

}

func GetLoanKey() *keystore.Key {
	RestoreKey()
	return key
}
