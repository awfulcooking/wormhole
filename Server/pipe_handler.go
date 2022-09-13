package wormhole

import (
	"context"
	"net/http"

	"nhooyr.io/websocket"
)

type PipeHandler struct {
	Controller *Controller
}

func (h PipeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		println("pipe accept err:", err.Error())
		return
	}

	client := WebsocketPipeClient{ws, r.Context()}
	pipe, err := h.Controller.RequestPipe(r.Context(), client, r.Header.Get("Sec-WebSocket-Protocol"))

	if err != nil {
		println("couldn't open pipe: ", err)
		ws.Close(websocket.StatusBadGateway, "couldn't open pipe")
		return
	}

	if err := pipe.Run(h.Controller); err != nil {
		ws.Close(websocket.StatusGoingAway, err.Error())
	} else {
		ws.Close(websocket.StatusNormalClosure, "thanks for flying awful.cooking/wormhole. why not leave a review on tripadvisor?")
	}
}

type WebsocketPipeClient struct {
	*websocket.Conn
	context.Context
}

var websocketMessageTypes = map[PipeDataType]websocket.MessageType{
	DataUTF8:   websocket.MessageText,
	DataBinary: websocket.MessageBinary,
}

var websocketDataTypes = map[websocket.MessageType]PipeDataType{
	websocket.MessageText:   DataUTF8,
	websocket.MessageBinary: DataBinary,
}

func (c WebsocketPipeClient) Read() ([]byte, PipeDataType, error) {
	msgType, data, err := c.Conn.Read(c.Context)
	return data, websocketDataTypes[msgType], err
}

func (c WebsocketPipeClient) Write(data []byte, dataType PipeDataType) error {
	return c.Conn.Write(c.Context, websocketMessageTypes[dataType], data)
}

func (c WebsocketPipeClient) Close(error) error {
	return c.Conn.Close(websocket.StatusNormalClosure, "thanks for flying awful.cooking/wormhole. why not leave a review on tripadvisor?")
}
