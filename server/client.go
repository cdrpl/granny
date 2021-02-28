package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second    // timeout for writing messages to clients
	pongWait       = 60 * time.Second    // allowed time between reads before client is terminated due to inactivity
	pingPeriod     = (pongWait * 9) / 10 // interval between sending ping messages to clientss
	maxMessageSize = 512                 // maximum message size from clients
)

// Client represents a user's WebSocket connection.
type Client struct {
	id   int64
	conn *websocket.Conn
	send chan []byte // channel used for sending data to the conn.
}

// CreateClient will create and return a Client instance.
func CreateClient(id int64, conn *websocket.Conn) *Client {
	return &Client{
		id:   id,
		conn: conn,
		send: make(chan []byte, 255),
	}
}

// ReadPump will continually read data from the conn and add it to the server.incoming channel.
func (c *Client) ReadPump(server *Server) {
	defer func() {
		server.Unregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				fmt.Fprintf(os.Stderr, "read pump error: %s\n", err)
			}
			break
		}

		if len(data) == 0 {
			fmt.Println("Received empty message from user", c.id)
			return
		}

		// Create the message
		message := Message{client: c, data: data}

		// Send message to incoming
		server.Incoming <- message
	}
}

// WritePump will receive data from the client.send channel and write it to the conn.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod) // ping ticker
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case data, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}

			// the server closed the channel
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// write the data
			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}

			_, err = w.Write(data)
			if err != nil {
				return
			}

			// close the writer
			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C: // ping interval using ticker
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
