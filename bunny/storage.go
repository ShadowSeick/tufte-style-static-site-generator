package bunny

import (
	"time"
)

type File struct {
	GUID string `json:"Guid"`
	StorageZoneName string `json:"StorageZoneName"`
	Path string `json:"Path"`
	ObjecName string `json:"ObjectName"`
	Length uint `json:"Length"`
	Checksum string `json:"Checksum,omitempty"`
	ContentType string `json:"ContentType,omitempty"`
	ReplicatedZones string `json:"ReplicatedZones,omitempty"`
	LastChanged time.Time `json:"LastChanged"`
	IsDirectory bool `json:"IsDirectory"`
	ServerID int `json:"ServerId"`
	UserID string `json:"UserId"`
	DateCreated time.Time `json:"DateCreated"`
	StorageZoneID int `json:"StorageZoneId"`
}
