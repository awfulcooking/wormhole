package wormhole

import (
	"net/http"
	"strings"
)

type Router struct {
	Config      ServerConfig
	Controllers *ControllerPool
}

func NewRouter(cfg ServerConfig) *Router {
	return &Router{
		Config:      cfg,
		Controllers: NewControllerPool(),
	}
}

func (h *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		ControllerHandler{
			Pool:          h.Controllers,
			ReadLimit:     h.Config.WebsocketReadLimit,
			NameGenerator: h.Config.NameGenerator,
		}.ServeHTTP(w, r)
	} else {
		name := strings.TrimPrefix(r.URL.Path, "/")
		controller := h.Controllers.Get(name)

		if controller == nil {
			w.WriteHeader(404)
			w.Write(nil)
			return
		}

		PipeHandler{
			Controller: controller,
			ReadLimit:  h.Config.WebsocketReadLimit,
		}.ServeHTTP(w, r)
	}
}
