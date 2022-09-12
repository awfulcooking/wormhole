package wormhole

import (
	"context"
	"encoding/json"
	"net/http"

	"nhooyr.io/websocket"
)

const WebsocketSubprotocol = "awful.cooking/wormhole"

type ControllerHandler struct {
	Printf func(string, ...any)
}

func (h ControllerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{WebsocketSubprotocol},
	})

	if err != nil {
		h.Printf("%v", err)
		return
	}

	defer ws.Close(websocket.StatusInternalError, "outlived handler")

	if ws.Subprotocol() != WebsocketSubprotocol {
		ws.Close(websocket.StatusPolicyViolation, "bad subprotocol")
		return
	}

	// todo: generate human id for controller: https://github.com/bojand/hri

	host := WebsocketJSONControllerHost{Conn: ws}
	controller := NewController(r.Context(), host)

	for {
		err = controller.ProcessNext()

		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		} else if err != nil {
			h.Printf("Controller error [%v]: %v", r.RemoteAddr, err)
			return
		}
	}
}

type WebsocketJSONControllerHost struct {
	*websocket.Conn
}

var _ ControllerHost = WebsocketJSONControllerHost{}

func (h WebsocketJSONControllerHost) ReadControllerPacket() (ControllerPacket, error) {
	var packet ControllerPacket

	if _, buf, err := h.Read(context.Background()); err != nil {
		return packet, err
	} else if err = json.Unmarshal(buf, &packet); err != nil {
		return packet, err
	} else {
		return packet, nil
	}
}

func (h WebsocketJSONControllerHost) WriteControllerPacket(packet ControllerPacket) error {
	if buf, err := json.Marshal(packet); err != nil {
		return err
	} else if err := h.Write(context.Background(), websocket.MessageText, buf); err != nil {
		return err
	}
	return nil
}

func (h WebsocketJSONControllerHost) Close() error {
	return h.Close()
}
