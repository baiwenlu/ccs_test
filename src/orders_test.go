package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestLoadOrders(t *testing.T) {
	var orders []Order
	LoadOrders("../data/orders.json", &orders)
	if len(orders) == 0 {
		t.Error("load order fail")
	}
}

func TestCreateOrder(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	chanTest := make(chan Order)
	go CreateOrder(chanTest, &wg)
	order := <-chanTest
	fmt.Println(order)
	// wg.Wait()
	if order.Id == "" {
		t.Error("create order error")
	}
}
