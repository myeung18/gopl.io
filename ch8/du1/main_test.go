package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func walkDir_wg(dir string, fileSizes chan<- int64, wg sync.WaitGroup) {

	//Done for the upper caller thread.
	//or done for this level
	defer wg.Done()

	for _, entry := range dirents_wg(dir) {
		if entry.IsDir() {
			wg.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir_wg(subdir, fileSizes, wg)
		} else {
			fileSizes <- entry.Size()
		}
	} //for
}

var sema = make(chan struct{}, 5)

// dirents returns the entries of directory dir.
func dirents_wg(dir string) []os.FileInfo {
	sema <- struct{}{}
	defer func() { <-sema }()

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}


func TestCall_v2(t *testing.T) {

	roots := []string{".", "C:\\Users\\yeung\\gomod\\gopl.io"}

	var wg sync.WaitGroup

	fileSizes := make(chan int64)
	go func() {
		for _, r := range roots {
			wg.Add(1)
			go walkDir_wg(r, fileSizes, wg)
		}
	}()

	go func() {
		wg.Wait()
		close(fileSizes)
	}()

	//
	var tick <-chan time.Time
	var total int64

	tick = time.Tick(600 * time.Millisecond)
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				fmt.Println("...")
				break loop
			}

			time.Sleep(100 * time.Millisecond)

			fmt.Println("size is: ", size)
			total += size
		case <-tick:
			fmt.Println("total is: ", total)
		}
	}

	fmt.Println("Grand Total...", total)

}

func testCall(t *testing.T) {

	roots := []string{".", "C:\\Users\\yeung\\gomod\\gopl.io"}

	fileSizes := make(chan int64)
	go func() {
		for _, r := range roots {
			walkDir(r, fileSizes)
		}
		close(fileSizes)
	}()

	var total int64
	for size := range fileSizes {
		total += size
	}
	fmt.Println("total is: ", total)

}
