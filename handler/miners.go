// Copyright (c) 2013-2014 The meeko-collector-pivotal-tracker AUTHORS
//
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package handler

import (
	"github.com/meeko-contrib/meeko-collector-pivotal-tracker/data"
)

type mineFunc func(resChange *data.ResourceChange) (typ string, body interface{})

var mineFuncs = [...]mineFunc{
	mineStoryStateChangedEvent,
}

//------------------------------------//
// pivotaltracker.story_state_changed //
//------------------------------------//

func mineStoryStateChangedEvent(resChange *data.ResourceChange) (typ string, body interface{}) {
	if !resChange.Resource.IsStory() {
		return "", nil
	}

	ch := resChange.Change.AsStoryChange()

	if ch.Type != "update" {
		return "", nil
	}

	if !ch.HasStateChanged() {
		return "", nil
	}

	return "pivotaltracker.story_state_changed", resChange
}
