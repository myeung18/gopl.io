package main

import (
	"fmt"
	"os"
)

func main() {

	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(done)
	}()

	close(done)
	done = nil
	cancelled()
}

var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		fmt.Println("done .... ")
		return true
	default:
		fmt.Println("done default .... ") //nil or
		return false
	}
}

var sema = make(chan struct{})

func bc() {

	select {
	case sema <- struct{}{}:
	case <-done:
		return
	}

}
