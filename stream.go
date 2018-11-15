package edkvs

import (
	"context"
	"log"
	"time"

	"github.com/simia-tech/errx"
)

type stream struct {
	ctx     context.Context
	cancel  context.CancelFunc
	network string
	address string
	updates chan streamUpdate
}

type streamUpdate struct {
	keyHash   keyHash
	container *container
}

func newStream(network, address string) *stream {
	ctx, cancel := context.WithCancel(context.Background())
	s := &stream{
		ctx:     ctx,
		cancel:  cancel,
		network: network,
		address: address,
	}
	go s.loop()
	return s
}

func (s *stream) loop() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			if err := s.connect(); err != nil {
				log.Printf("stream [%s %s]: %v", s.network, s.address, err)
			}
			time.Sleep(time.Second)
		}
	}
}

func (s *stream) connect() error {
	conn, err := Dial(s.network, s.address)
	if err != nil {
		return errx.Annotatef(err, "dial")
	}
	defer conn.Close()

	updates := make(chan streamUpdate)
	defer func() {
		s.updates = nil
		close(updates)
	}()
	s.updates = updates

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case u := <-updates:
			bytes, err := u.container.MarshalBinary()
			if err != nil {
				return errx.Annotatef(err, "marshal binary")
			}
			if err := conn.setContainer(u.keyHash, bytes); err != nil {
				return errx.Annotatef(err, "set container")
			}
		}
	}
}

func (s *stream) update(kh keyHash, container *container) {
	if s.updates == nil {
		return
	}
	s.updates <- streamUpdate{kh, container}
}

func (s *stream) close() {
	s.cancel()
}
