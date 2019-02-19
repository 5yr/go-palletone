package _txspool

import "github.com/palletone/go-palletone/common"

type set interface {
	exist(elem common.Hash) bool
	insert(elem common.Hash)
	delete(elem common.Hash)
	size() int64
	loop() map[common.Hash]bool
}

type txHashSet struct {
	s map[common.Hash]bool
}

func newTxHashSet() *txHashSet {
	return &txHashSet{
		s: make(map[common.Hash]bool, 0),
	}
}

// loop : return a map inside of  set
func (set *txHashSet) loop() map[common.Hash]bool {
	return set.s
}

func (set *txHashSet) exist(elem common.Hash) bool {
	_, status := set.s[elem]
	return !status
}

func (set *txHashSet) insert(elem common.Hash) {
	set.s[elem] = true
}

func (set *txHashSet) insertList(elems []common.Hash) {
	for _, elem := range elems {
		set.insert(elem)
	}
}

func (set *txHashSet) size() int64 {
	return int64(len(set.s))
}

func (set *txHashSet) delete(elem common.Hash) {
	if set.exist(elem) {
		delete(set.s, elem)
	}
}

func (set *txHashSet) merge(rset *txHashSet) {
	for re := range rset.s {
		set.s[re] = true
	}
}

func (set *txHashSet) replaceBy(rset *txHashSet) {
	set.s = rset.s
}
