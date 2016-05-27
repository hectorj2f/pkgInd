package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/hectorj2f/pkgind/log"
	"github.com/hectorj2f/pkgind/transport"
)

type Client struct {
	listenIP    string
	listenPort  int
	enableDebug bool
}

func (c Client) sendMessage(msg string) (transport.MessageResponseCode, error) {
	conn, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", c.listenIP, c.listenPort))
	defer conn.Close()

	fmt.Fprintf(conn, msg+"\n")

	message, err := bufio.NewReader(conn).ReadString('\n')
	if c.enableDebug {
		log.Logger().Debugf("sendMessage: message %s", string(message))
	}

	return transport.MessageResponseCode(message), err
}

func (c Client) executeIndex(packageName string, dependencies []string) (transport.MessageResponseCode, error) {
	msg := fmt.Sprintf("INDEX|%s|", packageName)

	for _, dependency := range dependencies {
		if dependency != "" {
			msg = msg + "," + dependency
		}
	}
	if c.enableDebug {
		log.Logger().Debugf("executeIndex: message %s", msg)
	}

	return c.sendMessage(msg)
}

func (c Client) executeRemove(packageName string) (transport.MessageResponseCode, error) {
	msg := fmt.Sprintf("REMOVE|%s|", packageName)

	if c.enableDebug {
		log.Logger().Debugf("executeRemove: message %s", msg)
	}

	return c.sendMessage(msg)
}

func (c Client) executeQuery(packageName string) (transport.MessageResponseCode, error) {
	msg := fmt.Sprintf("QUERY|%s|", packageName)

	if c.enableDebug {
		log.Logger().Debugf("executeQuery: message %s", msg)
	}

	return c.sendMessage(msg)
}
