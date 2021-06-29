package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type dreamMetadata struct {
	Note     int8     `valid:"range(0,4)" json:"note,omitempty"`
	Lucid    bool     `json:"lucid,omitempty"`
	Peoples  []string `json:"peoples,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	ToReview *bool    `json:"toReview,omitempty"`
	TextNote *string  `json:"textNote,omitempty"`
}

type techMetadata struct {
	LastChange *time.Time `json:"lastChange,omitempty" bson:"lastChange"`
}

type dream struct {
	Id            string        `json:"id,omitempty"`
	Name          string        `valid:"required" json:"name,omitempty"`
	Text          string        `valid:"required" json:"text,omitempty"`
	DreamMetadata dreamMetadata `json:"dreamMetadata,omitempty" bson:"dreamMetadata"`
}

type DreamDay struct {
	Id           string       `valid:"required" json:"id,omitempty"`
	Date         *time.Time   `json:"date,omitempty"`
	TechMetadata techMetadata `json:"techMetadata,omitempty" bson:"techMetadata"`
	UserId       string       `json:"userId,omitempty" bson:"userId"`
	Dreams       []dream      `json:"dreams,omitempty"`
}

func (dreamDay *DreamDay) HandleDefault() {
	if dreamDay.Date == nil {
		date := time.Now()
		dreamDay.Date = &date
	}

	if dreamDay.Id == "" {
		dreamDay.Id = uuid.NewV4().String()
	}

	for i, dream := range dreamDay.Dreams {
		if dream.Id == "" {
			dreamDay.Dreams[i].Id = uuid.NewV4().String()
		}
	}
}
