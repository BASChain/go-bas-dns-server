package fastsearch

import (
	"github.com/pkg/errors"
	"fmt"
	"strconv"
)

var ByteIndex [128]int
var gFastSearchMem []byte
var gFastSearchPTR [][]byte


const MaxChars int = 6
const MaxBitsLen int = 40

func init()  {
	lz:="abcdefghijklmnopqrstuvwxyz0123456789-_$"

	for i:=0;i<len(lz);i++{
		ByteIndex[lz[i]] = i
	}

	BLz:="ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for i:=0;i<len(BLz);i++{
		ByteIndex[BLz[i]] = i
	}

	sum := 0
	for i:=1;i<=MaxChars;i++{
		sum += GetNextBytesCnt(i)
	}

	fmt.Println("Max Memory: "+strconv.Itoa(sum/1024/1024)+"M")

	gFastSearchMem = make([]byte,sum)


	gFastSearchPTR = make([][]byte,MaxChars)

	curlen:=0
	startpos:=0

	for i:=0;i<MaxChars; i++{

		startpos += curlen

		gFastSearchPTR[i] = gFastSearchMem[startpos/8:]

		if curlen == 0 {
			curlen = MaxBitsLen
		}else {
			curlen *= MaxBitsLen
		}
	}
}



func GetNextBytesCnt(round int) int {
	root:=MaxBitsLen
	sum:=1

	for  i:=0;i<round;i++{
		sum *= root
	}

	return sum / 8
}

func GetNextBitsCnt(round int) int {
	root:=MaxBitsLen
	sum:=1

	if round < 1{
		return 0
	}

	for  i:=0;i<round;i++{
		sum *= root
	}

	return sum
}

func tag(i int, idx int)  {
	bytecnt:=idx>>3
	bits:=idx-(bytecnt<<3)

	tags := gFastSearchPTR[i]

	a := tags[bytecnt]
	a |= byte(1 << byte(bits))

	tags[bytecnt] = a
}

func Insert(domain string) error {
	if len(domain) > MaxChars-1{
		return errors.New("Domain name length too large")
	}

	domain = domain + "$"

	preidx := 0

	for i:=0;i<len(domain);i++{
		idx:=ByteIndex[domain[i]]

		preidx = preidx*MaxBitsLen + idx

		tag(i,preidx)
	}

	return nil

}

func InsertNonEnd(domain string) error {
	if len(domain) > MaxChars{
		return errors.New("Domain name length too large")
	}

	preidx := 0

	for i:=0;i<len(domain);i++{
		idx:=ByteIndex[domain[i]]

		preidx = preidx*MaxBitsLen + idx

		tag(i,preidx)
	}

	return nil

}


func InsertBytes(bts []byte) error {
	if len(bts) > MaxChars-1{
		return errors.New("Domain name length too large")
	}

	bts = append(bts,byte(38))

	preidx := 0

	for i:=0;i<len(bts);i++{
		idx:=int(bts[i])

		preidx = preidx*MaxBitsLen + idx

		tag(i,preidx)
	}

	return nil
}


func istag(i int, idx int) bool {
	bytecnt:=idx>>3
	bits:=idx-(bytecnt<<3)

	tags := gFastSearchPTR[i]

	a := tags[bytecnt]
	b := byte(1 << byte(bits))

	if (a&b) == b{
		return true
	}
	return false
}

func Find(domain string) bool{
	if len(domain) > MaxChars-1{
		return false
	}

	domain = domain + "$"

	preidx:=0
	var i int
	for ;i<len(domain);i++{
		idx:=ByteIndex[domain[i]]
		preidx = preidx*MaxBitsLen + idx

		if !istag(i,preidx){
			break
		}

	}

	if i==len(domain){
		return true
	}

	return false
}

func FindNonEnd(domain string) bool  {
	if len(domain) > MaxChars{
		return false
	}

	preidx:=0
	var i int
	for ;i<len(domain);i++{
		idx:=ByteIndex[domain[i]]
		preidx = preidx*MaxBitsLen + idx

		if !istag(i,preidx){
			break
		}

	}

	if i==len(domain){
		return true
	}

	return false

}


