package admin

import (
	"github.com/dracory/cdn"
)

// uiLayout creates the main layout HTML
func uiLayout(title string, content string) string {
	html := `
<!DOCTYPE html>
<html>
    <head>
        <title>` + title + ` - Media Manager</title>
        <link href="` + cdn.BootstrapIconsCss_1_10_2() + `" rel="stylesheet" type="text/css" />
		<link href="` + cdn.BootstrapCss_5_2_3() + `" rel="stylesheet" type="text/css" />
        <script src="` + cdn.Jquery_3_6_4() + `"></script>
        <script src="` + cdn.BootstrapJs_5_2_3() + `"></script>
		<script src="` + cdn.Notify_0_4_2() + `"></script>
        <style>
            html,body{
                padding-top: 40px;
            }
        </style>
    </head>
    <body>
        <!-- Fixed navbar -->
        <nav class="navbar navbar-expand-lg bg-light fixed-top"  data-bs-theme="dark">
            <div class="container">
				<a class="navbar-brand" href="#">
					Media Manager
				</a>
				<button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarTogglerDemo01" aria-controls="navbarTogglerDemo01" aria-expanded="false" aria-label="Toggle navigation">
					<span class="navbar-toggler-icon"></span>
				</button>
            </div>
        </nav>
        <div class="container">` + content + ` </div>
    </body>
</html>
	`

	return html
}
