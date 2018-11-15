package edkvs

import (
	"context"
	"log"
	"time"

	"github.com/simia-tech/errx"
)

type stream struct {
	ctx                   context.Context
	cancel                context.CancelFunc
	peerURL               string
	peerReconnectInterval time.Duration
	updates               chan streamUpdate
}

type streamUpdate struct {
	keyHash   keyHash
	container *container
}

func newStream(peerURL string, peerReconnectInterval time.Duration) *stream {
	ctx, cancel := context.WithCancel(context.Background())
	s := &stream{
		ctx:                   ctx,
		cancel:                cancel,
		peerURL:               peerURL,
		peerReconnectInterval: peerReconnectInterval,
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
				log.Printf("stream [%s]: %v", s.peerURL, err)
			}
			time.Sleep(s.peerReconnectInterval)
		}
	}
}

func (s *stream) connect() error {
	conn, err := Dial(s.peerURL)
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
