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

	auth "github.com/abbot/go-http-auth"
	collector "github.com/salsita/cider-abstract-webhook"
)

const (
	statusUnprocessableEntity = 422
	maxBodySize               = int64(10 << 20)
)

func main() {
	collector.ListenAndServe(handlePTActivityHook)
}

func handlePTActivityHook(w http.ResponseWriter, r *http.Request) {
	// Expecting JSON. We could wait for json.Unmarshal to fail, but...
	if ct := r.Header["Content-Type"]; len(ct) != 1 || ct != "application/json" {
		http.Error(w, "Json Expected", http.StatusUnsupportedMediaType)
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

	// Unmarshal the event object.
	var event map[string]interface{}
	err := json.Unmarshal(body, &event)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	// Publish the event.
	kind, ok := event["kind"]
	if !ok {
		http.Error(w, "Kind Field Missing", statusUnprocessableEntity)
		return
	}

	if err := collector.Publish("pivotaltracker."+kind, event); err != nil {
		http.Error(w, "Event Not Published", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
