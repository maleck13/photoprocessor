package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"io"
	"log"
	"strings"
	"time"
)

type Message struct {
	File   string
	User   string
	RESKEY string
	Name   string
}

type AmqpCon struct {
	CONNECTION *amqp.Connection
}

func connect() (*amqp.Connection, error) {

	conn, err := amqp.Dial(CONF.GetRabbitURL())
	//connClos <- amqp.ErrClosed
	return conn, err
}

func connectionListener(connErr chan *amqp.Error, stopChan chan bool) {
	ErrorLog.Println("Started connection Listener")
	errMess := <-connErr
	if errMess.Code == amqp.FrameError || errMess.Code == amqp.ChannelError {
		fmt.Println("recieed amqp error trying reconnect", errMess)
		close(stopChan)

		var conn *amqp.Connection
		var err error
		for {
			fmt.Println("Trying to reconnect")
			conn, err = connect()
			if err != nil {
				fmt.Println(" failed to connect " + err.Error())
				time.Sleep(2000)
			} else {
				break
			}

		}
		startConsuming(conn)

	}

}

func StartUp() {
	conn, err := connect()
	FailOnError(err, "failed to connect")
	startConsuming(conn)
}

func startConsuming(conn *amqp.Connection) {

	defer conn.Close()
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"pics", // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		3,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")
	//casues the thread to block
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
				d.Ack(true)
				ErrorLog.Println("error with rabbit msg " + err.Error())
			}
			fmt.Printf("%s \n", m.File)
			updates := make(chan string)
			go UpdateJob(conn, m.RESKEY, updates)
			go ProcessImg(m.Name, Picture{}, m.User, CONF, updates, m.RESKEY)
			fmt.Printf("finished with file %s \n", m.File)

			d.Ack(true)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	//casues the thread to block
	connClos := make(chan *amqp.Error)
	conn.NotifyClose(connClos)
	connectionListener(connClos, forever)
	<-forever

	fmt.Println("stopping listening ")
}

func UpdateJob(conn *amqp.Connection, resKey string, messages chan string) {
	channel, err := conn.Channel()
	if err != nil {
		FailOnError(err, "Failed to open a channel")
	}

	defer channel.Close()

	if err != nil {
		FailOnError(err, "failed to declare que ")
	}

	for m := range messages {

		fmt.Println("message ready " + m + "publishing to " + resKey)

		if err := channel.Publish(
			"amq.topic",             // publish to an exchange
			"picjob.update."+resKey, // routing to 0 or more queues
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType:     "application/json",
				ContentEncoding: "utf8",
				Body:            []byte(m),
				DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
			},
		); err != nil {
			fmt.Printf("Exchange Publish: %s", err)
		}

	}
}

type UPDATE_MESSAGE struct {
	Message string
	Status  string
	Jobid   string
	Type    string
}

func CreateMessage(message, status, jobid string) string {
	msg := UPDATE_MESSAGE{message, status, jobid, "PICTURE"}
	json, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error " + err.Error())
	}
	return string(json)
}
