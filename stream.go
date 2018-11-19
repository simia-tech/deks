package edkvs

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/simia-tech/errx"
)

type stream struct {
	ctx                   context.Context
	cancel                context.CancelFunc
	peerURL               string
	peerPingInterval      time.Duration
	peerReconnectInterval time.Duration
	updates               chan streamUpdate
	updatesMutex          sync.Mutex
	metric                Metric
}

type streamUpdate struct {
	keyHash   keyHash
	container *container
}

func newStream(
	peerURL string,
	peerPingInterval time.Duration,
	peerReconnectInterval time.Duration,
	m Metric,
) *stream {
	ctx, cancel := context.WithCancel(context.Background())
	s := &stream{
		ctx:                   ctx,
		cancel:                cancel,
		peerURL:               peerURL,
		peerPingInterval:      peerPingInterval,
		peerReconnectInterval: peerReconnectInterval,
		metric:                m,
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
			s.metric.PeerConnected(s.peerURL)
			if err := s.connect(); err != nil {
				log.Printf("stream [%s]: %v", s.peerURL, err)
			}
			s.metric.PeerDisconnected(s.peerURL)
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

	ticker := time.NewTicker(s.peerPingInterval)

	updates := make(chan streamUpdate)
	defer func() {
		s.updatesMutex.Lock()
		s.updates = nil
		close(updates)
		s.updatesMutex.Unlock()

		ticker.Stop()
	}()
	s.updatesMutex.Lock()
	s.updates = updates
	s.updatesMutex.Unlock()

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-ticker.C:
			if err := conn.Ping(); err != nil {
				return errx.Annotatef(err, "ping")
			}
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
	s.updatesMutex.Lock()
	if s.updates == nil {
		s.updatesMutex.Unlock()
		return
	}
	s.updates <- streamUpdate{kh, container}
	s.updatesMutex.Unlock()
}

func (s *stream) close() {
	s.cancel()
}
