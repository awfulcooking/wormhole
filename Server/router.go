package wormhole

import (
	"net/http"
	"strings"
)

type Router struct {
	Printf      func(string, ...interface{})
	Controllers *ControllerPool
}

func NewRouter(printf func(string, ...any)) *Router {
	return &Router{
		Printf:      printf,
		Controllers: NewControllerPool(),
	}
}

func (h *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		ControllerHandler{Printf: h.Printf}.ServeHTTP(w, r)
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
		}.ServeHTTP(w, r)
	}
}
