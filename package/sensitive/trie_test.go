package sensitive_test

import (
	"douyin/package/sensitive"
	"fmt"
	"testing"
)

func TestTrie(t *testing.T) {

	trie := sensitive.NewTrie()
	trie.Add("傻逼", nil)
	trie.Add("你妈的", nil)
	trie.Add("蓝色", nil)

	result, str := trie.Check("你是大傻 逼，逼，你妈，你妈的，我们这里有一个黄色的灯泡，他存在了很久。他是蓝色的。", "***")

	fmt.Printf("result:%#v, str:%v\n", result, str)

}
