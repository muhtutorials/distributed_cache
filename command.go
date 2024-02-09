package main

//import (
//	"errors"
//	"fmt"
//	"strconv"
//	"strings"
//	"time"
//)
//
//type Command string
//
//const (
//	CmdSet Command = "SET"
//	CmdGet Command = "GET"
//)
//
//type Message struct {
//	Cmd   Command
//	Key   []byte
//	Value []byte
//	TTL   time.Duration
//}
//
//func (m Message) ToBytes() []byte {
//	switch m.Cmd {
//	case CmdSet:
//		cmdStr := fmt.Sprintf("%s %s %s %d", m.Cmd, m.Key, m.Value, m.TTL)
//		return []byte(cmdStr)
//	case CmdGet:
//		cmdStr := fmt.Sprintf("%s %s", m.Cmd, m.Key)
//		return []byte(cmdStr)
//	default:
//		panic("unknown command")
//	}
//}
//
//func parseMessage(rawCmd []byte) (*Message, error) {
//	rawCmdStr := string(rawCmd)
//
//	parts := strings.Split(rawCmdStr, " ")
//	if len(parts) == 0 {
//		return nil, errors.New("invalid protocol format")
//	}
//
//	msg := &Message{
//		Cmd: Command(parts[0]),
//		Key: []byte(parts[1]),
//	}
//
//	if msg.Cmd == CmdSet {
//		if len(parts) < 4 {
//			return nil, errors.New("invalid SET command")
//		}
//		msg.Value = []byte(parts[2])
//
//		ttl, err := strconv.Atoi(parts[3])
//		if err != nil {
//			return nil, errors.New("invalid SET TTL command")
//		}
//		msg.TTL = time.Duration(ttl) * time.Millisecond
//	}
//
//	return msg, nil
//}
