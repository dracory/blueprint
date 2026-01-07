package home

import (
	"fmt"
	"html"
	"net/http"
	"project/internal/registry"
)

// == CONSTRUCTOR ==============================================================

func NewHomeController(registry registry.RegistryInterface) *homeController {
	return &homeController{
		registry: registry,
	}
}

// == CONTROLLER ===============================================================

type homeController struct {
	registry registry.RegistryInterface
}

// == PUBLIC METHODS ===========================================================

func (controller *homeController) Handler(w http.ResponseWriter, r *http.Request) string {
	appName := "Dracory"
	if controller != nil && controller.registry != nil && controller.registry.GetConfig() != nil {
		if controller.registry.GetConfig().GetAppName() != "" {
			appName = controller.registry.GetConfig().GetAppName()
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	escapedAppName := html.EscapeString(appName)

	return fmt.Sprintf(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>%s</title>
    <style>
      :root {
        --bg1: #0b1220;
        --bg2: #0f1b33;
        --card: rgba(255, 255, 255, 0.06);
        --border: rgba(255, 255, 255, 0.12);
        --text: rgba(255, 255, 255, 0.92);
        --muted: rgba(255, 255, 255, 0.70);
        --primary: #ff2d20;
        --primary2: #ff5a51;
      }
      * { box-sizing: border-box; }
      body {
        margin: 0;
        font-family: ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial;
        color: var(--text);
        background: radial-gradient(1200px circle at 20%% 20%%, #172a52 0%%, transparent 55%%),
                    radial-gradient(900px circle at 85%% 15%%, #2a174f 0%%, transparent 50%%),
                    linear-gradient(180deg, var(--bg1), var(--bg2));
        min-height: 100vh;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 32px 16px;
      }
      .container { width: 100%%; max-width: 980px; }
      .top {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 16px;
        margin-bottom: 18px;
      }
      .brand {
        display: flex;
        align-items: center;
        gap: 10px;
        font-weight: 600;
        letter-spacing: 0.2px;
      }
      .dot {
        width: 10px;
        height: 10px;
        border-radius: 999px;
        background: linear-gradient(135deg, var(--primary), var(--primary2));
        box-shadow: 0 0 18px rgba(255, 45, 32, 0.5);
      }
      .pill {
        font-size: 12px;
        color: var(--muted);
        padding: 6px 10px;
        border: 1px solid var(--border);
        border-radius: 999px;
        background: rgba(255, 255, 255, 0.03);
      }
      .card {
        border: 1px solid var(--border);
        background: var(--card);
        border-radius: 16px;
        overflow: hidden;
        backdrop-filter: blur(10px);
        box-shadow: 0 20px 60px rgba(0, 0, 0, 0.35);
      }
      .card-inner {
        padding: 28px;
      }
      h1 {
        margin: 0 0 10px;
        font-size: 34px;
        line-height: 1.15;
        letter-spacing: -0.5px;
      }
      p {
        margin: 0;
        color: var(--muted);
        line-height: 1.6;
        max-width: 70ch;
      }
      .actions {
        display: flex;
        flex-wrap: wrap;
        gap: 10px;
        margin-top: 18px;
      }
      .btn {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        padding: 10px 14px;
        border-radius: 10px;
        border: 1px solid var(--border);
        color: var(--text);
        text-decoration: none;
        background: rgba(255, 255, 255, 0.04);
      }
      .btn-primary {
        border-color: rgba(255, 45, 32, 0.55);
        background: linear-gradient(135deg, rgba(255, 45, 32, 0.95), rgba(255, 90, 81, 0.90));
        color: #fff;
      }
      .grid {
        display: grid;
        grid-template-columns: repeat(3, minmax(0, 1fr));
        gap: 12px;
        margin-top: 18px;
      }
      .item {
        border: 1px solid var(--border);
        border-radius: 12px;
        padding: 14px;
        background: rgba(255, 255, 255, 0.03);
      }
      .item h3 { margin: 0 0 6px; font-size: 14px; }
      .item p { font-size: 13px; }
      .item a { color: var(--text); }
      .footer {
        margin-top: 14px;
        font-size: 12px;
        color: rgba(255, 255, 255, 0.55);
      }
      @media (max-width: 820px) { .grid { grid-template-columns: 1fr; } }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="top">
        <div class="brand"><span class="dot"></span><span>%s</span></div>
        <div class="pill">Welcome</div>
      </div>

      <div class="card">
        <div class="card-inner">
          <h1>Welcome to %s</h1>
          <p>
            Your application is running. This starter includes routing, controllers, and optional modules
            you can enable as you build.
          </p>

          <div class="actions">
            <a class="btn btn-primary" href="/swagger">API</a>
            <a class="btn" href="/blog">Blog</a>
            <a class="btn" href="/user">Dashboard</a>
          </div>

          <div class="grid">
            <div class="item">
              <h3>Routes</h3>
              <p>Website routes are defined under <code>internal/controllers/website</code>.</p>
            </div>
            <div class="item">
              <h3>Config</h3>
              <p>Use environment variables / config to set <code>AppName</code>, stores, and integrations.</p>
            </div>
            <div class="item">
              <h3>Next</h3>
              <p>Replace this page with your real landing page or enable CMS pages.</p>
            </div>
          </div>

          <div class="footer">Go · Dracory</div>
        </div>
      </div>
    </div>
  </body>
</html>
`, escapedAppName, escapedAppName, escapedAppName)
}
