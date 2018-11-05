package edkvs

import (
	"math/big"

	"github.com/simia-tech/conflux"
	"github.com/simia-tech/conflux/recon"
	"github.com/simia-tech/errx"
)

var p *big.Int

func init() {
	p = big.NewInt(0)
	p.SetString("530512889551602322505127520352579437339", 10)
}

type Set struct {
	prefixTree *recon.MemPrefixTree
}

func NewSet() *Set {
	prefixTree := &recon.MemPrefixTree{}
	prefixTree.Init()
	return &Set{prefixTree: prefixTree}
}

func (s *Set) Insert(item []byte) error {
	if err := s.prefixTree.Insert(conflux.Zb(p, item)); err != nil {
		return errx.Annotatef(err, "insert")
	}
	return nil
}

func (s *Set) Remove(item []byte) error {
	if err := s.prefixTree.Remove(conflux.Zb(p, item)); err != nil {
		return errx.Annotatef(err, "remove")
	}
	return nil
}

func (s *Set) Items() [][]byte {
	items := s.prefixTree.Items()
	results := make([][]byte, len(items))
	for index, item := range items {
		results[index] = item.Bytes()
	}
	return results
}

func (s *Set) Len() int {
	return s.prefixTree.Len()
}
