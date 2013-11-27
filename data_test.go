package main

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestActivityItem_Activities(t *testing.T) {
	var (
		rawItem []byte
		item    ActivityItem
	)

	Convey("Given a raw Pivotal Tracker Activity Item", t, func() {
		rawItem = []byte(`
		{
		  "message": "Ondrej Kupka started this feature",
		  "occurred_at": 1385550496000,
		  "primary_resources": [
			{
			  "name": "New Story",
			  "story_type": "feature",
			  "url": "http://www.pivotaltracker.com/story/show/61540514",
			  "kind": "story",
			  "id": 61540514
			}
		  ],
		  "changes": [
			{
			  "name": "New Story",
			  "new_values": {
				"updated_at": 1385550496000,
				"current_state": "started",
				"owned_by_id": 537831
			  },
			  "story_type": "feature",
			  "original_values": {
				"updated_at": 1385550486000,
				"current_state": "unscheduled",
				"owned_by_id": null
			  },
			  "change_type": "update",
			  "kind": "story",
			  "id": 61540514
			}
		  ],
		  "highlight": "started",
		  "project": {
			"name": "Workflow Test Project",
			"kind": "project",
			"id": 959942
		  },
		  "kind": "story_update_activity",
		  "project_version": 4,
		  "guid": "959942_4",
		  "performed_by": {
			"name": "Ondrej Kupka",
			"initials": "OK",
			"kind": "person",
			"id": 537831
		  }
		}
		`)

		Convey("A list of activities should be extracted", func() {
			err := json.Unmarshal(rawItem, &item)
			if err != nil {
				t.Error(err)
				return
			}

			var (
				project = Project{
					Id:   959942,
					Name: "Workflow Test Project",
				}
				kind        = "story_update_activity"
				message     = "Ondrej Kupka started this feature"
				highlight   = "started"
				performedBy = Person{
					Id:       537831,
					Name:     "Ondrej Kupka",
					Initials: "OK",
				}
				occurredAt = float64(1385550496000)
				resource   = Resource{
					Kind:      "story",
					Id:        61540514,
					Name:      "New Story",
					StoryType: "feature",
					URL:       "http://www.pivotaltracker.com/story/show/61540514",
				}
				change = Change{
					ResourceId:   61540514,
					ResourceKind: "story",
					Type:         "update",
					OriginalValues: map[string]interface{}{
						"current_state": "unscheduled",
						"updated_at":    float64(1385550486000),
						"owned_by_id":   nil,
					},
					NewValues: map[string]interface{}{
						"current_state": "started",
						"updated_at":    float64(1385550496000),
						"owned_by_id":   float64(537831),
					},
				}
			)

			So(item.Activities(), ShouldResemble, []*Activity{
				{
					Project:     &project,
					Kind:        &kind,
					Message:     &message,
					Highlight:   &highlight,
					PerformedBy: &performedBy,
					OccurredAt:  &occurredAt,
					Resource:    &resource,
					Change:      &change,
				},
			})
		})
	})
}
