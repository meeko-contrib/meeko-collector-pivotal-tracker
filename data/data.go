// Copyright (c) 2013-2014 The meeko-collector-pivotal-tracker AUTHORS
//
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package data

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

type ResourceChange struct {
	Project     *Project
	Kind        *string
	Message     *string
	Highlight   *string
	PerformedBy *Person  `json:"performed_by" codec:"performed_by"`
	OccurredAt  *float64 `json:"occurred_at"  codec:"occurred_at"`
	Resource    *Resource
	Change      *Change
}

func (item *ActivityItem) ExtractChanges() []*ResourceChange {
	changes := make([]*ResourceChange, 0)

	for _, ch := range item.Changes {
		for _, r := range item.Resources {
			if ch.ResourceId == r.Id {
				changes = append(changes, &ResourceChange{
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

	return changes
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

func (ch *Change) IsStoryChange() bool {
	return ch.isChangeOf("story")
}

func (ch *Change) AsStoryChange() *StoryChange {
	ch.mustBeChangeOf("story")
	return (*StoryChange)(ch)
}

// Helpers ---------------------------------------------------------------------

func (ch *Change) isChangeOf(kind string) bool {
	return ch.ResourceKind == kind
}

func (ch *Change) mustBeChangeOf(kind string) {
	if !ch.isChangeOf(kind) {
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
			UpdatedAt: ch.OriginalValues["updated_at"].(float64),
		},
		Current: &StringValue{
			Value:     current.(string),
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
