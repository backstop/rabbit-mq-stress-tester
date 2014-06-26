rabbit-mq-stress-tester
=======================


Compiling
---------

    go build tester.go producer.go consumer.go

That will produce the executible `tester`

Running
-------

    $ ./tester -h
	NAME:
	   tester - Make the rabbit cry

	USAGE:
	   tester [global options] command [command options] [arguments...]

	VERSION:
	   0.0.0

	COMMANDS:
	   help, h	Shows a list of commands or help for one command

	GLOBAL OPTIONS:
	   --server, -s 'rabbit-mq-test.cs1cloud.internal'	Hostname for RabbitMQ server
	   --producer, -p '0'				Number of messages to produce, -1 to produce forever
	   --wait, -w '0'					Number of nanoseconds to wait between publish events
	   --consumer, -c '-1'				Number of messages to consume. 0 consumes forever
	   --bytes, -b '0'					number of extra bytes to add to the RabbitMQ message payload. About 50K max
	   --concurrency, -n '50'			number of reader/writer Goroutines
	   --quiet, -q						Print only errors to stdout
	   --wait-for-ack, -a				Wait for an ack or nack after enqueueing a message
	   --version, -v					print the version
	   --help, -h						show help

Examples
--------


Consume messages forever:

	./tester -s rabbit-mq-test.cs1cloud.internal -c 0


Produce 100,000 messages of 10KB each, using 50 concurrent goroutines, waiting 100 nanoseconds between each message. Only print to stdout if there is a nack or when you finish.

	./tester -s rabbit-mq-test.cs1cloud.internal -p 100000 -b 10000 -w 100 -n 50 -q