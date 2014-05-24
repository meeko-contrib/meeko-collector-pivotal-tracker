// Copyright (c) 2013-2014 The meeko-collector-pivotal-tracker AUTHORS
//
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package data

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
