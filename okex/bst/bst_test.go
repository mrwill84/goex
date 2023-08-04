package bst

import (
	"fmt"
	"testing"
)

type PriceBlock struct {
	PriceIdx int64
	Price    float64
	Amount   float64
}

func TestBST(t *testing.T) {
	bstree := BSTree{}
	pb := PriceBlock{
		PriceIdx: 1,
		Price:    1,
		Amount:   9,
	}
	pb2 := PriceBlock{
		PriceIdx: 2,
		Price:    2,
		Amount:   1,
	}
	bstree.Upsert("1", &pb)
	bstree.Upsert("2", &pb2)
	bstree.Upsert("3", &pb2)
	bstree.Upsert("4", &pb2)
	for i := range bstree.RIter() {
		fmt.Println(i)
	}
	fmt.Println("/////")
	fmt.Println(bstree.Value("1"))
	fmt.Println("/////")
	bstree.Delete("1")
	for i := range bstree.Iter() {
		fmt.Println(i)
	}
	fmt.Println("/////")
	bstree.Delete("2")
	for i := range bstree.Iter() {
		fmt.Println(i)
	}
	//fmt.Println(PreOrder(&bstree))
	//a := bstree.Search(1)
	//fmt.Println(a)
	//bstree.Delete(2)

}

//go test -v -timeout 30s -run ^TestBST$ github.com/mrwill84/goex/okex/bst
