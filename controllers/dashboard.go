package controllers

import (
	"net/http"
    "github.com/oskarszura/smarthome/controllers/utils"
	"github.com/oskarszura/gowebserver/router"
	"github.com/oskarszura/gowebserver/session"
)

func CtrDashboard(w http.ResponseWriter, r *http.Request, opt router.UrlOptions, sm session.ISessionManager) {
    utils.RenderTemplate(w, r, "dashboard", sm)
}
