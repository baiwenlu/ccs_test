package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

// courier is responsible to pick up order and deliver order
type Courier struct {
	Num int
}

func NewCourier(num int) *Courier {
	return &Courier{Num: num}
}

//Initialize the courier according to the configuration CourierCount
func InitCouriers(chanGridOrders chan *GridOrder, wg *sync.WaitGroup) {
	for i := 1; i <= CourierCount; i++ {
		wg.Add(1)
		go func(i int, ch chan *GridOrder, wg *sync.WaitGroup) {
			courier := NewCourier(i)
			for gOrder := range ch {
				courier.SpendTime()
				courier.PickUp(gOrder)
			}
			wg.Done()
		}(i, chanGridOrders, wg)
	}
}

//pick up the order
func (c *Courier) PickUp(gOrder *GridOrder) {

	if order, err := gOrder.Shf.PickUpOrder(gOrder.Odr.Id); err == nil {
		log.Printf("courier %d picked up order %s from %s. Order value is %f", c.Num, order.Id, gOrder.Shf.GetName(), gOrder.Shf.OrderLife(*gOrder))
	} else {
		log.Printf("order %s has been discarded by shelf %s.", gOrder.Odr.Id, gOrder.Shf.GetName())
	}
	// gOrder.Shf.ShowAllOrders()
}

//spend time on the road
func (c *Courier) SpendTime() {
	rand.Seed(time.Now().UnixNano())
	rndSec := rand.Int63n(CourierMaxTime-CourierMinTime+1) + CourierMinTime
	time.Sleep(time.Duration(rndSec * 1e9))
}
