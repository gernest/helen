package helen

import (
	"bytes"
	"io"
	"path/filepath"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

var (
	mediaStylesheet = "text/css"
	mediaJavascript = "text/javascript"

	minifier = func() *minify.M {
		m := minify.New()
		m.AddFunc(mediaStylesheet, css.Minify)
		m.AddFunc(mediaJavascript, js.Minify)
		return m
	}()
)

// MinifyChain minifies the given src. Files supported are javascript and stylesheet.
func MinifyChain(src *StaticFile) *StaticFile {
	buf := &bytes.Buffer{}
	switch filepath.Ext(src.Name()) {
	case ".css":
		minifier.Minify(mediaStylesheet, buf, src)
	case ".js":
		minifier.Minify(mediaJavascript, buf, src)
	default:
		io.Copy(buf, src)
	}
	src.SetReadSeeker(bytes.NewReader(buf.Bytes()))
	return src
}
