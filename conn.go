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

// Set sets the provided value at the provided key.
func (c *Conn) Set(key, value []byte) error {
	response := c.client.Cmd(cmdSet, key, value)
	if !isOK(response) {
		return errx.Errorf("set command failed")
	}
	return nil
}

// Get returns the value at the provided key.
func (c *Conn) Get(key []byte) ([]byte, error) {
	response := c.client.Cmd(cmdGet, key)
	if !response.IsType(redis.Str) {
		return nil, errx.Errorf("get item command failed")
	}
	bytes, err := response.Bytes()
	if err != nil {
		return nil, errx.Annotatef(err, "response bytes")
	}
	return bytes, nil
}

// Delete removes the value at the provided key.
func (c *Conn) Delete(key []byte) error {
	response := c.client.Cmd(cmdDelete, key)
	if !isOK(response) {
		return errx.Errorf("delete command failed")
	}
	return nil
}

// Keys returns a slice containing all keys.
func (c *Conn) Keys() ([][]byte, error) {
	response := c.client.Cmd(cmdKeys)
	items, err := response.Array()
	if err != nil {
		return nil, errx.Annotatef(err, "response array")
	}
	keys := make([][]byte, len(items))
	for index, item := range items {
		key, err := item.Bytes()
		if err != nil {
			return nil, errx.Annotatef(err, "response bytes")
		}
		keys[index] = key
	}
	return keys, nil
}

// Tidy cleans up the store.
func (c *Conn) Tidy() error {
	response := c.client.Cmd(cmdTidy)
	if !isOK(response) {
		return errx.Errorf("tidy command failed")
	}
	return nil
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

func (c *Conn) setContainer(kh keyHash, item []byte) error {
	response := c.client.Cmd(cmdSetContainer, kh[:], item)
	if !isOK(response) {
		return errx.Errorf("set container command failed")
	}
	return nil
}

func (c *Conn) getContainer(kh keyHash) ([]byte, error) {
	response := c.client.Cmd(cmdGetContainer, kh[:])
	if !response.IsType(redis.Str) {
		return nil, errx.Errorf("get container command failed")
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
