package server

import (
	"github.com/go-kit/kit/endpoint"
)

// KafkaConsumer encapsulates all members used for the kafka server.
type KafkaConsumer struct {
	//cfg    *configuration.KafkaConfig
	//config *sarama.Config
	//master sarama.Consumer
	topics map[string]Topic
	quit   chan bool
}

// Topic encapsulates a kafka channel name and a map of endpoints used on that channel
type Topic struct {
	Name               string
	ServiceEndpointMap ServiceEndpointMap
}

// MyTopic returns a new Topic instance
func MyTopic(name string) (topic Topic) {
	return Topic{
		Name:               name,
		ServiceEndpointMap: make(ServiceEndpointMap),
	}
}

// MyKafkaConsumer takes in configuration object and returns a new KafkaConsumer instance.
func MyKafkaConsumer() (kc *KafkaConsumer, err error) {
	//config := sarama.NewConfig()
	//config.Consumer.Return.Errors = true

	//master, err := sarama.NewConsumer(cfg.Brokers, config)
	if err != nil {
		return
	}

	kc = &KafkaConsumer{
		//cfg,
		//config,
		//master,
		make(map[string]Topic),
		make(chan bool),
	}

	return
}

// Start is an implementation of the Server Start() method.
func (kc *KafkaConsumer) Start() (err error) {
	for _, topic := range kc.topics {
		go func(topic Topic) {
			//consumer(kc.master, kc.quit, topic, ByteHandler, kc.cfg)
			select {
			case <-kc.quit:
				return
			}
		}(topic)
	}
	return
}

// Stop is an implementation of the Server Stop() method
func (kc *KafkaConsumer) Stop() (err error) {
	close(kc.quit)
	//kc.master.Close()
	return
}

// RegisterNamespace is an implementation of the Server RegisterNamespace method.
// In kafka a namespace is the same as the channel name.
func (kc *KafkaConsumer) RegisterNamespace(name string) {
	kc.topics[name] = MyTopic(name)
	return
}

// RegisterService is an implementation of the server RegisterService method.
// In kafka this is used to demux the messages coming in on a single channel/topic.
func (kc *KafkaConsumer) RegisterService(namespace string, service Service, ep endpoint.Endpoint) {
	kc.topics[namespace].ServiceEndpointMap[service] = ep
	return
}
