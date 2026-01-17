package render

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

// ------------------------------------------------
// TEST SESSION (mocked for render tests)
// ------------------------------------------------
var session *scs.SessionManager

func init() {
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false
}

// ------------------------------------------------
// TEST RESPONSE WRITER
// ------------------------------------------------
type myWriter struct{}

func (mw *myWriter) Header() http.Header {
	return make(http.Header)
}

func (mw *myWriter) WriteHeader(statusCode int) {}

func (mw *myWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
