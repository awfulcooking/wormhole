package wormhole

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"nhooyr.io/websocket"
)

const WebsocketSubprotocol = "wormhole.awful.cooking"

type ControllerHandler struct {
	Pool      *ControllerPool
	ReadLimit int64
	NameGenerator
}

func (h ControllerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
		Subprotocols:   []string{WebsocketSubprotocol},
	})

	if err != nil {
		log.Printf("controller accept error: %v", err)
		return
	}

	ws.SetReadLimit(h.ReadLimit)

	defer ws.Close(websocket.StatusInternalError, "zombie websocket")

	if ws.Subprotocol() != WebsocketSubprotocol {
		ws.Close(websocket.StatusPolicyViolation, "wrong subprotocol")
		return
	}

	host := &WebsocketJSONControllerHost{Conn: ws}
	controller := NewController(r.Context(), host)

	if err := h.NameGenerator.Assign(h.Pool, controller, 3); err != nil {
		log.Printf("couldn't generate unique controller name: %v", err)
		ws.Close(websocket.StatusTryAgainLater, "")
		return
	}

	defer h.Pool.Delete(controller.Name)
	defer controller.Shutdown()

	log.Println(controller, "ready")
	controller.SendWelcome()

	for {
		err = controller.ProcessNext()

		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			log.Println(controller, "host exited normally")
			return
		} else if err != nil {
			log.Printf("%v error (%v): %v", controller, r.RemoteAddr, err.Error())
			return
		}
	}
}

type WebsocketJSONControllerHost struct {
	*websocket.Conn
	sync.Mutex
}

var _ ControllerHost = &WebsocketJSONControllerHost{}

func (h *WebsocketJSONControllerHost) ReadControllerPacket() (ControllerPacket, error) {
	var packet ControllerPacket

	if _, buf, err := h.Read(context.Background()); err != nil {
		return packet, err
	} else if err = json.Unmarshal(buf, &packet); err != nil {
		return packet, err
	} else {
		return packet, nil
	}
}

func (h *WebsocketJSONControllerHost) WriteControllerPacket(packet ControllerPacket) error {
	if buf, err := json.Marshal(packet); err != nil {
		return err
	} else {
		return h.write(buf)
	}
}

func (h *WebsocketJSONControllerHost) write(buf []byte) error {
	h.Lock()
	defer h.Unlock()

	return h.Write(context.Background(), websocket.MessageText, buf)
}

func (h *WebsocketJSONControllerHost) Close() error {
	return h.Conn.Close(websocket.StatusNormalClosure, "closing")
}
