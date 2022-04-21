package main

import (
	"github.com/Shopify/sarama"
	"github.com/jgunvalson/kafka-producer/config"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)


func main() {
	producerConfig, err := config.LoadConfigOrDie()
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("connecting to kafka cluster at address=%v for topics=%v", producerConfig.BootstrapServer, producerConfig.Topics)
	producer, err := sarama.NewAsyncProducer(producerConfig.BootstrapServer, producerConfig.Config)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = producer.Close()
		log.Error(err)
	}()

	produceTicker := time.NewTicker(time.Millisecond * time.Duration(producerConfig.ProduceBackoffMs))
	var count = 1
	go func () {
		for {
			select {
			case <-produceTicker.C:
				log.Infof("%v: producing message n=%d to topic...", time.Now(), count)
				producer.Input() <- &sarama.ProducerMessage{
					Topic: producerConfig.Topics[0],
					Key:   sarama.StringEncoder("key"),
					Value: sarama.StringEncoder("value"),
				}
				count++
			}
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(
		sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	select {
		case sig := <-sigChan:
			log.Infof("received signal: %s", sig.Signal)
			os.Exit(0)
		case err = <-producer.Errors():
			log.Fatalf("received error producing: %s", err.Error())
			err.Error()
	}
}
