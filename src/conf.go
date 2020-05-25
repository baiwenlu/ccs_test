package main

const (
	CreateOrderRate                = 2                     //order's created count per second
	CourierCount                   = 5                     //total couriers
	SingleTemperatureCap           = 10                    //single temperature shelf's capacity
	AnyTemperatureCap              = 15                    //overflow shelf's capacity
	SingleTemperatureDecayModifier = 1                     // single temperature shelf's decay rate
	AnyTemperatureDecayModifier    = 2                     // overflow shelf's decay rate
	CourierMaxTime                 = 6                     // courier‘s max seconds spend on the road
	CourierMinTime                 = 2                     //courier‘s min seconds spend on the road
	OrderFile                      = "../data/orders.json" // source orders's data file
)
