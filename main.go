package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// User ..
type User struct {
	ID  int       `json:"id"`
	DOB time.Time `json:"dob"`
}

func main() {
	dataChan := make(chan []byte)
	var done = make(chan bool)
	go producer(dataChan, done)
	go consumer(dataChan)
	<-done
}

func producer(dataChan chan []byte, done chan bool) {
	for i := 0; i < 10; i++ {
		user := User{
			ID:  i,
			DOB: time.Now(),
		}
		userBytes, err := json.Marshal(user)
		if err != nil {
			log.Fatal(err)
		}
		dataChan <- userBytes
	}
	done <- true
}
func consumer(dataChan chan []byte) {
	for {
		data := <-dataChan
		fmt.Println(string(data))
	}
}
