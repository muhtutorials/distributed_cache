package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Command byte

const (
	CmdNone Command = iota
	CmdJoin
	CmdSet
	CmdGet
	CmdDel
)

type Status byte

const (
	StatusNone Status = iota
	StatusOK
	StatusError
	StatusKeyNotFound
)

type CommandJoin struct {
}

type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int64
}

type CommandGet struct {
	Key []byte
}

type ResponseSet struct {
	Status Status
}

type ResponseGet struct {
	Status Status
	Value  []byte
}

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusError:
		return "Error"
	case StatusKeyNotFound:
		return "Key Not Found"
	default:
		return "No Status"
	}
}

func ParseCommand(r io.Reader) (any, error) {
	var cmd Command
	if err := binary.Read(r, binary.LittleEndian, &cmd); err != nil {
		return nil, err
	}

	switch cmd {
	case CmdJoin:
		return CommandJoin{}, nil
	case CmdSet:
		return parseSetCommand(r), nil
	case CmdGet:
		return parseGetCommand(r), nil
	default:
		return nil, errors.New("invalid command")
	}
}

func (r *ResponseSet) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, r.Status)

	return buf.Bytes()
}

func ParseResponseSet(r io.Reader) (*ResponseSet, error) {
	res := &ResponseSet{}

	err := binary.Read(r, binary.LittleEndian, &res.Status)

	return res, err
}

func (r *ResponseGet) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, r.Status)
	binary.Write(buf, binary.LittleEndian, int32(len(r.Value)))
	binary.Write(buf, binary.LittleEndian, r.Value)

	return buf.Bytes()
}

func ParseResponseGet(r io.Reader) (*ResponseGet, error) {
	res := &ResponseGet{}
	err := binary.Read(r, binary.LittleEndian, &res.Status)
	if err != nil {
		return nil, err
	}

	var valueLen int32
	err = binary.Read(r, binary.LittleEndian, &valueLen)
	if err != nil {
		return nil, err
	}

	res.Value = make([]byte, valueLen)
	err = binary.Read(r, binary.LittleEndian, &res.Value)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *CommandSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdSet)

	binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)

	binary.Write(buf, binary.LittleEndian, int32(len(c.Value)))
	binary.Write(buf, binary.LittleEndian, c.Value)

	binary.Write(buf, binary.LittleEndian, c.TTL)

	return buf.Bytes()
}

func parseSetCommand(r io.Reader) *CommandSet {
	cmd := &CommandSet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &cmd.Key)

	var valueLen int32
	binary.Read(r, binary.LittleEndian, &valueLen)
	cmd.Value = make([]byte, valueLen)
	binary.Read(r, binary.LittleEndian, &cmd.Value)

	var ttl int64
	binary.Read(r, binary.LittleEndian, &ttl)
	cmd.TTL = ttl

	return cmd
}

func (c *CommandGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdGet)

	binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)

	return buf.Bytes()
}

func parseGetCommand(r io.Reader) *CommandGet {
	cmd := &CommandGet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &cmd.Key)

	return cmd
}
