package main

import (
	"fmt"
	"github.com/BASChain/go-bas-dns-server/fastsearch2"
)

func main()  {

	//testfs()

	testfs2()


}

func testfs()  {
	//for i:=0;i<100000;i++{
	//
	//	bts := make([]byte,5)
	//
	//	b1:=rand.Uint32()
	//	bts[0] = byte(b1%38)
	//	b2:=b1+rand.Uint32()
	//	bts[1] = byte(b2%38)
	//	b3:=b2+rand.Uint32()
	//	bts[2] = byte(b3%38)
	//	b4:=b3+rand.Uint32()
	//	bts[3] = byte(b4%38)
	//	b5:=b4+rand.Uint32()
	//	bts[4] = byte(b5%38)
	//}
	//
	//
	//fmt.Println()
	//
	//fmt.Println(fastsearch.Find("ABcde"))
	//
	//fmt.Println(fastsearch.Find("abdde"))
	//
	//fmt.Println(fastsearch.Find("bcdab"))
	//
	//fmt.Println(fastsearch.Find("abEcc"))
	//
	//fmt.Println(fastsearch.Find("aceaa"))
	//
	//fmt.Println(fastsearch.Find("googl"))
}

func testfs2()  {

	//fastsearch2.Insert("abcd")
	//fastsearch2.Insert("sina")
	fmt.Println(fastsearch2.Insert("googlex"))

	//for i:=0;i<100000000;i++{
	//
	//	bts := make([]byte,6)
	//
	//	b1:=rand.Uint32()
	//	bts[0] = byte(b1%38)
	//	b2:=b1+rand.Uint32()
	//	bts[1] = byte(b2%38)
	//	b3:=b2+rand.Uint32()
	//	bts[2] = byte(b3%38)
	//	b4:=b3+rand.Uint32()
	//	bts[3] = byte(b4%38)
	//	b5:=b4+rand.Uint32()
	//	bts[4] = byte(b5%38)
	//	b6:=b5+rand.Uint32()
	//	bts[5] = byte(b6%38)
	//
	//	//fmt.Println(bts)
	//	fastsearch2.InsertIdxs(bts)
	//}

	fmt.Println(fastsearch2.Find("google"))

}