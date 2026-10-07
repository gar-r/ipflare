package container

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var restartPath = regexp.MustCompile(`^/v[\d.]+/containers/([^/]+)/restart$`)

func newFakeDaemon(t *testing.T, requests *[]string) *httptest.Server {
	t.Helper()
	var mu sync.Mutex
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/_ping" {
			w.Header().Set("Api-Version", "1.44")
			w.WriteHeader(http.StatusOK)
			return
		}
		match := restartPath.FindStringSubmatch(r.URL.Path)
		if match == nil || r.Method != http.MethodPost {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		mu.Lock()
		*requests = append(*requests, match[1])
		mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestRestart(t *testing.T) {
	t.Run("posts restart for each container", func(t *testing.T) {
		var requests []string
		srv := newFakeDaemon(t, &requests)
		t.Setenv("DOCKER_HOST", "tcp://"+srv.Listener.Addr().String())

		Restart([]string{"nginx", "caddy"})

		assert.Equal(t, []string{"nginx", "caddy"}, requests)
	})

	t.Run("no requests when no containers configured", func(t *testing.T) {
		var requests []string
		srv := newFakeDaemon(t, &requests)
		t.Setenv("DOCKER_HOST", "tcp://"+srv.Listener.Addr().String())

		Restart(nil)

		assert.Empty(t, requests)
	})
}
