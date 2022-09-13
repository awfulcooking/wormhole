package wormhole

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bojand/hri"
	"nhooyr.io/websocket"
)

const WebsocketSubprotocol = "awful.cooking/wormhole"

type ControllerHandler struct {
	Pool *ControllerPool
}

func (h ControllerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:   []string{WebsocketSubprotocol},
		OriginPatterns: []string{"*"},
	})

	if err != nil {
		log.Printf("%v", err)
		return
	}

	defer ws.Close(websocket.StatusInternalError, "zombie websocket")

	if ws.Subprotocol() != WebsocketSubprotocol {
		ws.Close(websocket.StatusPolicyViolation, "bad subprotocol")
		return
	}

	// todo: generate human id for controller: https://github.com/bojand/hri
	// todo: store controller in map by id

	host := WebsocketJSONControllerHost{Conn: ws}
	controller := NewController(r.Context(), host)

	slug := hri.Random()
	h.Pool.SetUniq(slug, controller)
	defer h.Pool.Delete(slug)

	controller.SendMeta(ControllerMeta{
		Slug: slug,
	})

	for {
		err = controller.ProcessNext()

		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		} else if err != nil {
			log.Printf("Controller error [%v]: %v", r.RemoteAddr, err.Error())
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
		log.Println("controller websocket packet read error")
		return packet, err
	} else if err = json.Unmarshal(buf, &packet); err != nil {
		log.Println("json unmarshal error")
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
