package event_consumer

import (
	"github.com/tauadam/reading_list-bot/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("Error fetching events: %v", err)
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Printf("Error handling events: %v", err)

			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, e := range events {
		log.Printf("Processing event: %v", e)
		if err := c.processor.Process(e); err != nil {
			log.Printf("Error processing event: %v", err)

			continue
		}
	}

	return nil
}
