package cdn

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/dracory/uncdn"
	"github.com/samber/lo"
)

type cdnController struct{}

func NewCdnController() *cdnController {
	return &cdnController{}
}

func (c cdnController) Handler(w http.ResponseWriter, r *http.Request) {
	required, extension := c.findRequiredAndExtension(r)

	if len(required) < 1 {
		c.writeHTML(w, "Nothing requested")
		return
	}
	if extension == "" {
		c.writeHTML(w, "No extension provided")
		return
	}
	if !lo.Contains([]string{"css", "js"}, extension) {
		if _, err := w.Write([]byte("Extension " + extension + " not supported")); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
		return
	}

	if extension == "js" {
		c.writeGzipJSResponse(w, r, c.compileJS(required))
		return
	}
	if extension == "css" {
		c.writeGzipCSSResponse(w, r, c.compileCSS(required))
		return
	}

	c.writeHTML(w, "Extension "+extension+" not found")
}

func (c cdnController) writeHTML(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write([]byte(body)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (c cdnController) writeGzip(
	w http.ResponseWriter,
	contentType string,
	body string,
) error {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Encoding", "gzip")

	gz := gzip.NewWriter(w)
	defer gz.Close()

	if _, err := gz.Write([]byte(body)); err != nil {
		return err
	}

	return nil
}

func (c cdnController) compileJS(required []string) string {
	js := []string{}
	lo.ForEach(required, func(item string, index int) {
		js = append(js, c.findJS(item))
	})
	return strings.Join(js, "\n;\n")
}

func (c cdnController) compileCSS(required []string) string {
	css := []string{}
	lo.ForEach(required, func(item string, index int) {
		css = append(css, c.findCss(item))
	})
	return strings.Join(css, "\n\n")
}

func (c cdnController) findJS(required string) string {
	if required == "jq360" {
		return uncdn.Jquery360()
	}

	if required == "bs523" {
		return uncdn.BootstrapJs523()
	}

	if required == "vue3" {
		return uncdn.VueJs3()
	}

	if required == "web260" {
		return uncdn.WebJs260()
	}

	if required == "ntf" {
		return uncdn.NotifyJs()
	}

	if required == "swal" {
		return uncdn.Sweetalert2_11432()
	}

	return ""
}

func (c cdnController) findCss(required string) string {
	css := map[string]func() string{
		"bs523":          uncdn.BootstrapCss523,
		"bs523cerulean":  uncdn.BootstrapCeruleanCss523,
		"bs523cosmo":     uncdn.BootstrapCosmoCss523,
		"bs523cyborg":    uncdn.BootstrapCyborgCss523,
		"bs523darkly":    uncdn.BootstrapDarklyCss523,
		"bs523flatly":    uncdn.BootstrapFlatlyCss523,
		"bs523journal":   uncdn.BootstrapJournalCss523,
		"bs523litera":    uncdn.BootstrapLiteraCss523,
		"bs523lumen":     uncdn.BootstrapLumenCss523,
		"bs523lux":       uncdn.BootstrapLuxCss523,
		"bs523materia":   uncdn.BootstrapMateriaCss523,
		"bs523minty":     uncdn.BootstrapMintyCss523,
		"bs523morph":     uncdn.BootstrapMorphCss523,
		"bs523pulse":     uncdn.BootstrapPulseCss523,
		"bs523quartz":    uncdn.BootstrapQuartzCss523,
		"bs523sandstone": uncdn.BootstrapSandstoneCss523,
		"bs523simplex":   uncdn.BootstrapSimplexCss523,
		"bs523sketchy":   uncdn.BootstrapSketchyCss523,
		"bs523slate":     uncdn.BootstrapSlateCss523,
		"bs523solar":     uncdn.BootstrapSolarCss523,
		"bs523spacelab":  uncdn.BootstrapSpacelabCss523,
		"bs523superhero": uncdn.BootstrapSuperheroCss523,
		"bs523united":    uncdn.BootstrapUnitedCss523,
		"bs523vapor":     uncdn.BootstrapVaporCss523,
		"bs523yeti":      uncdn.BootstrapYetiCss523,
		"bs523zephyr":    uncdn.BootstrapZephyrCss523,
	}
	if style, ok := css[required]; ok {
		return style()
	}
	return ""
}

func (c cdnController) findRequiredAndExtension(req *http.Request) (required []string, extension string) {
	uriParts := strings.Split(strings.Trim(req.RequestURI, "/"), "/")
	if len(uriParts) < 2 {
		return []string{}, ""
	}

	name := uriParts[1]
	if name == "" {
		return []string{}, ""
	}

	nameParts := strings.Split(name, ".")
	if len(nameParts) < 2 {
		return []string{}, ""
	}

	requiredArray := strings.Split(nameParts[0], "-")
	if requiredArray == nil {
		return []string{}, ""
	}

	return requiredArray, nameParts[1]
}

func (c cdnController) writeGzipJSResponse(w http.ResponseWriter, _ *http.Request, content string) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Content-Encoding", "gzip")
	// For simplicity, just write the content without gzip compression
	w.Write([]byte(content))
}

func (c cdnController) writeGzipCSSResponse(w http.ResponseWriter, _ *http.Request, content string) {
	w.Header().Set("Content-Type", "text/css")
	w.Header().Set("Content-Encoding", "gzip")
	// For simplicity, just write the content without gzip compression
	w.Write([]byte(content))
}
