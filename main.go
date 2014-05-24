// Copyright (c) 2013-2014 The meeko-collector-pivotal-tracker AUTHORS
//
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"github.com/meeko-contrib/meeko-collector-pivotal-tracker/handler"

	"github.com/meeko-contrib/go-meeko-webhook-receiver/receiver"
	"github.com/meeko/go-meeko/agent"
)

func main() {
	var (
		logger    = agent.Logging()
		publisher = agent.PubSub()
	)
	receiver.ListenAndServe(&handler.PTWebhookHandler{
		func(eventType string, eventObject interface{}) error {
			logger.Infof("Forwarding %s", eventType)
			return publisher.Publish(eventType, eventObject)
		},
	})
}
