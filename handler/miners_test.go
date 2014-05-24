// Copyright (c) 2013-2014 The meeko-collector-pivotal-tracker AUTHORS
//
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package handler

import (
	"testing"

	"github.com/salsita-cider/cider-collector-pivotal-tracker/data"
	. "github.com/smartystreets/goconvey/convey"
)

//------------------------------------//
// pivotaltracker.story_state_changed //
//------------------------------------//

func Test_mineStoryStateChangedEvent(t *testing.T) {
	var change *data.ResourceChange

	Convey("Given a story state update activity", t, func() {
		change = storyStateChangedResourceChange

		Convey("A story state changed event should be generated", func() {
			typ, body := mineStoryStateChangedEvent(change)

			So(typ, ShouldEqual, "pivotaltracker.story_state_changed")
			So(body, ShouldResemble, change)
		})
	})
}
