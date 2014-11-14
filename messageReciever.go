package main

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
	"encoding/json"
	"strings"
	"io"
	"time"
)

type Message struct {
	File string
	User string
	RESKEY string
}

type AmqpCon struct {
	CONNECTION *amqp.Connection
}




func StartConsuming() {
	fmt.Println("start consuming")
	conn, err := amqp.Dial(CONF.GetRabbitURL())

	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()



	q,err := ch.QueueDeclare(
		"pics", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
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
			updates:= make(chan string)
			go UpdateJob(conn,m.RESKEY,updates)
			go ProcessImg(m.File,Picture{},m.User,CONF, updates)
			fmt.Printf("finished with file %s \n", m.File)

			d.Ack(true)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	//casues the thread to block
	<-forever
}


func setUpResponseExchange(connection * amqp.Connection) (*amqp.Channel,error){
	channel, err := connection.Channel()
	if err != nil {
		return nil,fmt.Errorf("Channel: %s", err)
	}

	return channel,err

}

func UpdateJob(conn * amqp.Connection, resKey string, messages chan string ){
	channel,err:= setUpResponseExchange(conn)
	FailOnError(err, "Failed to open a channel")
	defer channel.Close()
	_,err =channel.QueueDeclare(resKey,false, false, false, false, amqp.Table{})
	if err != nil{

	}

	for m:= range messages{

		fmt.Println("message ready " + m)

		if err:= channel.Publish(
			"",   // publish to an exchange
			resKey , // routing to 0 or more queues
			false,      // mandatory
			false,      // immediate
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

	fmt.Println("finished recieving message updated giving an extra 5 mins before removing que");
		//already in go routine so should be ok
		time.Sleep(time.Minute * 5)
		fmt.Println("deleting q")
		_, err = channel.QueueDelete(resKey, false, true, false);
		if err != nil{
			fmt.Println("deleting q err " + err.Error() )
		}
}

type UPDATE_MESSAGE struct {
	Message string
	Status string
}

func CreateMessage(message,status string)string{
	msg:=UPDATE_MESSAGE{message,status}
	json,err:=json.Marshal(msg)
	if err !=nil{
		fmt.Println("error " + err.Error())
	}
    return string(json)
}
