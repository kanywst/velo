package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Measurement is the database model
type Measurement struct {
	gorm.Model
	Timestamp     time.Time
	DownloadSpeed float64
	UploadSpeed   float64
	Latency       float64
	IPAddress     string
}

// MeasurementDTO is the struct exposed to the frontend
type MeasurementDTO struct {
	ID            uint    `json:"id"`
	Timestamp     string  `json:"timestamp"` // ISO8601 string
	DownloadSpeed float64 `json:"download_speed"`
	UploadSpeed   float64 `json:"upload_speed"`
	Latency       float64 `json:"latency"`
	IPAddress     string  `json:"ip_address"`
}

// ToDTO converts the database model to the frontend DTO
func (m Measurement) ToDTO() MeasurementDTO {
	return MeasurementDTO{
		ID:            m.ID,
		Timestamp:     m.Timestamp.Format(time.RFC3339),
		DownloadSpeed: m.DownloadSpeed,
		UploadSpeed:   m.UploadSpeed,
		Latency:       m.Latency,
		IPAddress:     m.IPAddress,
	}
}