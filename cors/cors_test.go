package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	. "github.com/smartystreets/goconvey/convey"
)

func test(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hello world"))
}

func TestUnitCors(t *testing.T) {
	r := mux.NewRouter()

	r.Path("/test").Methods("GET", "OPTIONS").Handler(Middleware(http.HandlerFunc(test)))

	Convey("test OPTIONS request", t, func() {
		req := httptest.NewRequest("OPTIONS", "http://localhost:3000/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		testCorsHeadersPresent(w)

		So(w.Header().Get("Access-Control-Max-Age"), ShouldEqual, "86400")
		So(w.Header().Get("Cache-Control"), ShouldNotEqual, "no-cache")
		So(w.Body.String(), ShouldBeEmpty)
	})

	Convey("test GET request", t, func() {
		req := httptest.NewRequest("GET", "http://localhost:3000/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		testCorsHeadersPresent(w)

		So(w.Header().Get("Access-Control-Max-Age"), ShouldNotEqual, "86400")
		So(w.Header().Get("Cache-Control"), ShouldEqual, "no-cache")
		So(w.Body.String(), ShouldEqual, "hello world")
	})

}

func testCorsHeadersPresent(w http.ResponseWriter) {
	So(w.Header().Get("Access-Control-Allow-Origin"), ShouldEqual, "*")
	So(w.Header().Get("Access-Control-Allow-Credentials"), ShouldEqual, "true")
	So(w.Header().Get("Access-Control-Allow-Headers"), ShouldEqual, "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, X-Requested-With")
	So(w.Header().Get("Access-Control-Allow-Methods"), ShouldEqual, "POST, OPTIONS, GET, PUT")
}
