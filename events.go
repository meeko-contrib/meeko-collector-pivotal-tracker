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

type Event struct {
	Type string
	Body interface{}
}

type mineFunc func(*Activity) *Event

var mineFuncs = [...]mineFunc{
	mineStoryStateChangedEvent,
}

//------------------------------------//
// pivotaltracker.story_state_changed //
//------------------------------------//

func mineStoryStateChangedEvent(activity *Activity) *Event {
	if !activity.Resource.IsStory() {
		return nil
	}

	ch := activity.Change.AsStoryChange()

	if ch.Type != "update" {
		return nil
	}

	if !ch.HasStateChanged() {
		return nil
	}

	return &Event{
		Type: "pivotaltracker.story_state_changed",
		Body: activity,
	}
}
