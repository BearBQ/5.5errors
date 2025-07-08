package main

import (
	"fmt"
	_ "fmt"
	"sync"
)

func errFunc(i int, intChan chan<- int, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	if i == 0 {

		errChan <- fmt.Errorf("Вводное число не может быть нулем")
		return
	}

	intChan <- i * i
}

func main() {
	intChan := make(chan int)
	errChan := make(chan error)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case num := <-intChan:
				fmt.Println("значение: ", num)
			case err := <-errChan:
				fmt.Println("Значение ошибки:", err)
			case <-done:
				return
			}

		}
	}()

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go errFunc(i, intChan, &wg, errChan)
	}

	wg.Wait()
	close(done)
}
