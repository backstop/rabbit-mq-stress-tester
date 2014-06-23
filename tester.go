package main

import (
	"github.com/codegangsta/cli"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
)

var totalTime int64 = 0
var totalCount int64 = 0

type MqMessage struct {
	TimeNow        time.Time
	SequenceNumber int
	Payload        string
}

func main() {
	app := cli.NewApp()
	app.Name = "RabbitMQ-Stress-Tester"
	app.Usage = "Make the rabbit cry"
	app.Flags = []cli.Flag{
		cli.StringFlag{"server, s", "rabbit-mq-test.cs1cloud.internal", "Hostname for RabbitMQ server"},
		cli.IntFlag{"producer, p", 0, "Number of messages to produce, -1 to produce forever"},
		cli.IntFlag{"wait, w", 0, "Number of nanoseconds to wait between publish events"},
		cli.IntFlag{"consumer, c", -1, "Number of messages to consume. 0 consumes forever"},
		cli.IntFlag{"bytes, b", 0, "number of extra bytes to add to the RabbitMQ message payload. About 50K max"},
		cli.IntFlag{"concurrency, n", 50, "number of reader/writer Goroutines"},
		cli.BoolFlag{"quiet, q", "Print only errors to stdout"},
	}
	app.Action = func(c *cli.Context) {
		runApp(c)
	}
	app.Run(os.Args)
}

func runApp(c *cli.Context) {
	println("Running!")
	uri := "amqp://guest:guest@" + c.String("server") + ":5672"

	if c.Int("consumer") > -1 {
		makeConsumers(uri, c.Int("concurrency"), c.Int("consumer"))
	}

	if c.Int("producer") != 0 {
		makeProducers(c.Int("producer"), uri, c.Int("wait"), c.Int("bytes"), c.Int("concurrency"), c.Bool("quiet"))
	}
}

func MakeQueue(c *amqp.Channel) amqp.Queue {
	q, err2 := c.QueueDeclare("stress-test-exchange", true, false, false, false, nil)
	if err2 != nil {
		panic(err2)
	}
	return q
}

func makeProducers(n int, uri string, wait int, bytes int, concurrency int, quiet bool) {

	taskChan := make(chan int)
	for i := 0; i < concurrency; i++ {
		go Produce(uri, taskChan, bytes, quiet)
	}

	start := time.Now()

	for i := 0; i < n; i++ {
		taskChan <- i
		time.Sleep(time.Duration(int64(wait)))
	}

	time.Sleep(time.Duration(10000))

	close(taskChan)

	log.Printf("Finished: %s", time.Since(start))
}

func makeConsumers(uri string, concurrency int, toConsume int) {

	doneChan := make(chan bool)

	for i := 0; i < concurrency; i++ {
		go Consume(uri, doneChan)
	}

	start := time.Now()

	if toConsume > 0 {
		for i := 0; i < toConsume; i++ {
			<-doneChan
			if i == 1 {
				start = time.Now()
			}
			log.Println("Consumed: ", i)
		}
	} else {

		for {
			<-doneChan
		}
	}

	log.Printf("Done consuming! %s", time.Since(start))
}
