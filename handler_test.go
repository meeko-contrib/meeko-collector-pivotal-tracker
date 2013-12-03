package main

import (
	"bytes"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/salsita-cider/cider-collector-pivotal-tracker/data"
	. "github.com/smartystreets/goconvey/convey"
)

//------------------------------------//
// pivotaltracker.story_state_changed //
//------------------------------------//

// A dummy raw Pivotal Tracker Activity Item that represents a story state change.
var rawStoryStateChangedItem = []byte(`
	{
	  "message": "Pepa Novak started this feature",
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
		"name": "Pepa Novak",
		"initials": "PN",
		"kind": "person",
		"id": 537831
	  }
	}
	`)

// A ResourceChange that represents the resource change contained in
// the activity item specified above.
var storyStateChangedResourceChange *data.ResourceChange

func init() {
	var (
		project = data.Project{
			Id:   959942,
			Name: "Workflow Test Project",
		}
		kind        = "story_update_activity"
		message     = "Pepa Novak started this feature"
		highlight   = "started"
		performedBy = data.Person{
			Id:       537831,
			Name:     "Pepa Novak",
			Initials: "PN",
		}
		occurredAt = float64(1385550496000)
		resource   = data.Resource{
			Kind:      "story",
			Id:        61540514,
			Name:      "New Story",
			StoryType: "feature",
			URL:       "http://www.pivotaltracker.com/story/show/61540514",
		}
		change = data.Change{
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

	storyStateChangedResourceChange = &data.ResourceChange{
		Project:     &project,
		Kind:        &kind,
		Message:     &message,
		Highlight:   &highlight,
		PerformedBy: &performedBy,
		OccurredAt:  &occurredAt,
		Resource:    &resource,
		Change:      &change,
	}
}

func TestPTWebhookHandler_HandleStoryStateChangedActivityItem(t *testing.T) {
	var (
		forwardedType string
		forwardedBody interface{}
	)

	handler := &PTWebhookHandler{
		Forward: func(eventType string, eventBody interface{}) error {
			forwardedType = eventType
			forwardedBody = eventBody
			return nil
		},
	}

	Convey("Receiving a Pivotal Tracker Activity Webhook", t, func() {
		req, err := http.NewRequest("POST", "http://example.com",
			bytes.NewReader(rawStoryStateChangedItem))
		if err != nil {
			t.Fatal(err)
		}

		req.Header = http.Header{
			"Content-Type": {"application/json"},
		}

		rw := httptest.NewRecorder()

		handler.ServeHTTP(rw, req)

		if rw.Code != http.StatusAccepted {
			t.Fatalf("Unexpected status code returned: expected %d, received %d %s",
				http.StatusAccepted, rw.Code, rw.Body.String())
		}

		Convey("A story_state_changed event with the correct payload should be emitted", func() {
			So(forwardedType, ShouldEqual, "pivotaltracker.story_state_changed")
			So(forwardedBody, ShouldResemble, storyStateChangedResourceChange)
		})
	})
}

// Helpers ---------------------------------------------------------------------

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getRandomPayload() map[string]interface{} {
	m := make(map[string]interface{}, rand.Intn(20))
	for i := 0; i < len(m); i++ {
		k := strconv.Itoa(rand.Int())
		v := rand.Int()
		m[k] = v
	}
	return m
}
