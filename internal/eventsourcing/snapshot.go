package eventsourcing

type Snapshot struct {
	BaseEvent

	Data string `gorm:"column:data"` // JSON-encoded snapshot data
}
