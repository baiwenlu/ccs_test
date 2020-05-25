package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestNewShelf(t *testing.T) {
	resShf := NewShelf("testname", "hot", 10, 1)

	if reflect.TypeOf(resShf).String() != "main.Shelf" {
		t.Error("shelf type error", reflect.TypeOf(resShf))
	}
	if resShf.name != "testname" {
		t.Error("return name error ", resShf.name)
	}
}

func TestDiscardOrderRand(t *testing.T) {
	testShf := NewShelf("testname", "hot", 1, 1)
	testGo := new(GridOrder)
	testGo.InTime = time.Now().Unix()
	testGo.Shf = &testShf
	testGo.Odr = &Order{Id: "4cc9d503-4e0e-42a3-b200-87c785468df9",
		Name:      "Coke",
		Temp:      "cold",
		ShelfLife: 240,
		DecayRate: 0.25,
	}
	testShf.capTable[testGo.Odr.Id] = testGo
	testShf.DiscardOrderRand()
	fmt.Println(len(testShf.capTable))
	if len(testShf.capTable) != 0 {
		t.Error("discard order err")
	}
}

func TestShelf(t *testing.T) {
	testShf := NewShelf("testname", "hot", 1, 1)
	testGo := new(GridOrder)
	testGo.InTime = time.Now().Unix()
	testGo.Shf = &testShf
	testGo.Odr = &Order{Id: "4cc9d503-4e0e-42a3-b200-87c785468df9",
		Name:      "Coke",
		Temp:      "cold",
		ShelfLife: 240,
		DecayRate: 0.25,
	}
	if testShf.GetName() != "testname" {
		t.Error("init shelf name error. Name should be testname")
	}
	if testShf.GetType() != "hot" {
		t.Error("init shelf type  error. type should be hot")
	}
	if testShf.GetCapacity() != 1 {
		t.Error("init shelf capacity error ")
	}
	if err := testShf.PutInOrder(testGo); err != nil {
		t.Error("PutInOrder error.", err)
	}
	if testShf.GetUsedCount() != 1 {
		t.Error("used count error.correct is 1,return ", testShf.GetUsedCount())
	}
	leftValue := testShf.OrderLife(*testGo)
	if leftValue > 1.0 || leftValue < 0 {
		t.Error("count left value error", leftValue)
	}

	if !testShf.IsFull() {
		t.Error("is full return false.should be true")
	}

	toShlf := NewShelf("target shelf", "overflow", 1, 1)
	if !testShf.migrateOrder(&toShlf, testGo) {
		t.Error("migrate order error")
	}

	if _, err := toShlf.PickUpOrder(testGo.Odr.Id); err != nil {
		t.Error("pick up order error.error ", err)
	}

	testShf.PutInOrder(testGo)
	testShf.DiscardOrder(testGo.Odr.Id)
	if _, ok := testShf.capTable[testGo.Odr.Id]; ok {
		t.Error("discard order grid order error")
	}
}
