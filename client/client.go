package client

import (
	"bufio"
	"fmt"
	"net"

	"github.com/hectorj2f/pkgInd/log"
	"github.com/hectorj2f/pkgInd/transport"
)

type Client struct {
	listenIP   string
	listenPort int
}

func (c Client) sendMessage(msg string) (transport.MessageResponseCode, error) {
	conn, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", c.listenIP, c.listenPort))
	defer conn.Close()

	fmt.Fprintf(conn, msg+"\n")

	message, err := bufio.NewReader(conn).ReadString('\n')
	log.Logger().Debugf("sendMessage: message %s", string(message))

	return transport.MessageResponseCode(message), err
}

func (c Client) executeIndex(packageName string, dependencies []string) (transport.MessageResponseCode, error) {
	msg := fmt.Sprintf("INDEX|%s|", packageName)

	for _, dependency := range dependencies {
		if dependency != "" {
			msg = msg + "," + dependency
		}
	}

	return c.sendMessage(msg)
}

func (c Client) executeRemove(packageName string) (transport.MessageResponseCode, error) {
	msg := fmt.Sprintf("REMOVE|%s|", packageName)

	return c.sendMessage(msg)
}

func (c Client) executeQuery(packageName string) (transport.MessageResponseCode, error) {
	msg := fmt.Sprintf("QUERY|%s|", packageName)

	return c.sendMessage(msg)
}
