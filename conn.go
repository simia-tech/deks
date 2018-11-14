package edkvs

import (
	"net"

	"github.com/mediocregopher/radix.v2/redis"
	"github.com/simia-tech/errx"
)

// Conn implements a edkvs client connection based on the redis protocol.
type Conn struct {
	conn   net.Conn
	client *redis.Client
}

// Dial establishes a connection to the server at the provided address.
func Dial(network, address string) (*Conn, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, errx.Annotatef(err, "dial [%s %s]", network, address)
	}

	client, err := redis.NewClient(conn)
	if err != nil {
		return nil, errx.Annotatef(err, "new client")
	}

	return &Conn{
		conn:   conn,
		client: client,
	}, nil
}

// Close tears down the connection.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// Reconsilate sets the server into reconsilation mode and returns the underlying connection.
func (c *Conn) Reconsilate() (net.Conn, error) {
	response := c.client.Cmd(cmdReconcilate)
	if !isOK(response) {
		return nil, errx.Errorf("reconsilate command failed")
	}
	c.client = nil
	return c.conn, nil
}

func (c *Conn) setItem(kh keyHash, item []byte) error {
	response := c.client.Cmd(cmdSetItem, kh[:], item)
	if !response.IsType(redis.Str) {
		return errx.Errorf("set item command failed")
	}
	s, err := response.Str()
	if err != nil {
		return errx.Annotatef(err, "response string")
	}
	if s != "OK" {
		return errx.Errorf("set item command failed")
	}
	return nil
}

func (c *Conn) getItem(kh keyHash) ([]byte, error) {
	response := c.client.Cmd(cmdGetItem, kh[:])
	if !response.IsType(redis.Str) {
		return nil, errx.Errorf("get item command failed")
	}
	bytes, err := response.Bytes()
	if err != nil {
		return nil, errx.Annotatef(err, "response bytes")
	}
	return bytes, nil
}

func isOK(response *redis.Resp) bool {
	if response.IsType(redis.Str) {
		if s, _ := response.Str(); s == "OK" {
			return true
		}
	}
	return false
}
