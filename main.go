package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// User ..
type User struct {
	ID  string    `json:"id"`
	DOB time.Time `json:"dob"`
}

var timeInterval int
var numProducer, numConsumer int

func init() {
	viper.SetConfigName("app")
	viper.SetConfigType("json")
	viper.AddConfigPath(os.Getenv("CONFIG_PATH"))
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		log.Println("Unable to Read file")
		return
	}
	timeInterval = viper.GetInt("time_interval")
	numConsumer = viper.GetInt("consumers")
	numProducer = viper.GetInt("producers")
	fmt.Println("ConfigPath", os.Getenv("CONFIG_PATH"))
	fmt.Println("TimeInt", timeInterval, "NumConsumers", numConsumer, "NumProducers", numProducer)
}
func main() {
	dataChan := make(chan []byte)
	var done = make(chan bool)
	for i := 0; i < numProducer; i++ {
		go producer(dataChan, done, i)
	}
	for j := 0; j < numConsumer; j++ {
		go consumer(dataChan, j)
	}
	<-done
	close(dataChan)
}

func producer(dataChan chan []byte, done chan bool, number int) {
	for i := 0; i < 100; i++ {
		user := User{
			ID:  strconv.Itoa(i) + ":" + strconv.Itoa(number),
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
func consumer(dataChan chan []byte, number int) {
	for {
		time.Sleep(time.Duration(int64(timeInterval) * int64(time.Millisecond)))
		data := <-dataChan
		fmt.Println("Consumer", number, string(data))
	}
}
