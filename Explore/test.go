package main

import (
	"fmt"
	"time"
)

type str struct {
	val int
}

func createJob(val int) (string, error) {
	st := &str{}
	errChan := make(chan error)
	go func(errChan chan error) {
		defer close(errChan)
		if val > 1 {
			st.val = val
			time.Sleep(time.Second * 2)
		} else {
			err := fmt.Errorf("invalid amount %v", val)
			errChan <- err
		}
	}(errChan)

	for err := range errChan {
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("Struct is: %v", st), nil
}
func main() {
	if s, err := createJob(-4); err != nil {
		panic(err)
	} else {
		fmt.Println(s)
	}

}
