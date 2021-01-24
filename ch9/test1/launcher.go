package main

import "sync"

func main() {

}

var (
	sema    = make(chan struct{}, 1)
	balance int

	mu sync.Mutex
)

func Deposit(amt int) {
	mu.Lock()
	balance = balance + amt
	mu.Unlock()
}

func Balance() int {

	mu.Lock()
	b := balance
	mu.Unlock()

	return b
}
