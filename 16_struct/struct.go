package main

import (
	"fmt"
	"time"
)

/*
In Go
A struct in Go is a user-defined data type that groups multiple related fields into a single unit
We make use of struct to create custom Data Structure
we don't have classes so we make use of struct to fullfil class requirements
*/

// Method 1: As we don't have inheritance in Go so we create another struct and mention it in parent struct
type Customer struct { // struct name should also start from CAPITAL Letters
	Name  string // Always start struct field name with CAPITAL Letters
	Phone string
}

type Order struct {
	ID        string
	Amount    float32
	Status    string
	CreatedAt int64 // nanosecond precision
	Customer        // used to create composition or relation between two struct
}

/*
Like in other OOPS languages we've class properties along with there behaviour in the form of methods
Here in Go, we can create custom function(or methods) and attach it to struct that will act as
struct behaviour.

In Struct we've two types of methods(which is written outside main function)
1. Value receiver
2. Pointer receiver
*/

// 1. Value Receiver: for fetching struct fields
// Note: No need to pass receiver type as pointer
func (o Order) getAmount() float32 {
	fmt.Println("Fetching Order Amount value...")
	return o.Amount
}

/*
2. Pointer Receiver - for updating struct fields
(o *Order) -> this is receiver type
Here we're changing the field value of struct than we've to use pointer receiver i.e., o *Order
As per convention use first letter of struct to name receiver type variable
*/
func (o *Order) updateStatus(status string) {
	fmt.Println("Updating status...")
	o.Status = status // no need to deference(i.e., *o.status) as Go will take care of it
}

// As we don't have constructors in Go so we create our own custom one in Go
// Convention is to use new as prefix in struct name with camelCase
func newOrder(id string, amount float32, status string) *Order {
	// intial setup go here
	myInitialOrder := Order{
		ID:     id,
		Amount: amount,
		Status: status,
	}

	// As return type of this constructor is pointer of Order struct, so we need to pass address
	return &myInitialOrder
}

func main() {
	// Method 1.1: Named Initialization
	fmt.Println("Method 1.1: Named Initialization")
	// creating struct order instance and initializing its value
	// no field is mandatory, if not defined they contain there zeroed values
	// As struct is like class in Go so we can create n no. of instance of order struct

	myOrder := Order{
		ID:     "1",
		Amount: 50.5,
		Status: "received",
	}

	// setting field value of order struct
	myOrder.CreatedAt = time.Now().Unix()

	fmt.Println("myOrder.Status:", myOrder.Status)
	fmt.Println("myOrder struct:", myOrder)

	// Method 1.2: Positional Initialization
	fmt.Println("Method 1.2: Positional Initialization")

	myOrder2 := Order{
		"2",
		23.8,
		"delivered",
		time.Now().Unix(),
		Customer{
			"Amit",
			"123456",
		},
	}

	fmt.Println("myOrder2 struct:", myOrder2)

	// Method 1.3: Using new keyword which returns pointer to the newly allocated memory to the
	// user defined struct
	fmt.Println("Method 1.3: Using new keyword")

	myOrder31 := new(Order)
	myOrder31.ID = "7"
	myOrder31.Amount = 56.123

	fmt.Println("myOrder31 struct returns address: ", myOrder31)

	// Method 1.4: Using Pointer - this also returns pointer
	fmt.Println("Method 1.4: Using Pointer")

	myOrder32 := &Order{
		ID:     "3",
		Amount: 78.09,
	}

	fmt.Println("myOrder32 struct also returns address:", myOrder32)

	// Updating the status of Order
	fmt.Println("Updating the status of Order using Pointer Receiver method")
	myOrder2.updateStatus("Paid")
	fmt.Println("After Updating Order Status:", myOrder2)

	// Now in order to fetch any order field value we can have its custom method
	fmt.Println("Fetching Order field value using Value Receiver method")
	fmt.Println("Fetched amount from myOrder2:", myOrder2.getAmount())

	// Cretaing order struct instance using custom constructor
	myOrder3 := newOrder("3", 45.5, "Initiated")

	fmt.Println("myOrder3 using custom constructor: ", myOrder3)

	// Method 2: Anonymous Struct - Creating struct and initializing at the same time
	language := struct {
		Name   string
		IsGood bool
	}{
		"golang",
		true,
	}

	fmt.Println("Anonymous struct:", language)

	fmt.Println("Compisiton in Go")
	myCustomer := Customer{
		Name:  "Amit",
		Phone: "1234567890",
	}

	myOrder4 := Order{
		ID:       "4",
		Amount:   12.00,
		Status:   "new",
		Customer: myCustomer,
	}

	// setting value
	myOrder4.Customer.Name = "sorr"

	fmt.Println("Struct using composition: ", myOrder4)
}
