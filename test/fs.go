package main

import (
	//"fmt"
	//"github.com/BASChain/go-bas-dns-server/fastsearch2"
	"time"
	//"github.com/BASChain/go-bas-dns-server/fastsearch2"
	"fmt"
	"github.com/kprc/nbsnetwork/common/list"
	"github.com/zserge/webview"
)

func main()  {

	//testfs()

	//bgi:=big.Int{}
	//
	//fmt.Println("===>",bgi.String())
	//
	////testfs2()
	//testmyTimeAfter()
	//testTimeAfter()
	//testlist()

	//time.After(time.Second)
	//time.AfterFunc()

		//debug := true
		w := webview.New(false)
		defer w.Destroy()
		w.SetTitle("Minimal webview example")
		w.SetSize(1024, 768, webview.HintNone)
		w.Navigate("https://www.baidu.com")
		w.Run()


}

func getNextTime(d int) <-chan time.Time {
	fmt.Println("start getNextTime",d)
	c := make(chan time.Time)

	time.Sleep(time.Second*time.Duration(d))

	c<-time.Now()

	time.Sleep(time.Second*time.Duration(1))

	fmt.Println("end getNextTime",d)
	return c
}

//func timeAfter() <-chan time.Time {
//	//fmt.Println("Start time after")
//
//
//
//
//
//	//fmt.Println("End time after")
//
//}

func testmyTimeAfter()  {
	for{
		select {
		case <-getNextTime(2):
			fmt.Println("AAAAAA")
		case <-getNextTime(3):
			fmt.Println("BBBBBB")
		}
	}
}

func testTimeAfter()  {
	for{
		select {
		case <-time.After(time.Second*2):
			fmt.Println("AAAAAA")
		case <-time.After(time.Second*2):
			fmt.Println("BBBBBB")
		}
	}
}

type a struct {
	x int
}

func testacmp(v1,v2 interface{}) int  {
	a1,a2 := v1.(*a),v2.(*a)


	if a1.x == a2.x{
		return  0
	}

	return 1
}

func testasort(v1,v2 interface{}) int {
	a1,a2 := v1.(*a),v2.(*a)

	if a1.x > a2.x{
		return 1
	}

	return -1

}

func testlist()  {
	l := list.NewList(testacmp)
	l.SetSortFunc(testasort)

	a1:=&a{3}
	a2:=&a{2}
	a3:=&a{4}

	l.AddValueOrder(a1)
	l.AddValueOrder(a2)
	l.AddValueOrder(a3)

	cusor:=l.ListIterator(0)

	for{
		n:=cusor.Next()
		if n == nil{
			break
		}
		fmt.Println(n.(*a).x)
	}


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
	//fmt.Println(fastsearch2.Insert("googlex"))
	//fmt.Println(fastsearch2.Insert("aaaaaa"))
	//fmt.Println(fastsearch2.Insert("bbbbbb"))
	//fmt.Println(fastsearch2.Insert("cccccc"))
	//fmt.Println(fastsearch2.Insert("______"))

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

	//fmt.Println(fastsearch2.Find("google"))
	//fmt.Println(fastsearch2.Find("aaaaaa"))
	//fmt.Println(fastsearch2.Find("______"))



}