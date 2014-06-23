package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func Consume(uri string, doneChan chan bool) {
	log.Println("Consuming...")
	connection, err := amqp.Dial(uri)
	if err != nil {
		println(err.Error())
		panic(err.Error())
	}
	defer connection.Close()

	channel, err1 := connection.Channel()
	if err1 != nil {
		println(err1.Error())
		panic(err1.Error())
	}
	defer channel.Close()

	q := MakeQueue(channel)

	msgs, err3 := channel.Consume(q.Name, "", true, false, false, false, nil)
	if err3 != nil {
		panic(err3)
	}

	for d := range msgs {
		doneChan <- true
		var thisMessage MqMessage
		err4 := json.Unmarshal(d.Body, &thisMessage)
		if err4 != nil {
			log.Printf("Error unmarshalling! %s", err.Error())
		}
		log.Printf("Message age: %s", time.Since(thisMessage.TimeNow))

	}

	log.Println("done recieving")

}
