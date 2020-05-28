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
	"github.com/hartfordfive/kafka-topic-tailer/lib"
	"github.com/hartfordfive/kafka-topic-tailer/version"
)

var (
	flagKafkaBrokers  = ""
	flagKafkaVersion  = ""
	flagTopics        = ""
	flagConsumerGroup = ""
	flagRegex         = ""
	flagLocalTZ       = ""
	flagFromOldest    = false
	flagVersion       = true
	flagIsJSON        = false
	flagDebug         = false
	brokers           = []string{}
	consumerGroup     = ""
)

func init() {
	flag.StringVar(&flagKafkaBrokers, "brokers", "", "Comma separated list of kafka brokers in IP:PORT format")
	flag.StringVar(&flagTopics, "topics", "", "Comma seperated list of topic to consume from")
	flag.StringVar(&flagConsumerGroup, "cg", "", "Custom consumer group name (auto-generated if left empty)")
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

	if len(flagTopics) < 1 {
		log.Fatal("Must specify a topic name!")
	}

	user, err := user.Current()
	if flagConsumerGroup != "" && len(flagConsumerGroup) >= 3 {
		consumerGroup = flagConsumerGroup
	} else {
		if err != nil {
			consumerGroup = fmt.Sprintf("default-topic-tailer-%s", lib.RndString(6))
		} else {
			consumerGroup = fmt.Sprintf("%s-topic-tailer-%s", user.Username, lib.RndString(6))
		}
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
		log.Printf("[DEBUG] flagTopics = %s", flagTopics)
		log.Printf("[DEBUG] flagRegex = %s", flagRegex)
		log.Printf("[DEBUG] consumerGroup = %s", consumerGroup)
		log.Printf("[DEBUG] flagLocalTZ = %s", flagLocalTZ)
		log.Printf("[DEBUG] flagFromOldest = %v", flagFromOldest)
		log.Printf("[DEBUG] flagIsJSON = %v", flagIsJSON)
		log.Printf("[DEBUG] flagDebug = %v", flagDebug)
	}

	// ---------------------------------

	topics := strings.Split(flagTopics, ",")
	client.Run(&client.Config{
		FilterRegex:   flagRegex,
		Brokers:       brokers,
		Topics:        topics,
		ConsumerGroup: consumerGroup,
		IsJSON:        flagIsJSON,
		Debug:         flagDebug,
		LocalTZ:       flagLocalTZ,
	}, config)

}
