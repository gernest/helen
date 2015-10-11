package helen

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/labstack/echo"
)

func TestStatic(t *testing.T) {
	static := NewStatic("fixtures")

	sample := []struct {
		name string
		code int
	}{
		{"/static/css/hello.css", http.StatusOK},
		{"/static/js/hello.js", http.StatusOK},
		{"/static/js/notfound.js", http.StatusNotFound},
		{"/static/", http.StatusNotFound},
	}
	for _, data := range sample {
		req, _ := http.NewRequest("GET", data.name, nil)
		w := httptest.NewRecorder()
		static.ServeHTTP(w, req)
		if w.Code != data.code {
			t.Errorf("expected %d got %d", data.code, w.Code)
		}
	}
}

func TestBind(t *testing.T) {
	static := NewStatic("fixtures")

	sample := []struct {
		name string
		code int
	}{
		{"/static/css/hello.css", http.StatusOK},
		{"/static/js/hello.js", http.StatusOK},
		{"/static/js/notfound.js", http.StatusNotFound},
		{"/static/", http.StatusNotFound},
	}
	for _, data := range sample {
		req, _ := http.NewRequest("GET", data.name, nil)
		w := httptest.NewRecorder()

		//
		// gorilla mux
		//
		gorilla := mux.NewRouter()
		static.Bind("/static", gorilla)
		gorilla.ServeHTTP(w, req)

		if w.Code != data.code {
			t.Errorf("expected %d got %d", data.code, w.Code)
		}

		//
		// labstack echo
		//
		e := echo.New()
		static.Bind("static/", e)
		req, _ = http.NewRequest("GET", data.name, nil)
		w = httptest.NewRecorder()
		e.ServeHTTP(w, req)
		if w.Code != data.code {
			t.Errorf("expected %d got %d", data.code, w.Code)
		}

		//
		// http.ServeMux
		//
		sm := http.NewServeMux()
		static.Bind("/static/", sm)
		req, _ = http.NewRequest("GET", data.name, nil)
		w = httptest.NewRecorder()
		sm.ServeHTTP(w, req)
		if w.Code != data.code {
			t.Errorf("expected %d got %d", data.code, w.Code)
		}

	}
	static.StripPrefix("mwanza")
	req, _ := http.NewRequest("GET", "/mwanza/static/css/hello.css", nil)
	w := httptest.NewRecorder()
	static.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expcted %d got %d", http.StatusOK, w.Code)
	}

}

func BenchmarkServeCached(b *testing.B) {
	static := NewStatic("fixtures")

	for n := 0; n < b.N; n++ {
		req, _ := http.NewRequest("GET", "/static/css/hello.css", nil)
		w := httptest.NewRecorder()
		static.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			b.Errorf("expected %d got %d", http.StatusOK, w.Code)
		}
	}
}

func BenchmarkServeNoCached(b *testing.B) {
	static := NewStatic("fixtures", true)
	static.noCache = true

	for n := 0; n < b.N; n++ {
		req, _ := http.NewRequest("GET", "/static/css/hello.css", nil)
		w := httptest.NewRecorder()
		static.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			b.Errorf("expected %d got %d", http.StatusOK, w.Code)
		}
	}
}
