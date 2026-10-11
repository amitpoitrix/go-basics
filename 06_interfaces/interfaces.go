package main

import "fmt"

/*
In Go,
An interface in Go is a type that defines a set of method signatures
Interfaces are same as abstract class/interface in Java
Here in Go, Interface is being implemented by struct but internally(i.e., we don't need to mention it
explicitly)
So whatever method which is being declared in interface needs to be implemented by another function
or methods of type struct
all they have to do is to maintain the same method signature.
Interface in Go promotes Open/Closed Principle and Dependency Injection
*/

/* creating interface which is nothing but a contract */
/* suffix interface name with "er" */
type Paymenter interface {
	/* so whoever struct methods is having same signature as this is implicitly implementing this 
	interface */
	pay(amount float32)
	refund(amount float32, account string)
}

/* Lets assume we've to make payment using payment gateway */
type Payment struct {
	paymentGateway Paymenter
}

func (p Payment) makePayment(amount float32) {
	/* Now creating instance of payment gateway & making the payment */
	// pgInstance := stripe{}
	// pgInstance.pay(amount)

	
	/* Now instead of creating stripe instance here we'll make it as field in Payment struct */
	p.paymentGateway.pay(amount)

	p.paymentGateway.refund(amount, "account1")
}


/* So in order to make payment we'll make use different payment gateway APIs */
/* Payment Gateway 1: Stripe */
type Stripe struct {}

func (s Stripe) pay(amount float32) {
	fmt.Println("Making payment using stripe:", amount)
}

/* Payment Gateway 2: RazorPay */
type RazorPay struct {}

func (r RazorPay) pay(amount float32) {
	fmt.Println("Making payment using razorpay:", amount)
}

func (r RazorPay) refund(amount float32, account string) {
	fmt.Println("Refund for account:", account, "of amount:", amount)
}

/* Payment Gateway 3: PayPal */
type PayPal struct {}

func (p PayPal) pay(amount float32) {
	fmt.Println("Making payment using paypal: ", amount)
}


func main() {
	/* creating payment struct instance to makePayment */
	// pgInstance := stripe{}
	pgInstance := RazorPay{}
	// pgInstance := PayPal{}

	paymentInstance := Payment{
		paymentGateway: pgInstance,
	}

	paymentInstance.makePayment(34.5)
}