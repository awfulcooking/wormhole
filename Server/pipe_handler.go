package wormhole

import (
	"net/http"

	"nhooyr.io/websocket"
)

type PipeHandler struct {
	Controller *Controller
}

func (h PipeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
	if err != nil {
		println("err: ", err)
		return
	}

	client := websocket.NetConn(r.Context(), ws, websocket.MessageBinary)
	pipe, err := h.Controller.RequestPipe(r.Context(), client)

	if err != nil {
		println("couldn't open pipe: ", err)
		ws.Close(websocket.StatusBadGateway, "couldn't open pipe")
		return
	}

	if err := pipe.Run(); err != nil {
		ws.Close(websocket.StatusGoingAway, err.Error())
	} else {
		ws.Close(websocket.StatusNormalClosure, "thanks for flying awful.cooking/wormhole. why not leave a review on tripadvisor?")
	}
}

type WebsocketPipeClient struct {
	*websocket.Conn
}

func (c WebsocketPipeClient) Close() error {
	c.Conn.Close(websocket.StatusNormalClosure, "thanks for flying awful.cooking/wormhole. why not leave a review on tripadvisor?")
}
