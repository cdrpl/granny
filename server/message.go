package main

// Message channel
const (
	// Channel for when a user want to join a room.
	JoinRoom = iota

	// Used to check if a channel is valid (n < End)
	End
)

// Message represents a socket message and links the client(sender) to the data.
type Message struct {
	channel int
	data    []byte
	client  *Client
}

// Create a new message.
func newMessage(data []byte, c *Client) Message {
	return Message{
		channel: int(data[0]),
		data:    data[1:],
		client:  c,
	}
}
