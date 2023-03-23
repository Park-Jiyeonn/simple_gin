package gee

import "strings"

type node struct {
	pattern  string // 待匹配路由
	part     string // 路由里面的一部分
	children []*node
	isWild   bool // 模糊匹配。part 要是含有 : 或者 * 的时候，为 true
}

// 在x的孩子中，查找第一个匹配上的节点
func (x *node) matchChild(part string) *node {
	for _, child := range x.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 在x中，找所有匹配上的节点
func (x *node) matchChildren(part string) []*node {
	ret := make([]*node, 0)
	for _, child := range x.children {
		if child.part == part || child.isWild {
			ret = append(ret, child)
		}
	}
	return ret
}

func (x *node) insert(pattern string, parts []string, step int) {
	// 走到底了，parts整个都遍历完了，可以return了
	if len(parts) == step {
		// 这里的这个赋值有什么用呢，
		// 针对 * 的匹配吧，一个 URL 要是凭借着 * 匹配到了某个位置，但是这里并不是正确的路径，那还是失败的匹配
		x.pattern = pattern
		return
	}
	part := parts[step]
	child := x.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		// 新增一个孩子
		x.children = append(x.children, child)
	}
	child.insert(pattern, parts, step+1)
}

func (x *node) search(parts []string, step int) *node {
	if len(parts) == step || strings.HasPrefix(x.part, "*") {
		if x.pattern == "" {
			return nil
		}
		return x
	}
	part := parts[step]
	children := x.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, step+1)
		if result != nil {
			return result
		}
	}
	return nil
}
