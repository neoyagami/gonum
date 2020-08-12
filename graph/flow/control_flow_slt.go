// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flow

import "neoyagami/gonum/graph"

// DominatorsSLT returns a dominator tree for all nodes in the flow graph
// g starting from the given root node using the sophisticated version of
// the Lengauer-Tarjan algorithm. The SLT algorithm may outperform the
// simple LT algorithm for very large dense graphs.
func DominatorsSLT(root graph.Node, g graph.Directed) DominatorTree {
	// The algorithm used here is essentially the
	// sophisticated Lengauer and Tarjan algorithm
	// described in
	// https://doi.org/10.1145%2F357062.357071

	lt := sLengauerTarjan{
		indexOf: make(map[int64]int),
		base:    sltNode{semi: -1},
	}
	lt.base.label = &lt.base

	// step 1.
	lt.dfs(g, root)

	for i := len(lt.nodes) - 1; i > 0; i-- {
		w := lt.nodes[i]

		// step 2.
		for _, v := range w.pred {
			u := lt.eval(v)

			if u.semi < w.semi {
				w.semi = u.semi
			}
		}

		lt.nodes[w.semi].bucket[w] = struct{}{}
		lt.link(w.parent, w)

		// step 3.
		for v := range w.parent.bucket {
			delete(w.parent.bucket, v)

			u := lt.eval(v)
			if u.semi < v.semi {
				v.dom = u
			} else {
				v.dom = w.parent
			}
		}
	}

	// step 4.
	for _, w := range lt.nodes[1:] {
		if w.dom.node.ID() != lt.nodes[w.semi].node.ID() {
			w.dom = w.dom.dom
		}
	}

	// Construct the public-facing dominator tree structure.
	dominatorOf := make(map[int64]graph.Node)
	dominatedBy := make(map[int64][]graph.Node)
	for _, w := range lt.nodes[1:] {
		dominatorOf[w.node.ID()] = w.dom.node
		did := w.dom.node.ID()
		dominatedBy[did] = append(dominatedBy[did], w.node)
	}
	return DominatorTree{root: root, dominatorOf: dominatorOf, dominatedBy: dominatedBy}
}

// sLengauerTarjan holds global state of the Lengauer-Tarjan algorithm.
// This is a mapping between nodes and the postordering of the nodes.
type sLengauerTarjan struct {
	// nodes is the nodes traversed during the
	// Lengauer-Tarjan depth-first-search.
	nodes []*sltNode
	// indexOf contains a mapping between
	// the id-dense representation of the
	// graph and the potentially id-sparse
	// nodes held in nodes.
	//
	// This corresponds to the vertex
	// number of the node in the Lengauer-
	// Tarjan algorithm.
	indexOf map[int64]int

	// base is the base label for balanced
	// tree path compression used in the
	// sophisticated Lengauer-Tarjan
	// algorith,
	base sltNode
}

// sltNode is a graph node with accounting for the Lengauer-Tarjan
// algorithm.
//
// For the purposes of documentation the ltNode is given the name w.
type sltNode struct {
	node graph.Node

	// parent is vertex which is the parent of w
	// in the spanning tree generated by the search.
	parent *sltNode

	// pred is the set of vertices v such that (v, w)
	// is an edge of the graph.
	pred []*sltNode

	// semi is a number defined as follows:
	// (i)  After w is numbered but before its semidominator
	//      is computed, semi is the number of w.
	// (ii) After the semidominator of w is computed, semi
	//      is the number of the semidominator of w.
	semi int

	// size is the tree size of w used in the
	// sophisticated algorithm.
	size int

	// child is the child node of w used in the
	// sophisticated algorithm.
	child *sltNode

	// bucket is the set of vertices whose
	// semidominator is w.
	bucket map[*sltNode]struct{}

	// dom is vertex defined as follows:
	// (i)  After step 3, if the semidominator of w is its
	//      immediate dominator, then dom is the immediate
	//      dominator of w. Otherwise dom is a vertex v
	//      whose number is smaller than w and whose immediate
	//      dominator is also w's immediate dominator.
	// (ii) After step 4, dom is the immediate dominator of w.
	dom *sltNode

	// In general ancestor is nil only if w is a tree root
	// in the forest; otherwise ancestor is an ancestor
	// of w in the forest.
	ancestor *sltNode

	// Initially label is w. It is adjusted during
	// the algorithm to maintain invariant (3) in the
	// Lengauer and Tarjan paper.
	label *sltNode
}

// dfs is the Sophisticated Lengauer-Tarjan DFS procedure.
func (lt *sLengauerTarjan) dfs(g graph.Directed, v graph.Node) {
	i := len(lt.nodes)
	lt.indexOf[v.ID()] = i
	ltv := &sltNode{
		node:   v,
		semi:   i,
		size:   1,
		child:  &lt.base,
		bucket: make(map[*sltNode]struct{}),
	}
	ltv.label = ltv
	lt.nodes = append(lt.nodes, ltv)

	to := g.From(v.ID())
	for to.Next() {
		w := to.Node()
		wid := w.ID()

		idx, ok := lt.indexOf[wid]
		if !ok {
			lt.dfs(g, w)

			// We place this below the recursive call
			// in contrast to the original algorithm
			// since w needs to be initialised, and
			// this happens in the child call to dfs.
			idx, ok = lt.indexOf[wid]
			if !ok {
				panic("path: unintialized node")
			}
			lt.nodes[idx].parent = ltv
		}
		ltw := lt.nodes[idx]
		ltw.pred = append(ltw.pred, ltv)
	}
}

// compress is the Sophisticated Lengauer-Tarjan COMPRESS procedure.
func (lt *sLengauerTarjan) compress(v *sltNode) {
	if v.ancestor.ancestor != nil {
		lt.compress(v.ancestor)
		if v.ancestor.label.semi < v.label.semi {
			v.label = v.ancestor.label
		}
		v.ancestor = v.ancestor.ancestor
	}
}

// eval is the Sophisticated Lengauer-Tarjan EVAL function.
func (lt *sLengauerTarjan) eval(v *sltNode) *sltNode {
	if v.ancestor == nil {
		return v.label
	}
	lt.compress(v)
	if v.ancestor.label.semi >= v.label.semi {
		return v.label
	}
	return v.ancestor.label
}

// link is the Sophisticated Lengauer-Tarjan LINK procedure.
func (*sLengauerTarjan) link(v, w *sltNode) {
	s := w
	for w.label.semi < s.child.label.semi {
		if s.size+s.child.child.size >= 2*s.child.size {
			s.child.ancestor = s
			s.child = s.child.child
		} else {
			s.child.size = s.size
			s.ancestor = s.child
			s = s.child
		}
	}
	s.label = w.label
	v.size += w.size
	if v.size < 2*w.size {
		s, v.child = v.child, s
	}
	for s != nil {
		s.ancestor = v
		s = s.child
	}
}
