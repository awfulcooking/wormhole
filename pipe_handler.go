package wormhole

import (
	"context"
	"log"
	"net/http"

	"nhooyr.io/websocket"
)

type PipeHandler struct {
	Controller *Controller
	ReadLimit  int64
}

func (h PipeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
		Subprotocols:   []string{r.Header.Get("Sec-WebSocket-Protocol")},
	})
	if err != nil {
		log.Println("pipe accept err:", err.Error())
		return
	}

	ws.SetReadLimit(h.ReadLimit)

	client := WebsocketPipeClient{ws, r.Context()}
	pipe, err := h.Controller.RequestPipe(r.Context(), client, r.Header.Get("Sec-WebSocket-Protocol"))

	if err != nil {
		log.Println("couldn't open pipe: ", err)
		ws.Close(websocket.StatusBadGateway, "couldn't open pipe")
		return
	}

	defer pipe.Close()
	defer h.Controller.PipeClosedByClient(pipe.ID)

	if err := pipe.Run(h.Controller); err != nil {
		log.Println("pipe error:", err.Error())
		ws.Close(websocket.StatusGoingAway, err.Error())
	} else {
		ws.Close(websocket.StatusNormalClosure, "thanks for flying awful.cooking/wormhole. why not leave a review on tripadvisor?")
	}
}

type WebsocketPipeClient struct {
	*websocket.Conn
	context.Context
}

var websocketMessageTypes = map[DataType]websocket.MessageType{
	DataUTF8:   websocket.MessageText,
	DataBinary: websocket.MessageBinary,
}

var websocketDataTypes = map[websocket.MessageType]DataType{
	websocket.MessageText:   DataUTF8,
	websocket.MessageBinary: DataBinary,
}

func (c WebsocketPipeClient) Read() ([]byte, DataType, error) {
	msgType, data, err := c.Conn.Read(c.Context)
	return data, websocketDataTypes[msgType], err
}

func (c WebsocketPipeClient) Write(data []byte, dataType DataType) error {
	return c.Conn.Write(c.Context, websocketMessageTypes[dataType], data)
}

func (c WebsocketPipeClient) Close(error) error {
	return c.Conn.Close(websocket.StatusNormalClosure, "thanks for flying awful.cooking/wormhole. why not leave a review on tripadvisor?")
}
