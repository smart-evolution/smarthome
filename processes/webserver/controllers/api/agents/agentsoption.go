package agents

import (
	"github.com/coda-it/gowebserver/router"
	"github.com/coda-it/gowebserver/session"
	"github.com/coda-it/gowebserver/store"
	"github.com/smart-evolution/shapi/processes/webserver/handlers"
	"net/http"
)

// CtrAgentsOptions - options handler
func (c *Controller) CtrAgentsOptions(w http.ResponseWriter, r *http.Request, opt router.URLOptions, sm session.ISessionManager, s store.IStore) {
	handlers.CorsHeaders(w, r)

	switch r.Method {
	case "OPTIONS":
		return
	}
}