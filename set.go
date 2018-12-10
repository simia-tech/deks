package kea

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

// Set defines a set of items.
type Set struct {
	pf *recon.MemPrefixTree
}

// NewSet returns a new empty set.
func NewSet() *Set {
	pf := &recon.MemPrefixTree{}
	pf.Init()
	return &Set{pf: pf}
}

// Insert adds the provided item to the set.
func (s *Set) Insert(item Item) error {
	if err := s.pf.Insert(conflux.Zb(p, item[:])); err != nil {
		return errx.Annotatef(err, "insert")
	}
	return nil
}

// Remove removes the provided item from the set.
func (s *Set) Remove(item Item) error {
	if err := s.pf.Remove(conflux.Zb(p, item[:])); err != nil {
		return errx.Annotatef(err, "remove")
	}
	return nil
}

// Items returns a slice of all items in the set.
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

// Len returns the length of the set.
func (s *Set) Len() int {
	return s.pf.Len()
}

func (s *Set) prefixTree() recon.PrefixTree {
	return s.pf
}

// ItemSize specifies the item size.
const ItemSize = keyHashSize + 8

// Item defines the set item.
type Item [ItemSize]byte
