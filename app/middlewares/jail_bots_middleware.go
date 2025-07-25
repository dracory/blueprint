package middlewares

import (
	"log/slog"
	"net/http"
	"project/config"
	"slices"
	"strings"
	"time"

	"github.com/dracory/base/req"
	"github.com/dracory/rtr"
	"github.com/gouniverse/responses"
	"github.com/samber/lo"
)

func JailBotsMiddleware(config JailBotsConfig) rtr.MiddlewareInterface {
	jb := new(jailBotsMiddleware)
	jb.exclude = config.Exclude

	m := rtr.NewMiddleware().
		SetName(jb.Name()).
		SetHandler(jb.Handler)

	return m
}

type JailBotsConfig struct {
	Exclude []string
}

type jailBotsMiddleware struct {
	exclude []string
}

func (j *jailBotsMiddleware) Name() string {
	return "Jail Bots Middleware"
}

func (m *jailBotsMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		ip := req.IP(r)

		if m.isJailed(ip) {
			w.WriteHeader(http.StatusForbidden)
			responses.HTMLResponse(w, r, "malicious access not allowed (jb)")
			return
		}

		jailable, reason := m.isJailable(uri)

		if jailable {
			m.jail(ip)

			config.Logger.Info("Jailed bot from "+ip+" for 5 minutes",
				slog.String("reason", reason),
				slog.String("uri", uri),
				slog.String("ip", ip),
				slog.String("useragent", r.UserAgent()),
			)

			w.WriteHeader(http.StatusForbidden)
			responses.HTMLResponse(w, r, "malicious access not allowed (jb)")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (j *jailBotsMiddleware) isJailed(ip string) bool {
	return config.CacheMemory.Has("jail:" + ip)
}

func (j *jailBotsMiddleware) jail(ip string) {
	config.CacheMemory.Set("jail:"+ip, ip, 5*time.Minute)
}

func (m *jailBotsMiddleware) isJailable(uri string) (jailable bool, reason string) {
	startsWithList := m.startsWithBlacklistedUriList()

	for i := 0; i < len(startsWithList); i++ {
		if strings.HasPrefix(uri, startsWithList[i]) {
			return true, "starts with " + startsWithList[i]
		}
	}

	containsList := m.containsBlacklistedUriList()

	for i := 0; i < len(containsList); i++ {
		if strings.Contains(uri, containsList[i]) {
			return true, "contains " + containsList[i]
		}
	}

	return false, ""
}

// containsBlacklistedUriList returns a list of strings
// which if they are found anywhere in the uri
// clearly indicate that there is a malicious bot/user
// trying to access them.
func (j *jailBotsMiddleware) containsBlacklistedUriList() []string {
	stopList := []string{
		"print(",
		"${print",
		".aws",
		".DS_Store",
		".env",
		".env.example",
		".git",
		".php",
		".vscode",
		".well-known/ALFA_DATA",
		".well-known/alfacgiapi",
		".well-known/cgialfa",
		"_ignition/health-check",
		"ALFA_DATA",
		"alfacgiapi",
		"search?folderIds=0",
		"aws/credentials",
		"backup",
		"backup/license.txt",
		"bc",
		"bk",
		"blog/license.txt",
		"bin",
		"cgialfa",
		"cloud-config.yml",
		"components/com_",
		"content/sitetree",
		"config.json",
		"cgi-bin",
		"credentials",
		"db",
		"ecp/Current/exporttool/microsoft.exchange.ediscovery.exporttool.application",
		"js/mage/cookies.js",
		"META-INF",
		"/main",
		"/new",
		"/old",
		"phpinfo",
		"server-status",
		"Telerik.Web.UI.WebResource.axd",
		"shop/license.txt",
		"sites/all/libraries/plupload/examples/upload.php",
		"simpla",
		"telescope/requests",
		"tmp/license.txt",
		"v2/_catalog",
		"wordpress",
		"wp",
		"www/license.txt",
	}

	// Check if we have any exclusion rules?

	if len(j.exclude) > 0 { // Check if exclude list is not empty
		stopList = lo.Filter(stopList, func(item string, index int) bool {
			return !slices.Contains(j.exclude, item)
		})
	}

	return stopList
}

// startsWithBlacklistedUriList returns a list of strings
// which if they are found at the start of the uri
// clearly indicate that there is a malicious bot/user
// trying to access them.
func (j jailBotsMiddleware) startsWithBlacklistedUriList() []string {
	return []string{
		"/content/sitetree",
		"/backup",
		"/bc",
		"/bk",
		"/main",
		"/new",
		"/old",
		"/tmp/",
		"/wordpress",
		"/wp",
		"/www",
	}
}
