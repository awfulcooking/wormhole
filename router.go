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

	if r.Header.Get("Upgrade") == "websocket" {
		log.Printf("[%d] WebSocket to %s", id, r.URL)
		defer log.Printf("[%d] WebSocket finished (%s)", id, r.URL)
	} else {
		log.Printf("[%d] Serving %s", id, r.URL)
		h.ServeStatic(w, r)
		return
	}

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

func (h *Router) ServeStatic(w http.ResponseWriter, r *http.Request) {
	if fs := h.Config.StaticFS; fs == nil {
		w.WriteHeader(404)
		w.Write(nil)
	} else {
		handler := http.FileServer(http.FS(fs))
		handler.ServeHTTP(w, r)
	}
}
