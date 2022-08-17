package config

import (
	"interview/pkg/queues"
)

const (
	GreetingKey = "greeting_key"

	Topic        = "queue_topic"
	Subscription = "topic_subscription"
)

var Subscriptions = []queues.Subscription{
	{
		Topic:        Topic,
		Subscription: Subscription,
	},
}

type Config struct {
	UnixSocket string `conf:"env:API_UNIX_SOCKET,default:/tmp/api.sock"`
	ProjectID  string `conf:"env:PROJECT_ID,default:interview"`
	Redis      redis
}

type redis struct {
	URL string `conf:"env:REDIS_URL"`
}
