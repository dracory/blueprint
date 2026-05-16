package taskadmin

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/taskstore"
)

type admin struct {
	externalLayout func(options LayoutOptions) string
	logger         *slog.Logger
	taskStore      taskstore.StoreInterface
}

func (a *admin) Router(w http.ResponseWriter, r *http.Request) {
	path := req.GetStringTrimmedOr(r, "path", "home")

	if path == "" {
		path = PathHome
	}

	ctx := context.WithValue(r.Context(), keyEndpoint, r.URL.Path)

	routeFunc := a.getRoute(path)
	routeFunc(w, r.WithContext(ctx))
}

func (a *admin) getRoute(route string) func(http.ResponseWriter, *http.Request) {
	routes := map[string]func(w http.ResponseWriter, r *http.Request){
		PathHome: newPageQueueManager(a).Handler,
		// PathQueueManager: a.pageQueueManager,
	}
	// log.Println(route)
	if val, ok := routes[route]; ok {
		return val
	}

	return routes[PathHome]
}

// func (a *admin) pageHome(w http.ResponseWriter, r *http.Request) {
// 	endpoint := r.Context().Value(keyEndpoint).(string)
// 	content := NewQueueManagerController().AnyIndex(w, r)
// 	a.layout(w, LayoutOptions{
// 		Title:   "Home",
// 		Content: content,
// 	})
// }

func (a *admin) layout(w http.ResponseWriter, options LayoutOptions) {
	content := ""

	options.ScriptURLs = []string{
		cdn.Sweetalert2_11(),
	}

	if a.externalLayout != nil {
		content = a.externalLayout(options)
	} else {
		content = Webpage(options.Title, options.Content).ToHTML()
	}

	w.Write([]byte(content))
}

func Webpage(title, content string) *hb.HtmlWebpage {
	faviconImgCms := `data:image/x-icon;base64,AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAmzKzAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABEQEAAQERAAEAAQABAAEAAQABAQEBEQABAAEREQEAAAERARARAREAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD//wAA//8AAP//AAD//wAA//8AAP//AAD//wAAi6MAALu7AAC6owAAuC8AAIkjAAD//wAA//8AAP//AAD//wAA`
	app := ""
	webpage := hb.NewWebpage()
	webpage.SetTitle(title)
	webpage.SetFavicon(faviconImgCms)

	webpage.AddStyleURLs([]string{
		cdn.BootstrapCss_5_2_3(),
	})
	webpage.AddScriptURLs([]string{
		cdn.BootstrapJs_5_2_3(),
		cdn.Jquery_3_6_4(),
		cdn.VueJs_3(),
		cdn.Sweetalert2_10(),
	})
	webpage.AddScripts([]string{
		app,
	})
	webpage.AddStyle(`html,body{height:100%;font-family: Ubuntu, sans-serif;}`)
	webpage.AddStyle(`body {
		font-family: "Nunito", sans-serif;
		font-size: 0.9rem;
		font-weight: 400;
		line-height: 1.6;
		color: #212529;
		text-align: left;
		background-color: #f8fafc;
	}
	.form-select {
		display: block;
		width: 100%;
		padding: .375rem 2.25rem .375rem .75rem;
		font-size: 1rem;
		font-weight: 400;
		line-height: 1.5;
		color: #212529;
		background-color: #fff;
		background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16'%3e%3cpath fill='none' stroke='%23343a40' stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M2 5l6 6 6-6'/%3e%3c/svg%3e");
		background-repeat: no-repeat;
		background-position: right .75rem center;
		background-size: 16px 12px;
		border: 1px solid #ced4da;
		border-radius: .25rem;
		-webkit-appearance: none;
		-moz-appearance: none;
		appearance: none;
	}`)
	webpage.AddChild(hb.NewHTML(content))
	return webpage
}
