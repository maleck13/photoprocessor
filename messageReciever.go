package main

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
	"encoding/json"
	"strings"
	"io"
)

type Message struct {
	File string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func StartConsuming() {
	fmt.Println("start consuming")
	conn, err := amqp.Dial("amqp://guest:guest@192.168.59.103:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		3,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			//validate the message and pass on to be processed
			dec := json.NewDecoder(strings.NewReader(string(d.Body)))
			var m Message
			if err := dec.Decode(&m); err == io.EOF {
				break
			} else if err != nil {
				d.Ack(false)
				log.Fatal(err)
			}
			fmt.Printf("%s \n", m.File)
		//	c:=make(chan int)
			ProcessImg(m.File)
			d.Ack(true)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
