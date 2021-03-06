package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
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
	size = 100
	queueLock = &sync.Mutex{}
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
}

func producer(dataChan chan []byte, done chan bool, number int) {
	for i := 0; i < 2; i++ {
		user := User{
			ID:  strconv.Itoa(i) + ":" + strconv.Itoa(number),
			DOB: time.Now(),
		}
		userBytes, err := json.Marshal(user)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("producer", user.ID)
		enqueue(userBytes)
	}
	done <- true
}
func consumer(dataChan chan []byte, number int) {
	for {
		time.Sleep(time.Duration(int64(timeInterval) * int64(time.Millisecond)))
		// fmt.Println(queue)
		data := dequeue()
		fmt.Println("Consumer", number, string(data))
	}
}

var queue [][]byte
var size int
var front, end int
var queueLock *sync.Mutex

func enqueue(data []byte) (err error) {
	queueLock.Lock()
	if len(queue) >= size {
		queueLock.Unlock()

		time.Sleep(time.Duration(int64(timeInterval) * int64(time.Millisecond)))

		return enqueue(data)
	}
	queue = append(queue, data)
	queueLock.Unlock()
	return
}
func dequeue() (data []byte) {
	queueLock.Lock()
	if len(queue) == 0 {
		time.Sleep(time.Duration(int64(timeInterval) * int64(time.Millisecond)))
		return dequeue()
	}
	dataArr := queue[:1]
	for _, data = range dataArr {
		break
	}
	queue = queue[1:]
	queueLock.Unlock()
	return
}
