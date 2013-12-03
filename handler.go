/*
   Copyright (C) 2013  Salsita s.r.o.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program. If not, see {http://www.gnu.org/licenses/}.
*/

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/salsita-cider/cider-collector-pivotal-tracker/data"
	receiver "github.com/salsita-cider/cider-webhook-receiver"
)

const (
	statusUnprocessableEntity = 422
	maxBodySize               = int64(10 << 20)
)

type PTWebhookHandler struct {
	Forward receiver.ForwardFunc
}

func (handler *PTWebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Expecting JSON. We could wait for json.Unmarshal to fail, but...
	if ct := r.Header["Content-Type"]; len(ct) != 1 || ct[0] != "application/json" {
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

	w.WriteHeader(http.StatusNoContent)
}
