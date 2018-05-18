// Implementation of an R-Way Trie data structure.
//
// A Trie has a root Node which is the base of the tree.
package trie

import (
	"sort"
	"unicode/utf8"
)

type NodeIterator func(key []byte, val interface{}) bool

type Node struct {
	code     rune         // code of node
	term     bool         // last node flag
	depth    int
	value interface{}  // property of node
	parent   *Node
	children map[rune]*Node
}

type Trie struct {
	root *Node
	size int
}

type ByRune []rune
func (a ByRune) Len() int           { return len(a) }
func (a ByRune) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRune) Less(i, j int) bool { return a[i] < a[j] }

// Creates a new Trie with an initialized root Node.
func NewTrie() *Trie {
	return &Trie{
		root: &Node{children: make(map[rune]*Node), depth: 0},
		size: 0,
	}
}

// Returns the root node for the Trie.
func (t *Trie) Root() *Node {
	return t.root
}

func (t *Trie) Size() int {
	return t.size
}

// ReplaceOrInsert adds the given key to the tree.  If an key in the tree
// already equals the given one, it is removed from the tree and returned.
// Otherwise, nil is returned.
func (t *Trie) ReplaceOrInsert(key []byte, value interface{}) *Node {
	if len(key) == 0 {
		return nil
	}
	node := t.root
	var pre *Node
	offset := 0
	for len(key[offset:]) > 0 {
		e, size := utf8.DecodeRune(key[offset:])
		if e == utf8.RuneError {
			return nil
		}
		offset += size
		if n, ok := node.children[e]; ok {
			node = n
			pre = n
		} else {
			node = node.NewChildNode(e, nil, false)
			pre = nil
		}
	}

	// new node
	if pre == nil {
		node.value = value
		node.term = true
		t.size++
	} else {
		node = &Node{
			code:     pre.code,
			term:     pre.term,
			value:    value,
			parent:   pre.parent,
			children: pre.children,
			depth:    pre.depth,
		}
		node.parent.ReplaceOrInsertChildNode(node)
	}
	return pre
}

// Finds and returns property data associated
// with `key`.
func (t *Trie) Find(key []byte) (*Node, bool) {
	keyRune := parseTextToRunes(key)
	node := findNode(t.Root(), keyRune)
	if node == nil {
		return nil, false
	}
	if !node.term {
		return nil, false
	}

	return &Node{
		code:     node.code,
		parent:   node.Parent(),
		depth:    node.Depth(),
		term:     node.Terminating(),
		value:    node.value,
	}, true
}

func (t *Trie) HasKeysWithPrefix(key []byte) bool {
	keyRune := parseTextToRunes(key)
	node := findNode(t.Root(), keyRune)
	return node != nil
}

// Removes a key from the trie.
// Return delete node if exist
// Note make sure the key is not only a prefix
func (t *Trie) Delete(key []byte) *Node {
	keyRune := parseTextToRunes(key)
	node := findNode(t.Root(), keyRune)
	var del *Node
	if node.term {
		t.size--
		del = node
		if len(node.children) > 0 {
			// we just flag the term
			node.term = false
		} else {
			// no children node, we need delete from parent node
			if node.Parent() != nil {
				node.parent.RemoveChild(node.code)
				// check the parent if the node has no children nodes
				for n := node.Parent(); n != nil; n = n.Parent() {
					if n.term {
						break
					}
					if len(n.children) > 0 {
						break
					}
					if n.Parent() != nil {
						n.parent.RemoveChild(n.code)
					}
				}
			}
		}
		return del
	} else {
		// not end node
		return nil
	}
}

// Returns all the keys currently stored in the trie.
func (t *Trie) Keys() [][]byte {
	var keys [][]byte
	iter := func(key []byte, val interface{}) bool {
		k := make([]byte, len(key))
		copy(k, key)
		keys = append(keys, k)
		return true
	}
	t.PrefixSearch(nil, iter)
	return keys
}

// Performs a prefix search against the keys in the trie.
// The key and value are only valid for the life of the iterator.
func (t *Trie) PrefixSearch(pre []byte, iter NodeIterator) {
	preRune := parseTextToRunes(pre)
	node := findNode(t.Root(), preRune)
	if node == nil {
		return
	}

	preTraverse(node, preRune, iter)
}

// Creates and returns a pointer to a new child for the node.
func (n *Node) NewChildNode(code rune, value interface{}, term bool) *Node {
	node := &Node{
		code:     code,
		term:     term,
		value:    value,
		parent:   n,
		children: make(map[rune]*Node),
		depth:    n.depth + 1,
	}
	n.children[code] = node
	return node
}

func (n *Node) ReplaceOrInsertChildNode(node *Node) {
	n.children[node.Code()] = node
}

func (n *Node) RemoveChild(r rune) {
	delete(n.children, r)
}

// Returns the parent of this node.
func (n Node) Parent() *Node {
	return n.parent
}

// Returns the children of this node.
func (n Node) Children() map[rune]*Node {
	return n.children
}

func (n Node) Terminating() bool {
	return n.term
}

func (n Node) Depth() int {
	return n.depth
}

func (n *Node) Code() rune {
	return n.code
}

func (n *Node) Value() interface{} {
	return n.value
}

func findNode(node *Node, key []rune) *Node {
	if node == nil {
		return nil
	}

	if len(key) == 0 {
		return node
	}

	n, ok := node.Children()[key[0]]
	if !ok {
		return nil
	}

	var subKey []rune
	if len(key) > 1 {
		subKey = key[1:]
	} else {
		subKey = key[0:0]
	}

	return findNode(n, subKey)
}

// Preorder traverse trie
func preTraverse(node *Node, prefix []rune, iter NodeIterator) {
	if node == nil {
		return
	}
	if node.term {
		if !iter(parseRunesToText(prefix), node.Value()) {
			return
		}
	}
	if len(node.Children()) == 0 {
		return
	}
	// sort key
	bs := make([]rune, 0, len(node.Children()))
	for val, _ := range node.children {
		bs = append(bs, val)
	}
	sort.Sort(ByRune(bs))
	for _, c := range bs {
		if n, ok := node.children[c]; ok {
			preTraverse(n, append(prefix, n.code), iter)
		}
	}
}

func parseTextToRunes(str []byte) []rune {
	if len(str) == 0 {
		return nil
	}

	var keyRune []rune
	offset := 0
	for len(str[offset:]) > 0 {
		e, size := utf8.DecodeRune(str[offset:])
		if e == utf8.RuneError {
			return nil
		}
		offset += size
		keyRune = append(keyRune, e)
	}
	return keyRune
}

func parseRunesToText(runes []rune) []byte {
	var l int
	for _, r := range runes {
		l += utf8.RuneLen(r)
	}
	_key := make([]byte, l)
	offset := 0
	for _, k := range runes {
		n := utf8.EncodeRune(_key[offset:], k)
		offset += n
	}
	return _key
}

