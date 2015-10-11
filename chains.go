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
	minifier = minify.New()

	mediaStylesheet = "text/css"
	mediaJavascript = "text/javascript"
)

// MinifyChain minifies the given src. Files supported are javascript and stylesheet.
func MinifyChain(src *StaticFile) *StaticFile {
	buf := &bytes.Buffer{}
	switch filepath.Ext(src.Name()) {
	case ".css":
		css.Minify(minifier, mediaStylesheet, buf, src)
	case ".js":
		js.Minify(minifier, mediaJavascript, buf, src)
	default:
		io.Copy(buf, src)
	}
	src.SetReadSeeker(bytes.NewReader(buf.Bytes()))
	return src
}
