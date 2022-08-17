package queues

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type (
	Message  any
	Callback func(ctx context.Context, data []byte) error
	Client   interface {
		Publish(ctx context.Context, topicID string, message Message) error
		Consume(ctx context.Context, subscriptionID string, callback Callback) error
		Close() error
	}
	Subscription struct {
		Topic        string
		Subscription string
	}
)

type pubSub struct {
	client *pubsub.Client
	logger *zap.Logger
}

func NewClient(projectID string, subscriptions []Subscription, logger *zap.Logger) (Client, error) {
	pubSubClient, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, errors.Wrap(err, "init pubsub client")
	}

	if err := createResources(context.Background(), pubSubClient, subscriptions); err != nil {
		return nil, errors.Wrap(err, "create resources")
	}

	return &pubSub{client: pubSubClient, logger: logger}, nil
}

func (c *pubSub) Publish(ctx context.Context, topicID string, message Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "marshal message")
	}

	_, err = c.client.Topic(topicID).Publish(ctx, &pubsub.Message{Data: payload}).Get(ctx)

	return errors.Wrap(err, "waiting for message publishing")
}

func (c *pubSub) Consume(ctx context.Context, subscriptionID string, callback Callback) error {
	subscription := c.client.Subscription(subscriptionID)

	err := subscription.Receive(ctx, func(receiverCtx context.Context, message *pubsub.Message) {
		fields := []zap.Field{zap.String("payload", string(message.Data))}

		if err := callback(receiverCtx, message.Data); err != nil {
			c.logger.With(append(fields, zap.Error(err))...).Error("callback returned error")

			message.Nack()

			return
		}

		message.Ack()
	})

	return errors.Wrap(err, "receiving messages")
}

func (c *pubSub) Close() error {
	return c.client.Close()
}

func createResources(ctx context.Context, client *pubsub.Client, subscriptions []Subscription) error {
	for _, subscription := range subscriptions {
		if _, err := createTopic(ctx, client, subscription.Topic); err != nil {
			return err
		}

		subscriptionExists, err := client.Subscription(subscription.Subscription).Exists(ctx)
		if err != nil {
			return errors.Wrap(err, "check subscription exists")
		}

		if !subscriptionExists {
			if err := createSubscription(ctx, client, subscription); err != nil {
				return errors.Wrap(err, "create subscription")
			}
		}
	}

	return nil
}

func createTopic(ctx context.Context, client *pubsub.Client, topicID string) (*pubsub.Topic, error) {
	topic := client.Topic(topicID)

	topicExists, err := topic.Exists(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "check topic exists")
	}

	if topicExists {
		return topic, nil
	}

	topic, err = client.CreateTopic(ctx, topicID)

	return topic, errors.Wrap(err, "create topic")
}

func createSubscription(ctx context.Context, client *pubsub.Client, subscription Subscription) error {
	config := pubsub.SubscriptionConfig{
		Topic: client.Topic(subscription.Topic),
	}

	_, err := client.CreateSubscription(ctx, subscription.Subscription, config)

	return errors.Wrap(err, "create subscription")
}
