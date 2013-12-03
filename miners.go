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
	"github.com/salsita-cider/cider-collector-pivotal-tracker/data"
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
