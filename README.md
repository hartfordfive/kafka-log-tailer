
## kafka-topic-tailer

### Description

This utility is intended to allow you to directly tail logs from kafka, which have been published there via filebeat, instead of waiting until they've been processed and become available at the final destination (typically Elasticsearch).  

### Flags

* `-brokers` : (string) Comma separated list of brokers in IP:PORT format
* `-json` : Parse log entry and only display `@timestamp`, `beat.hostname` and `message` fields (default: `false`)
* `-kver` : (string) Version of Kafka (default: `2.1.0`)
* `-oldest` : Start the kafka consumer from oldest ofset (default: `true`)
* `-topic` : (string) Name of the kafka topic to consume from
* `-r` : (string) Regex to filter for specific messages
* `-tz` :  (string) Convert the **@timestamp** value to a timezone of your choice. (default: `Etc/UTC`)
* `-d` : Enable debug mode logging
* `-v` : Print version info and exit

### Command usage example

./kafka-topic-tailer -brokers "kafka01:2181,kafka02:2181,kafka03:2181" -topic "my-log-topic" -tz "America/Chicago"

### Download pre-compiled builds




### Building from source

Building for your specific OS: 
* `make build`

Building for the three main OS varieties (Linux, Darwin & Windows):
* `make build-all`

Building with debug symbols present:
* `make build-debug`

*See Makefile for all available options*

