package beater

import (
	"context"
	"fmt"
	"time"

	"github.com/S1mplee/eventstorebeat/config"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/vectorhacker/goro"
)

// Eventstorebeat configuration.
type Eventstorebeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of eventstorebeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig

	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Eventstorebeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts eventstorebeat.
func (bt *Eventstorebeat) Run(b *beat.Beat) error {
	logp.Info("eventstorebeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	client := goro.Connect(bt.config.Brokers[0], goro.WithBasicAuth(bt.config.Username, bt.config.Password))

	ctx := context.Background()
	reader := client.FowardsReader(bt.config.Stream)
	catchupSubscription := client.CatchupSubscription(bt.config.Stream, 0) // start from 0

	go func() {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		messages := catchupSubscription.Subscribe(ctx)

		for message := range messages {
			if err := message.Error; err != nil {
				return
			}

			event2 := beat.Event{
				Timestamp: time.Now(),
				Fields: common.MapStr{
					"type":        b.Info.Name,
					"eventType":   message.Event.Type,
					"eventID":     message.Event.ID,
					"eventAuthor": message.Event.Author,
					"eventData":   string(message.Event.Data),
				},
			}

			bt.client.Publish(event2)
		}
	}()

	events, err := reader.Read(ctx, 0, 1)
	if err != nil {
		panic(err)
	}

	for _, event := range events {
		fmt.Printf("%s\n", event.Data)
	}

	return nil

}

// Stop stops eventstorebeat.
func (bt *Eventstorebeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
