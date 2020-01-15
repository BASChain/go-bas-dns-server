package fastsearch2

import (
	"strconv"
	"fmt"
)

var ByteIndex [128]int


type BitPtr struct {
	Level int
	Start int
	Length int
}


const MaxChars int = 6
var MaxBitsLen int

type FastSearchBitPtr struct {
	fastSearchBitPtr []*BitPtr
	startMemory []byte
}

var gFastBitPtr FastSearchBitPtr
var gFastSearchEndPoint FastSearchBitPtr


func init()  {

	lz := "abcdefghijklmnopqrstuvwxyz0123456789-_"

	for i:=0;i<len(lz);i++{
		ByteIndex[lz[i]] = i
	}

	BLz := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for i:=0;i<len(BLz);i++{
		ByteIndex[BLz[i]] = i
	}

	MaxBitsLen = len(lz)

	length := 0
	roundbits := 1

	i:=0

	for ;i<MaxChars;i++{
		roundbits *= MaxBitsLen
		bp := &BitPtr{}
		bp.Level = i
		bp.Start = length
		bp.Length = roundbits

		length += roundbits
		gFastBitPtr.fastSearchBitPtr = append(gFastBitPtr.fastSearchBitPtr,bp)
	}

	for j:=0;j<len(gFastBitPtr.fastSearchBitPtr);j++{
		bp:=gFastBitPtr.fastSearchBitPtr[j]
		gFastSearchEndPoint.fastSearchBitPtr = append(gFastSearchEndPoint.fastSearchBitPtr,bp.Clone())
	}

	memlen := length >> 3

	if length - (memlen << 3) > 0{
		memlen ++
	}

	gFastBitPtr.startMemory = make([]byte,memlen)
	gFastSearchEndPoint.startMemory = make([]byte,memlen)

	fmt.Println("total memory :",strconv.Itoa((2*memlen)/1024/1024),"M","/ max domain name:",strconv.Itoa(length))
}

func (bp *BitPtr)Clone() *BitPtr {
	bp1:=&BitPtr{}

	bp1.Length = bp.Length
	bp1.Start = bp.Start
	bp1.Level = bp.Level

	return bp1
}

func tag(idx,bitpos int, isEnd bool)  {
	bp:=gFastBitPtr.fastSearchBitPtr[idx]

	pos := bp.Start + bitpos

	bytescnt := pos >> 3

	bitscnt := pos - (bytescnt << 3)

	a := gFastBitPtr.startMemory[bytescnt]
	a |= byte(1 << byte(bitscnt))

	gFastBitPtr.startMemory[bytescnt] = a


	if isEnd{
		b := gFastSearchEndPoint.startMemory[bytescnt]
		b |= byte(1 << byte(bitscnt))
		gFastSearchEndPoint.startMemory[bytescnt] = b
	}
}


func Insert(domain string) (isEnd bool,predix int)  {

	var inslen int
	if len(domain) > MaxChars{
		inslen = MaxChars
	}else{
		isEnd = true
		inslen = len(domain)
	}

	predix = 0

	for i:=0;i<inslen;i++{
		bitpos:=ByteIndex[domain[i]]

		predix = predix * MaxBitsLen + bitpos

		if i == (inslen - 1){
			tag(i,predix,isEnd)
		}else{
			tag(i,predix,false)
		}
	}

	return
}

func InsertIdxs(idxs []byte) (isEnd bool,predix int) {
	var inslen int
	if len(idxs) > MaxChars{
		inslen = MaxChars
	}else{
		isEnd = true
		inslen = len(idxs)
	}

	predix = 0
	for i:=0;i<inslen;i++{
		bitpos:=int(idxs[i])

		predix = predix * MaxBitsLen + bitpos

		if i == (inslen - 1){
			tag(i,predix,isEnd)
		}else{
			tag(i,predix,false)
		}
	}

	return
}


func istag(idx,bitpos int, isEnd bool) (ok bool) {
	bp:=gFastBitPtr.fastSearchBitPtr[idx]

	pos := bp.Start + bitpos

	bytescnt := pos >> 3

	bitscnt := pos - (bytescnt << 3)

	a := gFastBitPtr.startMemory[bytescnt]
	b := byte(1 << byte(bitscnt))


	if (a & b) == b {
		ok = true
	}

	if isEnd{
		c := gFastSearchEndPoint.startMemory[bytescnt]
		if (c & b) == b{
			ok = true
		}else{
			ok = false
		}
	}

	return
}

func Find(domain string)  (isEnd bool,isFind bool,predix int) {

	var inslen int
	if len(domain) > MaxChars{
		inslen = MaxChars
	}else{
		isEnd = true
		inslen = len(domain)
	}

	predix = 0

	i:=0

	for ;i<inslen;i++{
		bitpos:=ByteIndex[domain[i]]

		predix = predix * MaxBitsLen + bitpos

		if i == (inslen - 1){
			if !istag(i,predix,isEnd){
				break
			}
		}else{
			if !istag(i,predix,false){
				break
			}
		}
	}
	if i == inslen {
		isFind = true
		predix += gFastSearchEndPoint.fastSearchBitPtr[i-1].Start
	}else{
		predix += gFastSearchEndPoint.fastSearchBitPtr[i].Start
	}

	return
}


