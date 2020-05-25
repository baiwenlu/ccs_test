package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Shelfer interface {
	GetName() string
	GetUsedCount() int
	GetCapacity() int
	GetType() string
	ShowOrderLife(ordr Order) float64
	DiscardOrderRand()
	PutInOrder(ordr *Order) error
}

type GridOrder struct {
	Odr    *Order
	Shf    *Shelf
	InTime int64
}

//Shelf is responsible for storing orders
//CapTable is the collection of grids to store order
type Shelf struct {
	name          string
	capacity      int    // shelf's capacity
	sfType        string // Suitable temperature for the orders
	capTable      map[string]*GridOrder
	decayModifier int
	mu            sync.RWMutex
}

//NewShelf initialize a shelf and return it
func NewShelf(name, stType string, capacity int, decayModifier int) Shelf {
	//gm := newGridMap()
	var res Shelf
	res.name = name
	res.sfType = stType
	res.capacity = capacity
	res.decayModifier = decayModifier
	res.capTable = make(map[string]*GridOrder)
	return res
}

//GetName returns shelf name
func (sh *Shelf) GetName() string {
	return sh.name
}

//GetType returns shelf type
func (sh *Shelf) GetType() string {
	return sh.sfType
}

//GetUsedCount returns current len of capTable
func (sh *Shelf) GetUsedCount() int {
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	return len(sh.capTable)
}

//GetCapacity returns shelf capacity
func (sh *Shelf) GetCapacity() int {
	return sh.capacity
}

//DiscardOrderRand discard a rand order when capTable is full
func (sh *Shelf) DiscardOrderRand() {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	for len(sh.capTable) >= sh.capacity {
		for oid := range sh.capTable {
			delete(sh.capTable, oid)
			log.Printf("discard an order %s from %s.\n", oid, sh.GetName())
			sh.showAllOrders()
			break
		}
	}
}

//OrderLife returns the order's order value
//order value is order's left shelflife divided by original shelflife
func (sh *Shelf) OrderLife(ordr GridOrder) float64 {
	if ordr.Odr.ShelfLife <= 0 {
		return 0.0
	}
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	leftTime := sh.countLeftTime(ordr)

	return leftTime / float64(ordr.Odr.ShelfLife)
}

//countLeftTime returns orders's left shelflife
func (sh *Shelf) countLeftTime(ordr GridOrder) float64 {
	leftTime := float64(ordr.Odr.ShelfLife) - ordr.Odr.DecayRate*
		float64(time.Now().Unix()-ordr.InTime)*
		float64(sh.decayModifier)
	return leftTime
}

//PutInOrder puts an order on shelf
//if shelf is full ,it returns error
func (sh *Shelf) PutInOrder(ordr *GridOrder) error {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sh.capacity <= len(sh.capTable) {
		return errors.New("shelf is full")
	}

	sh.capTable[ordr.Odr.Id] = ordr
	if orderInfo, err := json.Marshal(&ordr.Odr); err == nil {
		log.Printf("%s received an order %s ", sh.GetName(), string(orderInfo))
		sh.showAllOrders()
	}
	return nil
}

//IsFull returns is or not is full of current shelf
func (sh *Shelf) IsFull() bool {
	if sh.capacity <= sh.GetUsedCount() {
		return true
	}
	return false
}

//MigrateOrWastOrder randomly migrates an order to another shelf or
//drops an order when it has no shelflife left
//return true or false
func (sh *Shelf) MigrateOrWastOrder() bool {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	var toSh Shelf
	for oid, ordr := range sh.capTable {
		if sh.countLeftTime(*ordr) <= 0.0 {
			delete(sh.capTable, oid)
			log.Printf("discard order %s because of life time zero from %s", oid, sh.GetName())
			sh.showAllOrders()
			return true
		}
		if ordr.Odr.Temp == "hot" {
			toSh = HotShelf
		} else if ordr.Odr.Temp == "cold" {
			toSh = ColdShelf
		} else if ordr.Odr.Temp == "frozen" {
			toSh = FrozenShelf
		} else {
			continue
		}
		if sh.migrateOrder(&toSh, ordr) {
			log.Printf("migrate order %s from %s to %s. Order left shelf life is %f", oid, sh.GetName(), toSh.GetName(), ordr.Odr.ShelfLife)
			sh.showAllOrders()
			return true
		} else {
			continue
		}
	}
	return false
}

//migrateOrder randomly migrates an order to another shelf
//return true or false
func (sh *Shelf) migrateOrder(toSh *Shelf, gOrder *GridOrder) bool {
	toSh.mu.Lock()
	defer toSh.mu.Unlock()
	if len(toSh.capTable) >= toSh.capacity {
		return false
	}
	gOrder.Shf = toSh
	gOrder.Odr.ShelfLife = sh.countLeftTime(*gOrder)
	gOrder.InTime = time.Now().Unix()
	toSh.capTable[gOrder.Odr.Id] = gOrder
	delete(sh.capTable, gOrder.Odr.Id)
	return true
}

//PickUpOrder picks up an order from shelf and return
//this order with error info
func (sh *Shelf) PickUpOrder(orderId string) (*Order, error) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if ordr, ok := sh.capTable[orderId]; ok {
		delete(sh.capTable, orderId)
		log.Printf("order %s was picked up", orderId)
		sh.showAllOrders()
		return ordr.Odr, nil
	} else {
		return &Order{}, errors.New("order empty")
	}
}

//DiscardOrder delete an order from shelf,
//orderId is Order's property Id
func (sh *Shelf) DiscardOrder(orderId string) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	delete(sh.capTable, orderId)
}

//ShowAllOrders prints all of orders's order value on this shelf
func (sh *Shelf) ShowAllOrders() {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	sh.showAllOrders()
}

//showAllOrders prints all of orders's order value on this shelf
func (sh *Shelf) showAllOrders() {
	message := fmt.Sprintf("\tthe contents of all orders in %s are as follows:\n", sh.GetName())
	i := 1
	for _, gOrder := range sh.capTable {
		if orderInfo, err := json.Marshal(&gOrder.Odr); err == nil {
			message += fmt.Sprintf("\t%d is %s. order value is %f\n", i, string(orderInfo), sh.countLeftTime(*gOrder)/gOrder.Odr.ShelfLife)
			i++
		}
	}
	fmt.Println(message)
}
