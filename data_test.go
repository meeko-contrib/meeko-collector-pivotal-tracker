package main

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestActivityItem_ExtractChanges(t *testing.T) {
	var (
		rawItem []byte
		item    ActivityItem
	)

	Convey("Given a raw Pivotal Tracker Activity Item", t, func() {
		rawItem = rawStoryStateChangedItem

		Convey("A correct list of resource changes should be extracted", func() {
			err := json.Unmarshal(rawItem, &item)
			if err != nil {
				t.Error(err)
				return
			}

			So(item.ExtractChanges(), ShouldResemble, []*ResourceChange{
				storyStateChangedResourceChange,
			})
		})
	})
}
