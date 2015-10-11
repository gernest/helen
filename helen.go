package helen

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/labstack/echo"
)

// File is an interface for a static file.
type File interface {
	os.FileInfo
	io.ReadSeeker
}

// Static is a http handle, that serves static files.
type Static struct {
	*http.ServeMux

	fs http.FileSystem

	//	prefix is the url prefix that needs to be removed inorder to match the file name
	//
	// Example if URL is /static/hello.world and prefix is /hello/
	// The the file name will resolve to hello.world.
	// NOTE There is no need to put the opening slash / it will be added automatically
	prefix string

	// noCache diables caching when set to true, the defaul value is false.
	noCache bool

	// chains is a collection of chainable manipulation of the static files
	chains []FileChain

	// middlewares is a slice of alice compatible middlewares.
	middlewares []alice.Constructor
	c           *cache
}

// NewStatic returns a new Static that serves dir conents as static files.
// The optional argument noCache if set true caching is diabled, default is false.
func NewStatic(dir string, noCache ...bool) *Static {
	var hasCache bool
	if len(noCache) > 0 {
		hasCache = noCache[0]
	}
	m := &Static{
		ServeMux: http.NewServeMux(),
		fs:       http.Dir(dir),
		c:        newCache(),
		noCache:  hasCache,
	}

	// register / route so we can match all the routes passing here.
	return m.Register("/")
}

// Register registers a pattern that will match for the static files. This is visible on the *Static
// level, meaning that after the route has successfull matched on the server, this can be used to filter
// out what routes you want to operate on.
func (m *Static) Register(pattern string) *Static {
	m.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		m.getAlice().ThenFunc(m.handleStatic).ServeHTTP(w, r)
	})
	return m
}

// Bind register the static handler at path patern on server. Server can be of the following types
//	=> *http.ServeMux from http package
// 	=> *mux.Router from gorilla mux package
// 	=> *echo.Echo from labstack echo package
// pattern should not contain any regular expressions or fancy stuffs. There is no need for
// opening or closing slash.
//
// So, the following will all register for the /static/ route
// 	=> /static
// 	=> /static/
// 	=> static/
// 	=> static
func (m *Static) Bind(pattern string, server interface{}) {
	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}
	if !strings.HasSuffix(pattern, "/") {
		pattern = pattern + "/"
	}
	switch server.(type) {
	case *mux.Router: // support gorilla mux
		router := server.(*mux.Router)
		router.PathPrefix(pattern).Handler(m)
	case *echo.Echo: // support labstack echo
		router := server.(*echo.Echo)
		router.Get(pattern+"*", m)
	case *http.ServeMux: // support http.ServeMux
		router := server.(*http.ServeMux)
		router.Handle(pattern, m)
	}
}

// getAlice returns alice Chain with all middleares loaded
func (m *Static) getAlice() alice.Chain {
	return alice.New(m.middlewares...)
}

// handleStatic is a http.HandlerFunc helper that serves static files.
func (m *Static) handleStatic(w http.ResponseWriter, r *http.Request) {
	name := filepath.Clean(r.URL.Path)
	m.serveFile(w, r, name)
}

// StripPrefix sets p for stripping from the URL Path.
func (m *Static) StripPrefix(p string) *Static {
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	m.prefix = p
	return m
}

func (m *Static) serveFile(w http.ResponseWriter, r *http.Request, name string) {
	if m.prefix != "" {
		name = strings.TrimPrefix(name, m.prefix)
	}
	if !m.noCache {
		if f, ok := m.c.GetOk(name); ok {
			http.ServeContent(w, r, f.Name(), f.ModTime(), f)
			return
		}

	}
	f, err := m.fs.Open(name)
	if err != nil {
		m.serveError(w, r, err)
		return
	}
	defer f.Close()
	nf, err := NewFromHTTPFile(f)
	if err != nil {
		m.serveError(w, r, err)
		return
	}

	supported := []string{".css", ".js"}
	var found bool
	xt := filepath.Ext(nf.Name())
	for _, v := range supported {
		if v == xt {
			found = true
			break
		}
	}
	defer func() {
		if !m.noCache && found {
			m.c.Set(name, nf)
		}
	}()

	if nf.IsDir() {
		http.NotFound(w, r)
		return
	}
	if len(m.chains) > 0 {
		if found {
			for k := range m.chains {
				nf = m.chains[k](nf)
			}
		}

	}
	http.ServeContent(w, r, nf.Name(), nf.ModTime(), nf)
}

func (m *Static) serveError(w http.ResponseWriter, r *http.Request, err error) {
	if os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}
}

// Use registers middleware
func (m *Static) Use(middleware ...alice.Constructor) *Static {
	m.middlewares = append(m.middlewares, middleware...)
	return m
}

type cache struct {
	mu    sync.RWMutex
	files map[string]File
}

func newCache() *cache {
	return &cache{
		mu:    sync.RWMutex{},
		files: make(map[string]File),
	}
}

func (c *cache) GetOk(key string) (file File, ok bool) {
	file, ok = c.files[key]
	return
}

func (c *cache) Set(key string, val File) {
	c.mu.RLock()
	c.files[key] = val
	c.mu.RUnlock()
}
