package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/hartfordfive/kafka-topic-tailer/client"
	"github.com/hartfordfive/kafka-topic-tailer/version"
)

var (
	flagKafkaBrokers = ""
	flagKafkaVersion = ""
	flagTopic        = ""
	flagRegex        = ""
	flagLocalTZ      = ""
	flagFromOldest   = false
	flagVersion      = true
	flagIsJSON       = false
	flagDebug        = false
	brokers          = []string{}
	consumerGroup    = ""
)

func init() {
	flag.StringVar(&flagKafkaBrokers, "brokers", "", "Comma separated list of kafka brokers in IP:PORT format")
	flag.StringVar(&flagTopic, "topic", "", "Name of the kafka topic to consume from")
	flag.StringVar(&flagKafkaVersion, "kver", "2.1.0", "Version of Kafka")
	flag.BoolVar(&flagFromOldest, "oldest", false, "Start the kafka consumer from oldest ofset")
	flag.BoolVar(&flagIsJSON, "json", false, "Parse log entry and only display `@timestamp`, `beat.hostname` and `message` fields")
	flag.StringVar(&flagRegex, "r", "", "Regex to isolate specific messages")
	flag.StringVar(&flagLocalTZ, "tz", "Etc/UTC", "Convert the `@timestamp` value to a timezone of your choice.")
	flag.BoolVar(&flagVersion, "v", false, "Print version info and exit")
	flag.BoolVar(&flagDebug, "d", false, "Enable debug mode logging")
	flag.Parse()

	if flagVersion {
		version.PrintVersion()
		os.Exit(0)
	}

	brokers = strings.Split(strings.Trim(flagKafkaBrokers, " "), ",")
	if len(brokers) < 1 {
		log.Fatal("At least one broker must be specified!")
	}

	if len(flagTopic) < 1 {
		log.Fatal("Must specify a topic name!")
	}

	user, err := user.Current()
	if err != nil {
		consumerGroup = "default-topic-tailer"
	} else {
		consumerGroup = fmt.Sprintf("%s-topic-tailer", user.Username)
	}

}

func main() {
	version, err := sarama.ParseKafkaVersion(flagKafkaVersion)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()
	config.Version = version
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	if flagFromOldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	if flagDebug {
		log.Printf("[DEBUG] flagKafkaBrokers = %s", flagKafkaBrokers)
		log.Printf("[DEBUG] flagKafkaVersion = %s", flagKafkaVersion)
		log.Printf("[DEBUG] flagTopic = %s", flagTopic)
		log.Printf("[DEBUG] flagRegex = %s", flagRegex)
		log.Printf("[DEBUG] consumerGroup = %s", consumerGroup)
		log.Printf("[DEBUG] flagLocalTZ = %s", flagLocalTZ)
		log.Printf("[DEBUG] flagFromOldest = %v", flagFromOldest)
		log.Printf("[DEBUG] flagIsJSON = %v", flagIsJSON)
		log.Printf("[DEBUG] flagDebug = %v", flagDebug)
	}

	// ---------------------------------

	client.Run(&client.Config{
		FilterRegex:   flagRegex,
		Brokers:       brokers,
		Topic:         flagTopic,
		ConsumerGroup: consumerGroup,
		IsJSON:        flagIsJSON,
		Debug:         flagDebug,
		LocalTZ:       flagLocalTZ,
	}, config)

}
