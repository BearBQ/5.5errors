package main

import (
	"fmt"
	_ "fmt"
	"sync"
	"time"
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
	var errors []error

	go func() { //выводим результат
		for {
			select {
			case num := <-intChan:
				fmt.Println("значение: ", num)
			case <-done:
				return

			}

		}
	}()

	go func() {
		select {
		case err := <-errChan:
			errors = append(errors, err)
		case <-done:
			return

		}

	}()

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go errFunc(i, intChan, &wg, errChan)
	}

	wg.Wait()

	time.Sleep(2 * time.Second)
	for _, err := range errors {
		fmt.Println(err)
	}
	close(done)
	close(intChan)
	close(errChan)
}
