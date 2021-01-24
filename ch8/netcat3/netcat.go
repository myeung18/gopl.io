// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 227.

// Netcat is a simple read/write client for TCP servers.
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

//!+
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		fmt.Println("send to server")
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		fmt.Println("done")

		done <- struct{}{} 		// signal the main goroutine
	}()

	//receive from server
	mustCopy(conn, os.Stdin)
	conn.Close()   	/** this has two effects  */
	<-done 			// wait for background goroutine to finish

	fmt.Println("after done")
}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func test() {
	conn, err := net.Dial("","")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		//inside a thread/goroutine
		fmt.Println("send a msg")
		io.Copy(os.Stdout, conn)
		fmt.Println("done")

		done <- struct{}{}
	}()

	mustCopy(conn, os.Stdin) //block call
	conn.Close()

	<- done
}