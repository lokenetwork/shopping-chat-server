package chat

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func init() {
	//connect_db()
}

// Chat server.
type Server struct {
	buyer_pattern   string
	shoper_pattern   string
	messages  []*CommandMessage

	buyerClients   map[uint32]*BuyerClient
	shoperClients   map[uint32]*ShoperClient

	addBuyerClientCh     chan *BuyerClient
	addShoperClientCh     chan *ShoperClient

	delBuyerClientCh     chan *BuyerClient
	delShoperClientCh     chan *ShoperClient

	sendBuyerClientAllCh chan *CommandMessage
	sendShoperClientAllCh chan *CommandMessage

	doneCh    chan bool
	errCh     chan error
}

// Create new chat server.
func NewServer(buyer_pattern string,shoper_pattern string) *Server {
	messages := []*CommandMessage{}
	buyer_clients := make(map[uint32]*BuyerClient)
	shoper_clients := make(map[uint32]*ShoperClient)

	addBuyerClientCh := make(chan *BuyerClient)
	addShoperClientCh := make(chan *ShoperClient)

	delBuyerClientCh := make(chan *BuyerClient)
	delShoprClientCh := make(chan *ShoperClient)

	sendBuyerClientAllCh := make(chan *CommandMessage)
	sendShoperClientAllCh := make(chan *CommandMessage)

	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		buyer_pattern,
		shoper_pattern,
		messages,

		buyer_clients,
		shoper_clients,

		addBuyerClientCh,
		addShoperClientCh,

		delBuyerClientCh,
		delShoprClientCh,

		sendBuyerClientAllCh,
		sendShoperClientAllCh,

		doneCh,
		errCh,
	}
}

func (s *Server) AddBuyerClient(c *BuyerClient) {
	s.addBuyerClientCh <- c
}

func (s *Server) AddShoperClient(c *ShoperClient) {
	s.addShoperClientCh <- c
}

func (s *Server) DelBuyerClient(c *BuyerClient) {
	s.delBuyerClientCh <- c
}

func (s *Server) DelShoperClient(c *ShoperClient) {
	s.delShoperClientCh <- c
}

func (s *Server) SendBuyerClientAll(msg *CommandMessage) {
	s.sendBuyerClientAllCh <- msg
}

func (s *Server) SendShoperClientAll(msg *CommandMessage) {
	s.sendShoperClientAllCh <- msg
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}


func (s *Server) sendAll(msg *CommandMessage) {
	for _, c := range s.buyerClients {
		c.Write(msg)
	}
}


// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {

	log.Println("Listening server...")

	// websocket handler
	onBuyerConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			log.Println("ws.Close()")

			if err != nil {
				s.errCh <- err
			}
		}()
		connect_success,buyer_client := NewBuyerClient(ws, s)
		if connect_success {
			s.AddBuyerClient(buyer_client)
			buyer_client.Listen()
		}
	}

	onShoperConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			log.Println("ws.Close()")

			if err != nil {
				s.errCh <- err
			}
		}()
		connect_success,shoper_client := NewShoperClient(ws, s)
		if connect_success {
			s.AddShoperClient(shoper_client)
			shoper_client.Listen()
		}
	}

	http.Handle(s.buyer_pattern, websocket.Handler(onBuyerConnected))

	http.Handle(s.shoper_pattern, websocket.Handler(onShoperConnected))


	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-s.addBuyerClientCh:
			log.Println("Added new buyer client")

			s.buyerClients[c.websocketId] = c
			log.Println("Now", len(s.buyerClients), "buyer clients connected.")

		case c := <-s.addShoperClientCh:
			log.Println("Added new shoper client")

			var shoperClientIndex uint32
			shoperClientIndex = uint32(len(s.shoperClients))

			s.shoperClients[shoperClientIndex] = c
			log.Println("Now", len(s.shoperClients), "shop clients connected.")

		// del a client
		case c := <-s.delBuyerClientCh:
			log.Println("Delete buyer client")
			delete(s.buyerClients, c.websocketId)

		case c := <-s.delShoperClientCh:
			log.Println("Delete shoper client")
			delete(s.shoperClients, c.websocketId)
		// broadcast message for all clients
		case msg := <-s.sendBuyerClientAllCh:
			log.Println("Send all:", msg)

			s.messages = append(s.messages, msg)
			s.sendAll(msg)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
