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
	keyHash keyHash
	item    *item
}

func newStream(network, address string) (*stream, error) {
	ctx, cancel := context.WithCancel(context.Background())
	s := &stream{
		ctx:     ctx,
		cancel:  cancel,
		network: network,
		address: address,
	}
	go s.loop()
	return s, nil
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
			bytes, err := u.item.MarshalBinary()
			if err != nil {
				return errx.Annotatef(err, "marshal binary")
			}
			if err := conn.setItem(u.keyHash, bytes); err != nil {
				return errx.Annotatef(err, "set item")
			}
		}
	}

	return nil
}

func (s *stream) update(kh keyHash, item *item) {
	if s.updates == nil {
		return
	}
	s.updates <- streamUpdate{kh, item}
}

func (s *stream) close() {
	s.cancel()
}
