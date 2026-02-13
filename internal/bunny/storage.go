package bunny

import (
	"time"
	"github.com/ShadowSeick/tufte-style-static-site-generator/internal/pkg/http"
)

type HTTPError struct {
	Code int `json:"HttpCode"`
	Message string `json:"Message"`
}

func (httpErr HTTPError) NotFound() bool {
	return httpErr.Code == http.StatusNotFound.Code()
}

type File struct {
	GUID string `json:"Guid"`
	StorageZoneName string `json:"StorageZoneName"`
	Path string `json:"Path"`
	ObjecName string `json:"ObjectName"`
	Length uint `json:"Length"`
	Checksum *string `json:"Checksum,omitempty"`
	ContentType string `json:"ContentType,omitempty"`
	ReplicatedZones *string `json:"ReplicatedZones,omitempty"`
	LastChanged time.Time `json:"LastChanged"`
	IsDirectory bool `json:"IsDirectory"`
	ServerID int `json:"ServerId"`
	UserID string `json:"UserId"`
	DateCreated time.Time `json:"DateCreated"`
	StorageZoneID int `json:"StorageZoneId"`
}
