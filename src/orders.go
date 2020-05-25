package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

type Order struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Temp      string  `json:"temp"` //The temperature at which the order should be stored
	ShelfLife float64 `json:"shelfLife"`
	DecayRate float64 `json:"decayRate"`
}

//load orders from source data file
func LoadOrders(file string, orders *[]Order) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(content), orders)
	if err != nil {
		panic(err)
	}
}

//CreateOrder imitates customers placing orders
//it put order into channel order
func CreateOrder(ch chan Order, wg *sync.WaitGroup) {
	var orders []Order
	LoadOrders(OrderFile, &orders)
	for _, order := range orders {
		time.Sleep(1e9 / CreateOrderRate)
		ch <- order
	}
	close(ch)
	wg.Done()
}

//CookAndDispatch imitates cookers cooking orders and
// puts orders into channel chanGridOrder to waiting for courier to pick up
func CookAndDispatch(ch chan Order, chanGridOrders chan *GridOrder, wg *sync.WaitGroup) {
	for order := range ch {
		DispatchOrder(order, chanGridOrders)
		//fmt.Printf("dispatch order %s to shelf %s,order temperature is: %s \n", order.Name, shType, order.Temp)
		// fmt.Println(HotShelf, ColdShelf, FrozenShelf, OverflowShelf)

	}
	close(chanGridOrders)
	wg.Done()
}

//DispatchOrder puts orders on the corresponding shelves
//and into channel chanGridOrder to waiting for courier to pick up
func DispatchOrder(ordr Order, chanGridOrders chan *GridOrder) {
	GO := new(GridOrder)
	GO.Odr = &ordr

	switch {
	case ordr.Temp == "hot" && !HotShelf.IsFull():
		GO.Shf = &HotShelf
	case ordr.Temp == "cold" && !ColdShelf.IsFull():
		GO.Shf = &ColdShelf
	case ordr.Temp == "frozen" && !FrozenShelf.IsFull():
		GO.Shf = &FrozenShelf
	case !OverflowShelf.IsFull():
		GO.Shf = &OverflowShelf
	default:
		if !OverflowShelf.MigrateOrWastOrder() {
			OverflowShelf.DiscardOrderRand()
		}
		GO.Shf = &OverflowShelf
	}
	GO.InTime = time.Now().Unix()
	if err := GO.Shf.PutInOrder(GO); err != nil {
		log.Fatal("put in order fail")
		return
	}
	chanGridOrders <- GO

	log.Printf("dispatch order %s to shelf %s,order temperature is: %s. Order value is %f", ordr.Id, GO.Shf.GetType(), ordr.Temp, GO.Shf.OrderLife(*GO))
}
