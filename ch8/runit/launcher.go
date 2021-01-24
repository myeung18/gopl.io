package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	wd_launcher(".")
}

func wd_launcher(dir string) {
	roots := []string{dir}
	fileSizes := make(chan int64)


	var n sync.WaitGroup
	for _, root := range roots {

		n.Add(1)
		go wd(root, fileSizes, &n)
	}

	go func() {
		n.Wait()
		close(fileSizes) //close!!
	}()

	//for si := range fileSizes { //stop if it is closed!!!
	//	fmt.Println("size: ", si)
	//}

	var tick <-chan time.Time
	tick = time.Tick(2 * time.Millisecond)

	var total int64

loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop
			}
			total += size
		case <-tick:
			fmt.Println("file in tick : ", total)
		}
	}
	fmt.Println("file in tick-final : ", total)
}

func wd(dir string, fileSizes chan<- int64, n *sync.WaitGroup) {
	defer n.Done()

	for _, entry := range wd_2(dir) {
		if entry.IsDir() {

			n.Add(1)

			subdir := filepath.Join(dir, entry.Name())
			go wd(subdir, fileSizes, n)

		} else {
			fileSizes <- entry.Size()
		}
	} //for
}

var sema = make(chan struct{}, 1)

func wd_2(dir string) []os.FileInfo {
	sema <- struct{}{}
	defer func() { <- sema}() //the last func to call!!!

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

/*********************************************************************/
func a() {
	ch := make(chan bool)
	go func() {
		os.Stdin.Read(make([]byte, 1))
		fmt.Println("aborted")
		ch <- false
	}()

	fmt.Println("commercing...")

	var test chan int
	test = make(chan int)
	test <- 10

	//tick := time.Tick(1 * time.Second)

	tick := time.NewTicker(3 * time.Second)
	for i := 10; i > 0; i-- {
		fmt.Print("timing ticker ....", i)
		select { //block and wait
		case <-tick.C:
			fmt.Println(" Do nothing..")

		case <-ch:
			fmt.Println("Launch aborted!! ")
			return
			//default:
			//	fmt.Println("no event yet")
		}
	}
	other()
}

func other() {
	fmt.Println("from others")
}

func b() {
	ch := make(chan int, 2)
	for i := 0; i < 10; i++ {
		select {
		case x := <-ch: /** it is empty at the beginning  */
			fmt.Println(x)
		case ch <- i:
			fmt.Println("sent")
		}
	} //for
}

func abort() {
	ch := make(chan bool)
	go func() {
		os.Stdin.Read(make([]byte, 1)) //block and wait
		fmt.Println("aborted")
		ch <- false
	}()
}
