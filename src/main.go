package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
)

//four shelfs
var (
	HotShelf      Shelf
	ColdShelf     Shelf
	FrozenShelf   Shelf
	OverflowShelf Shelf
)

//it recovers from panics and prints the stacktrace
func RecoverFromPanic() {
	if err := recover(); err != nil {

		var stacktrace string
		for i := 1; ; i++ {
			_, f, l, got := runtime.Caller(i)
			if !got {
				break
			}

			stacktrace += fmt.Sprintf("%s:%d\n", f, l)
		}
		// when stack finishes
		logMessage := fmt.Sprintf("Trace: %s\n", err)
		logMessage += fmt.Sprintf("\n%s", stacktrace)
		log.Fatal(logMessage)
	}
}

func main() {
	defer RecoverFromPanic()
	//Initialize the shelfs
	HotShelf = NewShelf("hot shelf", "hot", SingleTemperatureCap, SingleTemperatureDecayModifier)
	ColdShelf = NewShelf("cold shelf", "cold", SingleTemperatureCap, SingleTemperatureDecayModifier)
	FrozenShelf = NewShelf("frozen shelf", "frozen", SingleTemperatureCap, SingleTemperatureDecayModifier)
	OverflowShelf = NewShelf("overflow shelf", "overflow", AnyTemperatureCap, AnyTemperatureDecayModifier)
	//new order's channel used to synchronize consumers and cookers
	chanOrders := make(chan Order, 100)
	//new grid order's channel used to synchronize cookers and couriers
	chanGridOrders := make(chan *GridOrder, 100)

	var wg sync.WaitGroup // synchronize goroutines
	wg.Add(2)
	go CreateOrder(chanOrders, &wg)                     //produce new order goroutine
	go CookAndDispatch(chanOrders, chanGridOrders, &wg) // cook and dispatch order goroutine
	InitCouriers(chanGridOrders, &wg)                   //Initialize the couriers
	wg.Wait()
	log.Print("all orders have been picked up")
}
