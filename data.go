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

type ActivityItem struct {
	Project     Project
	Kind        string
	Message     string
	Highlight   string
	PerformedBy Person     `json:"performed_by"`
	OccurredAt  float64    `json:"occurred_at"`
	Resources   []Resource `json:"primary_resources"`
	Changes     []Change
}

type Activity struct {
	Project     *Project
	Kind        *string
	Message     *string
	Highlight   *string
	PerformedBy *Person  `json:"performed_by" codec:"performed_by"`
	OccurredAt  *float64 `json:"occurred_at"  codec:"occurred_at"`
	Resource    *Resource
	Change      *Change
}

func (item *ActivityItem) Activities() []*Activity {
	as := make([]*Activity, 0)

	for _, ch := range item.Changes {
		for _, r := range item.Resources {
			if ch.ResourceId == r.Id {
				as = append(as, &Activity{
					Project:     &item.Project,
					Kind:        &item.Kind,
					Message:     &item.Message,
					Highlight:   &item.Highlight,
					PerformedBy: &item.PerformedBy,
					OccurredAt:  &item.OccurredAt,
					Resource:    &r,
					Change:      &ch,
				})
			}
		}
	}

	return as
}

type Project struct {
	Id   int
	Name string
}

type Person struct {
	Id       int
	Name     string
	Initials string
}

//------------------------------------------------------------------------------
// Resource
//------------------------------------------------------------------------------

type Resource struct {
	Kind        string
	Id          int    `json:",omitempty"           codec:",omitempty"`
	Number      string `json:",omitempty"           codec:",omitempty"`
	Name        string `json:",omitempty"           codec:",omitempty"`
	Description string `json:",omitempty"           codec:",omitempty"`
	Text        string `json:",omitempty"           codec:",omitempty"`
	StoryType   string `json:"story_type,omitempty" codec:"story_type,omitempty"`
	URL         string
}

func (res *Resource) IsStory() bool {
	return res.Kind == "story"
}

func (res *Resource) MustBeStory() {
	if !res.IsStory() {
		panic("Resource not a story")
	}
}

func (res *Resource) AsStory() *StoryResource {
	return &StoryResource{res}
}

// Resource specialization: Story resource -------------------------------------

type StoryResource struct {
	res *Resource
}

func (res *StoryResource) Id() int {
	return res.res.Id
}

func (res *StoryResource) Name() string {
	return res.res.Name
}

func (res *StoryResource) Type() string {
	return res.res.StoryType
}

func (res *StoryResource) URL() string {
	return res.res.URL
}

//------------------------------------------------------------------------------
// Change
//------------------------------------------------------------------------------

type Change struct {
	ResourceId     int                    `json:"id,-"            codec:"id,-"`
	ResourceKind   string                 `json:"kind,-"          codec:"kind,-"`
	Type           string                 `json:"change_type"     codec:"change_type"`
	OriginalValues map[string]interface{} `json:"original_values" codec:"original_values"`
	NewValues      map[string]interface{} `json:"new_values"      codec:"new_values"`
}

func (ch *Change) AsStoryChange() *StoryChange {
	ch.mustBeChangeOf("story")
	return (*StoryChange)(ch)
}

// Helpers ---------------------------------------------------------------------

func (ch *Change) mustBeChangeOf(kind string) {
	if ch.ResourceKind != kind {
		panic("Not a change of a " + kind)
	}
}

func (ch *Change) hasValueChanged(key string) bool {
	if _, ok := ch.NewValues[key]; ok {
		return true
	} else {
		return false
	}
}

func (ch *Change) getStringValueChange(key string) *StringValueChange {
	original, ok := ch.OriginalValues[key]
	if !ok {
		return nil
	}

	current := ch.NewValues[key]

	return &StringValueChange{
		Original: &StringValue{
			Value:     original.(string),
			OwnerId:   ch.OriginalValues["owned_by_id"].(float64),
			UpdatedAt: ch.OriginalValues["updated_at"].(float64),
		},
		Current: &StringValue{
			Value:     current.(string),
			OwnerId:   ch.NewValues["owned_by_id"].(float64),
			UpdatedAt: ch.NewValues["updated_at"].(float64),
		},
	}
}

// Change specialization: Story change -----------------------------------------

type StoryChange Change

func (ch *StoryChange) HasStateChanged() bool {
	return ((*Change)(ch)).hasValueChanged("current_state")
}

func (ch *StoryChange) StateChange() *StringValueChange {
	return ((*Change)(ch)).getStringValueChange("current_state")
}
