package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type dreamMetadata struct {
	Note     int8     `binding:"range(0,4)" json:"note,omitempty"`
	Lucid    bool     `json:"lucid,omitempty"`
	Peoples  []string `json:"peoples,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	ToReview *bool    `json:"toReview,omitempty"`
	TextNote *string  `json:"textNote,omitempty"`
}

type techMetadata struct {
	LastChange *time.Time `binding:"required" json:"lastChange,omitempty" bson:"lastChange"`
}

type dream struct {
	Id            string        `json:"id,omitempty"`
	Name          string        `binding:"required" json:"name,omitempty"`
	Text          string        `binding:"required" json:"text,omitempty"`
	DreamMetadata dreamMetadata `json:"dreamMetadata,omitempty" bson:"dreamMetadata"`
}

type DreamDay struct {
	Id           string       `json:"id,omitempty"`
	Date         *time.Time   `json:"date,omitempty"`
	TechMetadata techMetadata `binding:"required" json:"techMetadata,omitempty" bson:"techMetadata"`
	UserId       string       `json:"userId,omitempty" bson:"userId"`
	Dreams       []dream      `binding:"required" json:"dreams"`
}

func (dreamDay *DreamDay) HandleDefault() {
	if dreamDay.Date == nil {
		now := time.Now()
		date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
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
