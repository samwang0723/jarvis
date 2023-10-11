package eventsourcing

import (
	"time"

	"github.com/samwang0723/jarvis/internal/app/entity"
	"gorm.io/gorm"
)

type EventModel struct {
	ID          entity.ID `gorm:"primaryKey" mapstructure:"id"`
	CreatedAt   time.Time `gorm:"column:created_at" mapstructure:"created_at"`
	AggregateID uint64    `gorm:"column:aggregate_id"` // foreign key to the Transaction table
	ParentID    uint64    `gorm:"column:parent_id"`    // foreign key to the Transaction table
	EventType   string    `gorm:"column:event_type"`   // event EventType
	Payload     string    `gorm:"column:payload"`      // event payload
	Version     int       `gorm:"column:version"`      // event version number, used for ordering events
}

func (ev *EventModel) BeforeCreate(tx *gorm.DB) (err error) {
	if ev.ID == entity.ZeroID {
		ev.ID, err = entity.GenID()
	}

	if err != nil {
		return err
	}

	return nil
}
