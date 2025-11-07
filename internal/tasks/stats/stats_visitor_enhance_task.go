package tasks

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"project/internal/types"

	"github.com/dracory/statsstore"
	"github.com/dracory/taskstore"
	"github.com/spf13/cast"

	"github.com/dracory/base/cfmt"

	"github.com/mileusna/useragent"
)

const (
	ipLookupEndpoint = "https://ip2c.org/"
	ipLookupTimeout  = 5 * time.Second
)

var ipLookupHTTPClient = &http.Client{
	Timeout: ipLookupTimeout,
}

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

	country := t.findCountryByIp(ctx, visitor.IpAddress())

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

func (t *statsVisitorEnhanceTask) findCountryByIp(ctx context.Context, ip string) string {
	if ip == "" || ip == "127.0.0.1" {
		return "UN"
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, ipLookupTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCtx, http.MethodGet, ipLookupEndpoint+ip, nil)
	if err != nil {
		log.Printf("Creating geo lookup request failed: %s", err)
		return "ER"
	}

	resp, err := ipLookupHTTPClient.Do(req)

	if err != nil {
		log.Printf("Request Failed: %s", err)
		return "ER" // error
	}

	if resp == nil {
		return "UN"
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Geo lookup returned status code %d", resp.StatusCode)
		return "ER"
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Closing response body failed: %s", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Reading geo lookup response failed: %s", err)
		return "ER"
	}

	parts := strings.Split(string(body), ";")
	if len(parts) > 2 {
		code := strings.TrimSpace(parts[1])
		if code == "" {
			return "UN"
		}
		return code
	}

	return "UN"
}
