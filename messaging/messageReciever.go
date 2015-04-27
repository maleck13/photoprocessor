package messaging

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"io"
	"log"
	"strings"
	"time"
	"github.com/maleck13/photoProcessor/logger"
	"github.com/maleck13/photoProcessor/processor"
	"github.com/maleck13/photoProcessor/errorHandler"
	"github.com/maleck13/photoProcessor/model"
	"github.com/maleck13/photoProcessor/conf"
)





type AmqpCon struct {
	CONNECTION *amqp.Connection
}

var conn * amqp.Connection;

func connect() (*amqp.Connection, error) {

	connect,err  := amqp.Dial(conf.CONF.GetRabbitURL())

	return connect, err
}

func connectionListener(connErr chan *amqp.Error, stopChan chan bool) {
	logger.ErrorLog.Println("Started connection Listener")
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
		startMessageConsuming(conn)
	}

}

func StartMessaging(){
	conn, err := connect()
	errorHandler.FailOnError(err, "failed to connect")
	startMessageConsuming(conn)
}

func startMessageConsuming(c *amqp.Connection) {
	conn = c;
	go func() {
		defer c.Close()
		ch, err := c.Channel()
		errorHandler.FailOnError(err, "Failed to open a channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"pics", // name
			true,   // durable
			false,  // delete when unused
			false,  // exclusive
			false,  // no-wait
			nil,    // arguments
		)
		errorHandler.FailOnError(err, "Failed to declare a queue")

		err = ch.Qos(
			3,     // prefetch count
			0,     // prefetch size
			false, // global
		)
		errorHandler.FailOnError(err, "Failed to set QoS")

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			false,  // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		errorHandler.FailOnError(err, "Failed to register a consumer")
		//casues the thread to block
		forever := make(chan bool)

		go func() {
			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)
				//validate the message and pass on to be processed
				dec := json.NewDecoder(strings.NewReader(string(d.Body)))
				var m model.Message
				if err := dec.Decode(&m); err == io.EOF {
					break
				} else if err != nil {
					d.Ack(true)
					logger.ErrorLog.Println("error with rabbit msg " + err.Error())
				}
				fmt.Printf("%s \n", m.File)
				updates := make(chan string)
				go UpdateJob(m.ResKey, updates)
				go processor.ProcessImg(m.Name, model.Picture{}, m.User, updates, m.ResKey)
				fmt.Printf("finished with file %s \n", m.File)

				d.Ack(true)
			}
		}()

		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		//casues the thread to block
		connClose := make(chan *amqp.Error)
		c.NotifyClose(connClose)
		connectionListener(connClose, forever)
		<-forever

		fmt.Println("stopping listening ")
	}()
}

func SetUpResponseQue(id string){
	ch, err := conn.Channel()
	errorHandler.FailOnError(err, "Failed to open a channel")
	//defer ch.Close()
	fmt.Println("decaring queue " + id)
	var timeout int32 = 60000 * 5
	_, err = ch.QueueDeclare(
		id, // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		amqp.Table{"x-expires":timeout});
	ch.QueueBind(id,"picjob.update."+id,"amq.topic",false,nil)
}

func UpdateJob(resKey string, messages chan string) {
	channel, err := conn.Channel()
	if err != nil {
		errorHandler.FailOnError(err, "Failed to open a channel")
	}

	defer channel.Close()

	if err != nil {
		errorHandler.FailOnError(err, "failed to declare que ")
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


