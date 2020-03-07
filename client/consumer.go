package client

import (
	"github.com/Shopify/sarama"
	"github.com/fatih/color"
	"github.com/pquerna/ffjson/ffjson"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	Ready  chan bool
	IsJSON bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	var msg map[string]interface{}
	cY := color.New(color.FgYellow).Add(color.Bold)
	cW := color.New(color.FgWhite)

	for message := range claim.Messages() {

		if consumer.IsJSON {
			ffjson.Unmarshal(message.Value, &msg)
			cY.Printf("[%s]", msg["@timestamp"].(string))
			cW.Printf(" %s\n", msg["message"].(string))
		} else {
			cW.Printf(" %s\n", string(message.Value))
		}
		session.MarkMessage(message, "")
	}

	return nil
}
