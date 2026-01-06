package backend

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/showwin/speedtest-go/speedtest"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"velo/backend/models"
)

type VeloApp struct {
	ctx  context.Context
	DB   *gorm.DB
	quit chan struct{}
}

func NewVeloApp() *VeloApp {
	return &VeloApp{}
}

func (a *VeloApp) Startup(ctx context.Context) {
	a.ctx = ctx
	a.quit = make(chan struct{})

	// Init DB
	var err error
	a.DB, err = gorm.Open("sqlite3", "velo_data.db")
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to connect database: %v", err)
	} else {
		a.DB.AutoMigrate(&models.Measurement{})
		// Start Scheduler
		go a.startScheduler()
	}
}

func (a *VeloApp) Shutdown(ctx context.Context) {
	if a.quit != nil {
		close(a.quit)
	}
	if a.DB != nil {
		a.DB.Close()
	}
}

func (a *VeloApp) startScheduler() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			runtime.LogInfo(a.ctx, "Auto-starting measurement...")
			a.RunMeasurement()
		case <-a.quit:
			return
		}
	}
}

func (a *VeloApp) RunMeasurement() *models.MeasurementDTO {
	runtime.LogInfo(a.ctx, "Starting measurement...")

	// Basic Speedtest Logic
	var speedTestClient = speedtest.New()
	serverList, _ := speedTestClient.FetchServers()
	targets, _ := serverList.FindServer([]int{})

	if len(targets) == 0 {
		runtime.LogError(a.ctx, "No servers found")
		return nil
	}

	s := targets[0]
	s.PingTest(func(latency time.Duration) {})
	s.DownloadTest()
	s.UploadTest()

	m := &models.Measurement{
		Timestamp:     time.Now(),
		DownloadSpeed: float64(s.DLSpeed) * 8 / 1000000, // bps to Mbps
		UploadSpeed:   float64(s.ULSpeed) * 8 / 1000000,   // bps to Mbps
		Latency:       float64(s.Latency.Milliseconds()),
	}

	a.DB.Create(m)

	dto := m.ToDTO()
	// Emit event to frontend to update graph
	runtime.EventsEmit(a.ctx, "measurement_complete", dto)

	return &dto
}

func (a *VeloApp) GetHistory() []models.MeasurementDTO {
	var history []models.Measurement
	var dtos []models.MeasurementDTO

	if a.DB != nil {
		a.DB.Order("timestamp asc").Find(&history)
		for _, m := range history {
			dtos = append(dtos, m.ToDTO())
		}
	}
	return dtos
}
