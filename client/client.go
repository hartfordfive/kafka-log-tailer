package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
)

// Config specifies the configuration to use for the consumer client
type Config struct {
	FilterRegex   string
	Brokers       []string
	Topics        []string
	ConsumerGroup string
	IsJSON        bool
	LocalTZ       string
	Debug         bool
}

// Run creates and configures the Kafka consumer to consume logs from the indicated topic
//func Run(brokers []string, topic, consumerGroup string, isJSON bool, config *sarama.Config) {
func Run(clientConfig *Config, config *sarama.Config) {
	consumer := Consumer{
		Ready:   make(chan bool),
		IsJSON:  clientConfig.IsJSON,
		Debug:   clientConfig.Debug,
		LocalTZ: clientConfig.LocalTZ,
	}

	_, regexpErr := regexp.Compile(clientConfig.FilterRegex)

	if regexpErr != nil {
		log.Printf("[ERROR] Regex compilation error: %s\n", regexpErr)
	}

	if clientConfig.FilterRegex != "" && regexpErr == nil {
		consumer.FilterRegex = clientConfig.FilterRegex
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(clientConfig.Brokers, clientConfig.ConsumerGroup, config)
	if err != nil {
		log.Fatalf("[FATAL] Could not create consumer group client: %v\n", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, clientConfig.Topics, &consumer); err != nil {
				log.Fatalf("[FATAL] From consumer: %v\n", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	log.Println(fmt.Sprintf("[INFO] Consuming logs from %s\n", strings.Join(clientConfig.Topics, ", ")))

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("[INFO] Terminating...")
	case <-sigterm:
		log.Println("[INFO] Shutting down...")
	}
	cancel()
	wg.Wait()

	log.Printf("[INFO] Total bytes consumed: %d\n", consumer.bytesConsumed)
	log.Printf("[INFO] Total bytes displayed: %d\n", consumer.bytesDisplayed)

	if err = client.Close(); err != nil {
		log.Fatalf("[FATAL] Error closing client: %v\n", err)
	}
}
