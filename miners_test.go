package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//------------------------------------//
// pivotaltracker.story_state_changed //
//------------------------------------//

func Test_mineStoryStateChangedEvent(t *testing.T) {
	var change *ResourceChange

	Convey("Given a story state update activity", t, func() {
		change = storyStateChangedResourceChange

		Convey("A story state changed event should be generated", func() {
			typ, body := mineStoryStateChangedEvent(change)

			So(typ, ShouldEqual, "pivotaltracker.story_state_changed")
			So(body, ShouldResemble, change)
		})
	})
}
