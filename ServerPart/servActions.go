package ServerPart

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
)

const (
	ClusterID = "test-cluster"
	NatsURL   = "0.0.0.0:4222"
)

func ConnectNats(clientID string) (error, stan.Conn) {
	sc, err := stan.Connect(ClusterID, clientID, stan.NatsURL(NatsURL))
	if err != nil {
		return err, nil
	}
	fmt.Printf("Connected Successfully to %s \n", NatsURL)
	return nil, sc
}

func SubscribeNATS(sc stan.Conn, receiveFunc func(m *stan.Msg)) error {
	_, err := sc.Subscribe(
		"foo", receiveFunc, stan.DeliverAllAvailable(),
	)
	if err != nil {
		return err
	}
	return nil
}

func PublishNATS(sc stan.Conn, subj string, data []byte) {
	err := sc.Publish(subj, data)
	if err != nil {
		log.Fatalf("Error publishing message: %v", err)
	}
}
