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
	flagFromOldest   = true
	flagVersion      = true
	flagIsJSON       = false
	brokers          = []string{}
	consumerGroup    = ""
)

func init() {
	flag.StringVar(&flagKafkaBrokers, "brokers", "", "Comma separated list of brokers in IP:PORT format")
	flag.StringVar(&flagTopic, "topic", "", "Name of the log topic to consume from")
	flag.StringVar(&flagKafkaVersion, "kver", "2.1.0", "Version of Kafka")
	flag.BoolVar(&flagFromOldest, "oldest", true, "Kafka consumer consume initial offset from oldest")
	flag.BoolVar(&flagIsJSON, "json", false, "Messages in the topic are json compliant payloads")
	flag.StringVar(&flagRegex, "r", "", "Regex to isolate specific messages")
	flag.BoolVar(&flagVersion, "v", false, "Print version info and exit")
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
		consumerGroup = "default-log-tailer"
	} else {
		consumerGroup = fmt.Sprintf("%s-log-tailer", user.Username)
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

	if flagFromOldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	// ---------------------------------

	client.Run(&client.Config{
		FilterRegex:   flagRegex,
		Brokers:       brokers,
		Topic:         flagTopic,
		ConsumerGroup: consumerGroup,
		IsJSON:        flagIsJSON,
	}, config)

	//client.Run(brokers, flagTopic, consumerGroup, flagIsJSON, config)
}
