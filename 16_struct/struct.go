package main

import (
	"fmt"
	"time"
)

/*
In Go
We make use of struct to create custom Data Structure
we don't have classes so we make use of struct to fullfil class requirements
*/

// As we don't have inheritance in Go so we create another struct and mention it in parent struct
type customer struct {
	name  string
	phone string
}

type order struct {
	id        string
	amount    float32
	status    string
	createdAt time.Time // nanosecond precision
	customer            // used to create composition or relation between two struct
}

/*
Like in other OOPS languages we've class properties along with there behaviour in the form of methods
Here in Go, we can create custom function and attch it to struct that will act as behaviour for that struct

(o *order) -> this is receiver type
Here we're changing the field value of struct than we've to use pointer i.e., o *order
As per convention use first letter of struct to name receiver type variable
*/
func (o *order) changeStatus(status string) {
	fmt.Println("Changing status")
	o.status = status // no need to deference(i.e., *o.status) as Go will take care of it
}

// While fetching value of any struct field no need to pass receiver type as pointer
func (o order) getAmount() float32 {
	return o.amount
}

// As we don't have constructors in Go so we create our own custom one in Go
// Convention is to use new as prefix in struct name with camelCase
func newOrder(id string, amount float32, status string) *order {
	// intial setup go here
	myInitialOrder := order{
		id:     id,
		amount: amount,
		status: status,
	}

	// As return type of this constructor is pointer of order struct, so we need to pass address
	return &myInitialOrder
}

func main() {
	// creating struct order instance and initializing its value
	// no field is mandatory, if not defined they contain there zeroed values
	// As struct is like class in Go so we can create n no. of instance of order struct
	myOrder := order{
		id:     "1",
		amount: 50.5,
		status: "received",
	}

	// setting field value of order struct
	myOrder.createdAt = time.Now()

	fmt.Println(myOrder.status)
	fmt.Println("myOrder struct", myOrder)

	// creating another instance of order struct
	myOrder2 := order{
		id:        "2",
		amount:    23.8,
		status:    "delivered",
		createdAt: time.Now(),
	}

	fmt.Println("myOrder2 struct", myOrder2)


	// created behaviour of order struct to change its status
	myOrder2.changeStatus("paid")
	fmt.Println("myOrder2 struct", myOrder2)

	// Now in order to fetch any order field value we can have its custom method
	fmt.Println("Fetched amount from myOrder2 ", myOrder2.getAmount())


	// Cretaing order struct instance using custom constructor
	myOrder3 := newOrder("3", 45.5, "Initiated")

	fmt.Println("myOrder3 using custom constructor: ", myOrder3)

	// Creating struct and initializing at the same time
	language := struct {
		name   string
		isGood bool
	}{"golang", true}

	fmt.Println("One liner struct: ", language)

	fmt.Println("Compisiton in Go")
	myCustomer := customer{
		name:  "Amit",
		phone: "1234567890",
	}

	
	myOrder4 := order{
		id:       "4",
		amount:   12.00,
		status:   "new",
		customer: myCustomer,
	}

	// setting value
	myOrder4.customer.name = "sorr"

	fmt.Println("Struct using composition: ", myOrder4)
}
