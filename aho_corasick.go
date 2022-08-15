package sensitive

// 通配符
const commonChar = '_'

// Ac 自动机
type Ac struct {
	results []map[int]string
}

func (ac *Ac) fail(node *Node, c rune) *Node {
	var next *Node
	for {
		next = ac.next(node.Failure, c)
		if next == nil {
			if node.IsRootNode() {
				return node
			}
			node = node.Failure
			continue
		}
		return next
	}
}

func (ac *Ac) next(node *Node, c rune) *Node {
	// 匹配普通字符
	next, ok := node.Children[c]
	if ok && (c != commonChar || next.isEscape) {
		return next
	}
	// 匹配通配符
	next, ok = node.Children[commonChar]
	if ok && !next.isEscape {
		return next
	}
	return nil
}

func (ac *Ac) output(node *Node, runes []rune, position int) {
	if node.IsRootNode() {
		return
	}

	if node.IsPathEnd() {
		ac.results = append(ac.results, map[int]string{position + 1 - node.depth: string(runes[position+1-node.depth : position+1])})
	}

	ac.output(node.Failure, runes, position)
}

func (ac *Ac) firstOutput(node *Node, runes []rune, position int) string {
	if node.IsRootNode() {
		return ""
	}

	if node.IsPathEnd() {
		return string(runes[position+1-node.depth : position+1])
	}

	return ac.firstOutput(node.Failure, runes, position)
}

func (ac *Ac) replace(node *Node, runes []rune, position int, replace rune) {
	if node.IsRootNode() {
		return
	}

	if node.IsPathEnd() {
		for i := position + 1 - node.depth; i < position+1; i++ {
			runes[i] = replace
		}
	}

	ac.replace(node.Failure, runes, position, replace)
}
