// Copyright (c) 2013-2014 The meeko-collector-pivotal-tracker AUTHORS
//
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/meeko-contrib/meeko-collector-pivotal-tracker/data"
)

const (
	statusUnprocessableEntity = 422
	maxBodySize               = int64(10 << 20)
)

type PTWebhookHandler struct {
	Forward func(eventType string, eventObject interface{}) error
}

func (handler *PTWebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Expecting JSON. We could wait for json.Unmarshal to fail, but...
	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		http.Error(w, "Json Content Type Expected", http.StatusUnsupportedMediaType)
		return
	}

	// Read the request body.
	bodyReader := http.MaxBytesReader(w, r.Body, maxBodySize)
	defer bodyReader.Close()

	body, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		http.Error(w, "Request Payload Too Large", http.StatusRequestEntityTooLarge)
		return
	}

	// Unmarshal the activity item.
	var item data.ActivityItem
	err = json.Unmarshal(body, &item)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	// Mine all the events that can be generated from the received item.
	for _, change := range item.ExtractChanges() {
		for _, mineFunc := range mineFuncs {
			eventType, eventBody := mineFunc(change)
			if eventType == "" {
				continue
			}

			if err := handler.Forward(eventType, eventBody); err != nil {
				http.Error(w, "Event Not Published", http.StatusInternalServerError)
				// This is a critical error, panic.
				panic(err)
			}
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
