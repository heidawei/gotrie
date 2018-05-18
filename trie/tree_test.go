package trie

import (
	"testing"
	"reflect"
)

func checkNode(t *testing.T, n *Node, ex string) {
	if n == nil {
		t.Fatal("find fault")
	}
	if n.Value() == nil {
		t.Fatalf("invalid node %v", n)
	}
	if n.Value().(string) != ex {
		t.Fatalf("invalid node %v, expect %s", n.Value().(string), ex)
	}
}

func insert(t *testing.T, tree *Trie, keys []string) {
	for _, key := range keys {
		n := tree.ReplaceOrInsert([]byte(key), key)
		if n != nil {
			t.Fatalf("insert failed, replace %v", n)
		}
	}
}

func insertKV(t *testing.T, tree *Trie, keys []string, values []string, ignoreReplace bool) {
	for i, key := range keys {
		n := tree.ReplaceOrInsert([]byte(key), values[i])
		if !ignoreReplace && n != nil {
			t.Fatalf("insert failed, replace %v", n)
		}
	}
}

func getCheck(t *testing.T, tree *Trie, keys []string) {
	for _, key := range keys {
		n, _ := tree.Find([]byte(key))
		checkNode(t, n, key)
	}
}

func getCheckKV(t *testing.T, tree *Trie, keys []string, values []string) {
	for i, key := range keys {
		n, _ := tree.Find([]byte(key))
		checkNode(t, n, values[i])
	}
}

func deleteCheck(t *testing.T, tree *Trie, keys []string) {
	for _, key := range keys {
		n := tree.Delete([]byte(key))
		checkNode(t, n, key)
	}
}

func TestInsert(t *testing.T) {
	tree := NewTrie()
	keys := []string{"a", "ab", "abc", "abd", "abcdef", "b", "bc", "bcd", "bce"}
	insert(t, tree, keys)
	getCheck(t, tree, keys)

	values := []string{"a1", "ab1", "abc1", "abd1", "abcdef1", "b1", "bc1", "bcd1", "bce1"}
	insertKV(t, tree, keys, values, true)
	getCheckKV(t, tree, keys, values)
}

func TestInsertGetDelete(t *testing.T) {
	tree := NewTrie()
	keys := []string{"a", "ab", "abc", "abd", "abcdef", "b", "bc", "bcd", "bce"}
	insert(t, tree, keys)
	getCheck(t, tree, keys)

	_, find := tree.Find([]byte("abcde"))
	if find {
		t.Fatal("insert failed")
	}

	// delete
	n := tree.Delete([]byte("a"))
	if n == nil {
		t.Fatal("delete failed")
	}
	checkNode(t, n, "a")
	_, find = tree.Find([]byte("a"))
	if find {
		t.Fatal("delete failed")
	}
	keys = []string{"ab", "abc", "abd", "b", "bc", "bcd", "bce"}
	getCheck(t, tree, keys)
}

func TestDelete(t *testing.T) {
	tree := NewTrie()
	keys := []string{"a", "ab", "abc", "abd", "abcdef"}
	insert(t, tree, keys)

	// delete
	n := tree.Delete([]byte("a"))
	if n == nil {
		t.Fatal("delete failed")
	}
	checkNode(t, n, "a")
	_, find := tree.Find([]byte("a"))
	if find {
		t.Fatal("delete failed")
	}
	keys = []string{"ab", "abc", "abd", "abcdef"}
	getCheck(t, tree, keys)

	n = tree.Delete([]byte("abcd"))
	if n != nil {
		t.Fatal("delete failed")
	}
	getCheck(t, tree, keys)

	n = tree.Delete([]byte("abc"))
	if n == nil {
		t.Fatalf("delete failed")
	}
	keys = []string{"ab", "abd", "abcdef"}
	getCheck(t, tree, keys)

	deleteCheck(t, tree, keys)
	if len(tree.Root().Children()) > 0 {
		t.Fatal("delete failed")
	}
}

func TestPrefixSearch(t *testing.T) {
	tree := NewTrie()
	keys := []string{"a", "ab", "abc", "abcdef", "abd", "b", "bc", "bcd", "bce"}
	insert(t, tree, keys)

	exKeys := []string{"a", "ab", "abc", "abcdef", "abd"}
	var iterKeys []string
	iter := func(key []byte, val interface{}) bool {
		iterKeys = append(iterKeys, string(key))
		return true
	}
	tree.PrefixSearch([]byte("a"), iter)
	if !reflect.DeepEqual(exKeys, iterKeys) {
		t.Fatalf("prefix search failed, ex %v, iter %v", exKeys, iterKeys)
	}

	exKeys = []string{"b", "bc", "bcd", "bce"}
	iterKeys = iterKeys[:0]
	iter = func(key []byte, val interface{}) bool {
		iterKeys = append(iterKeys, string(key))
		return true
	}
	tree.PrefixSearch([]byte("b"), iter)
	if !reflect.DeepEqual(exKeys, iterKeys) {
		t.Fatalf("prefix search failed, ex %v, iter %v", exKeys, iterKeys)
	}

	iterKeys = iterKeys[:0]
	all := tree.Keys()
	for _, key := range all {
		iterKeys = append(iterKeys, string(key))
	}
	if !reflect.DeepEqual(keys, iterKeys) {
		t.Fatalf("keys failed, ex %v, iter %v", keys, iterKeys)
	}
}
