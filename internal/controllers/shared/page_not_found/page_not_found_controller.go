package page_not_found

import (
	"net/http"
)

// == CONSTRUCTOR =============================================================

func PageNotFoundController() *pageNotFoundController {
	return &pageNotFoundController{}
}

// == CONTROLLER ==============================================================

type pageNotFoundController struct{}

// PUBLIC METHODS =============================================================

func (controller *pageNotFoundController) Handler(w http.ResponseWriter, r *http.Request) string {
	w.WriteHeader(http.StatusNotFound)

	// Create a beautiful 404 page
	return controller.pageHTML()
}

func (controller *pageNotFoundController) pageHTML() string {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>404 - Page Not Found</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.0/font/bootstrap-icons.css" rel="stylesheet">
    <style>
        body {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        }
        .error-container {
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            border-radius: 20px;
            padding: 3rem;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
            text-align: center;
            max-width: 600px;
            width: 90%;
        }
        .error-code {
            font-size: 6rem;
            font-weight: 900;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            margin-bottom: 1rem;
            line-height: 1;
        }
        .error-title {
            font-size: 2.5rem;
            font-weight: 700;
            color: #2d3748;
            margin-bottom: 1rem;
        }
        .error-message {
            font-size: 1.1rem;
            color: #718096;
            margin-bottom: 2rem;
            line-height: 1.6;
        }
        .error-icon {
            font-size: 4rem;
            color: #667eea;
            margin-bottom: 1.5rem;
        }
        .home-button {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border: none;
            border-radius: 50px;
            padding: 1rem 2rem;
            font-size: 1.1rem;
            font-weight: 600;
            color: white;
            text-decoration: none;
            transition: all 0.3s ease;
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
        }
        .home-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 20px rgba(102, 126, 234, 0.3);
            color: white;
            text-decoration: none;
        }
        .back-button {
            background: transparent;
            border: 2px solid #667eea;
            border-radius: 50px;
            padding: 0.8rem 1.5rem;
            font-size: 1rem;
            font-weight: 500;
            color: #667eea;
            text-decoration: none;
            transition: all 0.3s ease;
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
            margin-right: 1rem;
        }
        .back-button:hover {
            background: #667eea;
            color: white;
            text-decoration: none;
        }
        .footer-links {
            margin-top: 2rem;
            padding-top: 2rem;
            border-top: 1px solid #e2e8f0;
        }
        .footer-links a {
            color: #718096;
            text-decoration: none;
            margin: 0 1rem;
            font-size: 0.9rem;
            transition: color 0.3s ease;
        }
        .footer-links a:hover {
            color: #667eea;
        }
    </style>
</head>
<body>
    <div class="error-container">
        <div class="error-icon">
            <i class="bi bi-compass"></i>
        </div>
        <div class="error-code">404</div>
        <h1 class="error-title">Oops! Page Not Found</h1>
        <p class="error-message">
            It looks like you've ventured into uncharted territory. The page you're looking for seems to have sailed away!
            Don't worry, even the best explorers get lost sometimes.
        </p>
        <div class="button-group">
            <a href="javascript:history.back()" class="back-button">
                <i class="bi bi-arrow-left"></i>
                Go Back
            </a>
            <a href="/" class="home-button">
                <i class="bi bi-house-door"></i>
                Back to Home
            </a>
        </div>
        <div class="footer-links">
            <a href="/">Home</a>
            <a href="/contact">Contact</a>
            <a href="/help">Help</a>
        </div>
    </div>
</body>
</html>`

	return html
}
