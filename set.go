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
	pf *recon.MemPrefixTree
}

func NewSet() *Set {
	pf := &recon.MemPrefixTree{}
	pf.Init()
	return &Set{pf: pf}
}

func (s *Set) Insert(item Item) error {
	if err := s.pf.Insert(conflux.Zb(p, item[:])); err != nil {
		return errx.Annotatef(err, "insert")
	}
	return nil
}

func (s *Set) Remove(item Item) error {
	if err := s.pf.Remove(conflux.Zb(p, item[:])); err != nil {
		return errx.Annotatef(err, "remove")
	}
	return nil
}

func (s *Set) Items() []Item {
	items := s.pf.Items()
	results := make([]Item, len(items))
	for index, item := range items {
		i := Item{}
		copy(i[:], item.Bytes())
		results[index] = i
	}
	return results
}

func (s *Set) Len() int {
	return s.pf.Len()
}

func (s *Set) prefixTree() recon.PrefixTree {
	return s.pf
}

const ItemSize = keyHashSize + 8

type Item [ItemSize]byte
