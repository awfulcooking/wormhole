package wormhole

import (
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type Router struct {
	Config       ServerConfig
	Controllers  *ControllerPool
	requestCount atomic.Uint64
}

func NewRouter(cfg ServerConfig) *Router {
	return &Router{
		Config:      cfg,
		Controllers: NewControllerPool(),
	}
}

func (h *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := h.requestCount.Add(1)

	log.Printf("[%d] Request to %s", id, r.URL)
	defer log.Printf("[%d] Request finished (%s)", id, r.URL)

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
