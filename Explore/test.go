package main

import (
	"fmt"
	"sync"
	"time"
)

type str struct {
	val1 int
	val2 int
}

func createJob(val1, val2 int) (string, error) {
	st := &str{}

	var wg sync.WaitGroup
	errChan1 := make(chan error, 2)
	//errChan2 := make(chan error)

	wg.Add(1)
	go func(wg *sync.WaitGroup, errChan chan error) {
		defer wg.Done()
		if val1 > 1 {
			st.val1 = val1
			time.Sleep(time.Second * 3)
			errChan <- nil
		} else {
			err := fmt.Errorf("too much for val1 : %v", val1)
			errChan <- err
		}
	}(&wg, errChan1)

	wg.Add(1)
	go func(errChan chan error) {
		defer wg.Done()

		if val2 < 0 {
			st.val2 = val2
			//time.Sleep(time.Second * 2)
			errChan <- nil
		} else {
			err := fmt.Errorf("too less for val2 : %v", val2)
			errChan <- err
		}
	}(errChan1)

	go func(wg *sync.WaitGroup, errChan chan error) {
		wg.Wait()
		close(errChan)
	}(&wg, errChan1)

	for err := range errChan1 {
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("Values are correct! Struct is: %v", st), nil
}

func display(s string, err error) {
	if err != nil {
		fmt.Println("ERROR!! >>> ", err)
	} else {
		fmt.Println(s)
	}
}

func main() {

	display(createJob(4, 2))
	display(createJob(-4, -2))
	display(createJob(-4, 2))
	display(createJob(4, -2))

}
