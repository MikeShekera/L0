package testdata

import (
	"02.08.2024-L0/services"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/stretchr/testify/mock"
	"testing"
)

type NATSMock struct {
	mock.Mock
}

func (m *NATSMock) Subscribe(subject string, cb stan.MsgHandler, opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	args := m.Called(subject, cb, opts)
	return nil, args.Error(0)
}

func (m *NATSMock) Publish(subj string, data []byte) error {
	args := m.Called(subj, data)
	return args.Error(0)
}

func TestSubscribeNATS(t *testing.T) {
	mockConn := new(NATSMock)
	mockConn.On("Subscribe", "foo", mock.Anything, mock.Anything).Return(nil)

	err := services.SubscribeNATS(mockConn, func(m *stan.Msg) {})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	mockConn.AssertExpectations(t)
}

func TestConnectNats(t *testing.T) {
	err, _ := services.ConnectNats("test-client")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPublishNATS(t *testing.T) {
	mockConn := new(NATSMock)
	mockConn.On("Publish", "test-subject", []byte("test-message")).Return(nil)

	services.PublishNATS(mockConn, "test-subject", []byte("test-message"))

	mockConn.AssertExpectations(t)
}

func (m *NATSMock) PublishAsync(subject string, data []byte, ah stan.AckHandler) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *NATSMock) QueueSubscribe(subject, qgroup string, cb stan.MsgHandler, opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	//TODO implement me
	panic("implement me")
}

func (m *NATSMock) Close() error {
	//TODO implement me
	panic("implement me")
}

func (m *NATSMock) NatsConn() *nats.Conn {
	//TODO implement me
	panic("implement me")
}
