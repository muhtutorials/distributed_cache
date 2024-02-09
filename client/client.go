package client

import (
	"context"
	"distributed_cache/protocol"
	"fmt"
	"net"
)

type Options struct {
}

type Client struct {
	conn net.Conn
}

func NewFromConn(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func New(endpoint string, opts Options) (*Client, error) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Set(ctx context.Context, key, value []byte, ttl int64) error {
	cmd := &protocol.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return err
	}

	res, err := protocol.ParseResponseSet(c.conn)
	if err != nil {
		return err
	}

	if res.Status != protocol.StatusOK {
		return fmt.Errorf("server responded with a non OK status [%s]", res.Status)
	}

	fmt.Println("Response SET status:", res.Status)

	return nil
}

func (c *Client) Get(ctx context.Context, key []byte) ([]byte, error) {
	cmd := &protocol.CommandGet{
		Key: key,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}

	res, err := protocol.ParseResponseGet(c.conn)
	if err != nil {
		return nil, err
	}

	if res.Status == protocol.StatusKeyNotFound {
		return nil, fmt.Errorf("key [%s] was not found: %s", key, res.Status)
	}

	if res.Status != protocol.StatusOK {
		return nil, fmt.Errorf("server responded with a non OK status [%s]", res.Status)
	}

	fmt.Println("Response GET status:", res.Status)
	fmt.Println("Response GET status:", string(res.Value))

	return res.Value, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
