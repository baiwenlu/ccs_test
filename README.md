# ccs_test

ccs_test is a engineering challenge homework
Its main roles or concepts are shelf, order, and courier.
The shelves are used to store the finished orders. 
Since it is public to all other roles, it is used as a global variable for convenience.
Of course, we can also pass the corresponding shelf through the interface.



It depends on golang version 1.14

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

1.
docker build -t ccs:test .
docker run -it --rm -v `pwd`/data:/data  ccs:test

2. 
cd bin
./ccs

3.
cd src
go run main.go conf.go courier.go orders.go shelf.go
