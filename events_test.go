package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

//------------------------------------//
// pivotaltracker.story_state_changed //
//------------------------------------//

func Test_mineStoryStateChangedEvent(t *testing.T) {
	var a *Activity

	Convey("Given a story state update activity", t, func() {
		a = &Activity{
			Resource: &Resource{
				Kind:      "story",
				Id:        61540514,
				Name:      "New Story",
				StoryType: "feature",
				URL:       "http://www.pivotaltracker.com/story/show/61540514",
			},
			Change: &Change{
				ResourceId:   61540514,
				ResourceKind: "story",
				Type:         "update",
				OriginalValues: map[string]interface{}{
					"current_state": "unscheduled",
					"updated_at":    1385550496000,
					"owned_by_id":   nil,
				},
				NewValues: map[string]interface{}{
					"current_state": "started",
					"updated_at":    1385550486000,
					"owned_by_id":   537831,
				},
			},
		}

		Convey("A story state changed event should be generated", func() {
			So(mineStoryStateChangedEvent(a), ShouldResemble, &Event{
				Type: "pivotaltracker.story_state_changed",
				Body: a,
			})
		})
	})
}
