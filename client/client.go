package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
)

// Config specifies the configuration to use for the consumer client
type Config struct {
	FilterRegex   string
	Brokers       []string
	Topic         string
	ConsumerGroup string
	IsJSON        bool
}

// Run creates and configures the Kafka consumer to consume logs from the indicated topic
//func Run(brokers []string, topic, consumerGroup string, isJSON bool, config *sarama.Config) {
func Run(clientConfig *Config, config *sarama.Config) {
	consumer := Consumer{
		Ready:  make(chan bool),
		IsJSON: clientConfig.IsJSON,
	}

	_, regexpErr := regexp.Compile(clientConfig.FilterRegex)

	if clientConfig.FilterRegex != "" && regexpErr != nil {
		consumer.FilterRegex = clientConfig.FilterRegex
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(clientConfig.Brokers, clientConfig.ConsumerGroup, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, []string{clientConfig.Topic}, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	log.Println(fmt.Sprintf("Consuming logs from %s\n", clientConfig.Topic))

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("Terminating...")
	case <-sigterm:
		log.Println("Shutting down...")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}
