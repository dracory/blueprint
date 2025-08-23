package tasks

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"project/internal/types"
	"strings"

	"github.com/gouniverse/statsstore"
	"github.com/gouniverse/taskstore"
	"github.com/spf13/cast"

	"github.com/mingrammer/cfmt"

	"github.com/mileusna/useragent"
)

// statsVisitorEnhanceTask enhances the visitor stats with the country
//
// =================================================================
// Example:
//
// go run main.go task stats-visitor-enhance
//
// =================================================================
type statsVisitorEnhanceTask struct {
	taskstore.TaskHandlerBase
	app types.AppInterface
}

// == CONSTRUCTOR =============================================================

func NewStatsVisitorEnhanceTask(app types.AppInterface) *statsVisitorEnhanceTask {
	if app == nil {
		log.Fatal("app is nil")
	}
	return &statsVisitorEnhanceTask{app: app}
}

// == IMPLEMENTATION ==========================================================

// var _ jobsshared.TaskInterface = (*statsVisitorEnhanceTask)(nil) // verify it extends the task interface
var _ taskstore.TaskHandlerInterface = (*statsVisitorEnhanceTask)(nil) // verify it extends the task interface

// == PUBLIC METHODS ==========================================================

func (t *statsVisitorEnhanceTask) Enqueue() (taskstore.QueueInterface, error) {
	if t.app == nil || t.app.GetTaskStore() == nil {
		return nil, errors.New("task store is nil")
	}
	return t.app.GetTaskStore().TaskEnqueueByAlias(t.Alias(), map[string]interface{}{})
}

func (t *statsVisitorEnhanceTask) Alias() string {
	return "StatsVisitorEnhanceTask"
}

func (t *statsVisitorEnhanceTask) Title() string {
	return "Stats Visitor Enhance"
}

func (t *statsVisitorEnhanceTask) Description() string {
	return "Enhances the visitor stats by adding the country"
}

func (t *statsVisitorEnhanceTask) Handle() bool {
	if t.app == nil || t.app.GetStatsStore() == nil {
		t.LogError("Task StatsVisitorEnhance. Store is nil")
		return false
	}

	ctx := context.Background()
	unprocessedEntries, err := t.app.GetStatsStore().VisitorList(ctx, statsstore.VisitorQueryOptions{
		Country: "empty",
		Limit:   10,
	})

	if err != nil {
		t.LogError("Task StatsVisitorEnhance. Error: " + err.Error())
		return false
	}

	if len(unprocessedEntries) < 1 {
		t.LogInfo("Task StatsVisitorEnhance. No entries to process")
		return true
	}

	t.LogInfo("Task StatsVisitorEnhance. Found: " + cast.ToString(len(unprocessedEntries)) + " entries to process")

	for i := 0; i < len(unprocessedEntries); i++ {
		entry := unprocessedEntries[i]
		t.processVisitor(ctx, entry)
	}

	return false
}

// == PRIVATE METHODS =========================================================

func (t *statsVisitorEnhanceTask) processVisitor(ctx context.Context, visitor statsstore.VisitorInterface) bool {
	if t.app == nil || t.app.GetStatsStore() == nil {
		t.LogError("Task StatsVisitorEnhance. Store is nil")
		return false
	}
	ua := useragent.Parse(visitor.UserAgent())
	userOs := ua.OS
	userOsVersion := ua.OSVersion
	userDevice := ua.Device
	userBrowser := ua.Name
	userBrowserVersion := ua.Version

	userDeviceType := ""

	if ua.Mobile {
		userDeviceType = "mobile"
	}
	if ua.Tablet {
		userDeviceType = "tablet"
	}
	if ua.Desktop {
		userDeviceType = "desktop"
	}
	if ua.Bot {
		userDeviceType = "bot"
	}

	country := t.findCountryByIp(visitor.IpAddress())

	visitor.SetCountry(country)
	visitor.SetUserBrowser(userBrowser)
	visitor.SetUserBrowserVersion(userBrowserVersion)
	visitor.SetUserDevice(userDevice)
	visitor.SetUserDeviceType(userDeviceType)
	visitor.SetUserOs(userOs)
	visitor.SetUserOsVersion(userOsVersion)

	errUpdated := t.app.GetStatsStore().VisitorUpdate(ctx, visitor)

	if errUpdated != nil {
		cfmt.Errorln(errUpdated.Error())
	}

	return false
}

func (t *statsVisitorEnhanceTask) findCountryByIp(ip string) string {
	if ip == "" || ip == "127.0.0.1" {
		return "UN"
	}

	resp, err := http.Get("https://ip2c.org/" + ip)

	if err != nil {
		log.Printf("Request Failed: %s", err)
		return "ER" // error
	}

	if resp == nil {
		return "UN"
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return "ER"
	}
	// Log the request body
	bodyString := string(body)
	cfmt.Infoln(bodyString)
	parts := strings.Split(bodyString, ";")
	if len(parts) > 2 {
		if parts[1] == "" {
			return "UN"
		}
		return parts[1]
	}

	return "UN"
}
