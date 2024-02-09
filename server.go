package main

import (
	"context"
	"distributed_cache/cache"
	"distributed_cache/client"
	"distributed_cache/protocol"
	"encoding/binary"
	"fmt"
	"go.uber.org/zap"
	"io"
	"log"
	"net"
	"time"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	cache   cache.Cacher
	members map[*client.Client]struct{}
	logger  *zap.SugaredLogger
}

func NewServer(opts ServerOpts, cache cache.Cacher) *Server {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	return &Server{
		ServerOpts: opts,
		cache:      cache,
		members:    make(map[*client.Client]struct{}),
		logger:     sugar,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	s.logger.Infow("server starting",
		"port", s.ListenAddr,
		"leader", s.IsLeader,
	)

	if !s.IsLeader {
		go func() {
			if err = s.dialLeader(); err != nil {
				log.Println(err)
			}
		}()
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) dialLeader() error {
	conn, err := net.Dial("tcp", s.LeaderAddr)
	if err != nil {
		return fmt.Errorf("failed to dial leader [%s]", s.ListenAddr)
	}

	s.logger.Infow("connected to leader", "addr", s.LeaderAddr)

	binary.Write(conn, binary.LittleEndian, protocol.CmdJoin)

	s.handleConn(conn)

	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("connection made:", conn.RemoteAddr())

	for {
		cmd, err := protocol.ParseCommand(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("parse command error", err)
			break
		}

		go s.handleCommand(conn, cmd)
	}

	fmt.Println("connection closed", conn.RemoteAddr())
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *protocol.CommandJoin:
		s.handleJoinCommand(conn, v)
	case *protocol.CommandSet:
		s.handleSetCommand(conn, v)
	case *protocol.CommandGet:
		s.handleGetCommand(conn, v)
	}
}

func (s *Server) handleJoinCommand(conn net.Conn, cmd *protocol.CommandJoin) error {
	fmt.Println("member just joined the cluster:", conn.RemoteAddr())

	s.members[client.NewFromConn(conn)] = struct{}{}

	return nil
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *protocol.CommandSet) error {
	go func() {
		for member := range s.members {
			err := member.Set(context.Background(), cmd.Key, cmd.Value, cmd.TTL)
			if err != nil {
				log.Println("forward to member error:", err)
			}
		}
	}()

	var res protocol.ResponseSet

	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)); err != nil {
		res.Status = protocol.StatusError
		conn.Write(res.Bytes())
		return err
	}

	res.Status = protocol.StatusOK
	conn.Write(res.Bytes())

	return nil
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *protocol.CommandGet) error {
	var res protocol.ResponseGet

	value, err := s.cache.Get(cmd.Key)
	if err != nil {
		res.Status = protocol.StatusError
		conn.Write(res.Bytes())
		return err
	}

	res.Status = protocol.StatusOK
	res.Value = value

	conn.Write(res.Bytes())

	return nil
}
