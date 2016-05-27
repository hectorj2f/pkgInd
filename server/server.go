package server

import (
	"bufio"
	"io"
	"net"
	"strings"

	"github.com/hectorj2f/pkgInd/indexer"
	"github.com/hectorj2f/pkgInd/log"
	"github.com/hectorj2f/pkgInd/transport"
)

type PkgIndServer struct {
	pkgMgmt *indexer.PackageManager

	addchan chan Client
	rmchan  chan net.Conn

	eCh     chan error
	reqChan chan Message
}

type Client struct {
	conn net.Conn
	ch   chan<- string
}

type Message struct {
	conn    net.Conn
	content string
}

// NewServer() creates a server initializing all its attributes
func NewServer() *PkgIndServer {
	server := &PkgIndServer{
		pkgMgmt: indexer.NewPackageManager(),
		addchan: make(chan Client),
		rmchan:  make(chan net.Conn),
		eCh:     make(chan error),
		reqChan: make(chan Message),
	}
	return server
}

// Start receives and address and scheme to create a listener and to start waiting
// for incoming connections.
// Returns an error if something went wrong
func (s *PkgIndServer) Start(addr, scheme string) error {
	ln, err := transport.NewListener(addr, scheme)
	if err != nil {
		log.Logger().Errorf("Error: %v\n", err)
		return err
	}

	// This function will handle the registration of new Clients
	go s.handleClientsRegistry(s.addchan, s.rmchan)

	// We handle here the message requests
	go func(reqChan <-chan Message, eCh <-chan error) {
		for {
			select {
			case req := <-reqChan:
				go s.handleMessageRequest(req.conn, req.content)
			case err := <-eCh:
				log.Logger().Error(err)
			}
		}
	}(s.reqChan, s.eCh)

	log.Logger().Info("Starting to accept connections ...")
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Logger().Errorf("Unable to accept a connection %v:", err)
			continue
		}

		go s.serveClient(conn, s.reqChan, s.eCh, s.addchan, s.rmchan)
	}

	return nil
}

// serveClient receives a connection and waits to process all the messages from this
// connection. Within this function, the messages are processed and clients are registered.
func (s *PkgIndServer) serveClient(conn net.Conn, ch chan Message, eCh chan error, addchan chan<- Client, rmchan chan<- net.Conn) {
	defer conn.Close()

	clientCh := make(chan string)
	addchan <- Client{conn, clientCh}

	for {
		rawMsg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err != io.EOF {
				eCh <- err
				return
			}
		}
		if len(rawMsg) == 0 {
			break
		}

		req := Message{
			conn:    conn,
			content: rawMsg,
		}
		ch <- req
	}
	rmchan <- conn
}

// handleClientsRegistry receives two channels to manages the clients connection/disconnection
func (s *PkgIndServer) handleClientsRegistry(addchan <-chan Client, rmchan <-chan net.Conn) {
	clients := make(map[net.Conn]chan<- string)

	for {
		select {
		case client := <-addchan:
			log.Logger().Infof("New client: %v\n", client.conn)
			clients[client.conn] = client.ch
		case conn := <-rmchan:
			log.Logger().Infof("Client disconnects: %v\n", conn)
			delete(clients, conn)
		}
	}
}

// handleMessageRequest receives a connection and the input messate and performs
// any required operation from: message validation til index,remove or query ops.
func (s *PkgIndServer) handleMessageRequest(conn net.Conn, rawMsg string) {
	stringMsg := strings.TrimRight(rawMsg, "\n")
	if err := transport.ValidateMessage(stringMsg); err != nil {
		log.Logger().Error(err)
		s.sendMessageResponse(conn, transport.ERROR)
		return
	}

	msg, err := transport.ExtractMessage(stringMsg)
	if err != nil {
		log.Logger().Error(err)
		s.sendMessageResponse(conn, transport.ERROR)
		return
	}

	switch msg.Command {
	case transport.IndexCmd:
		log.Logger().Debugf("Received INDEX message: '%v'", msg)
		resp, err := s.pkgMgmt.Index(msg)
		if err != nil {
			log.Logger().Errorf("unable to index the package %s: %v", stringMsg, err)
		}
		s.sendMessageResponse(conn, resp)

	case transport.RemoveCmd:
		log.Logger().Debugf("Received REMOVE message: '%v'", msg)
		resp, err := s.pkgMgmt.Remove(msg)
		if err != nil {
			log.Logger().Errorf("unable to remove the package %s: %v", stringMsg, err)
		}
		s.sendMessageResponse(conn, resp)

	case transport.QueryCmd:
		log.Logger().Debugf("Received QUERY message: '%v'", msg)
		resp, err := s.pkgMgmt.Query(msg)
		if err != nil {
			log.Logger().Errorf("unable to query the package %s: %v", stringMsg, err)
		}
		s.sendMessageResponse(conn, resp)
	}
}

// sendMessageResponse receives a connection and response message code to be
// transferred to the clients
func (s *PkgIndServer) sendMessageResponse(conn net.Conn, msg transport.MessageResponseCode) error {
	log.Logger().Debugf("Sending response message '%s' ...", msg)
	_, err := conn.Write([]byte(msg + "\n"))
	if err != nil {
		return err
	}
	return nil
}
