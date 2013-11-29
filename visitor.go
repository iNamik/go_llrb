package llrb

import . "github.com/iNamik/go_pkg/debug/assert"

import "fmt"

import (
	"github.com/iNamik/go_bst/visitor"
	"github.com/iNamik/go_cmp"
)

// tree::Visit
func (t *tree) Visit(key interface{}, f visitor.F) (value interface{}, result visitor.Result) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.root, value, result = visit(t.root, key, t.fcmp, f)
	if t.root != nil {
		t.root.red = false
	}
	if result == visitor.INSERTED {
		t.size++
	} else if result == visitor.REMOVED {
		t.size--
	}
	Assert((t.root != nil && t.size > 0) || (t.root == nil && t.size == 0))
	return value, result
}

// visit is based on sedgewick's original remove logic.
// NOTE: Inserting a node can create a new scenario that must be checked during fixup.
func visit(h *node, key interface{}, fcmp cmp.F, f visitor.F) (_ *node, value interface{}, result visitor.Result) {
	if h == nil {
		var action visitor.Action
		value, action = f(nil, false)
		switch action {
		case visitor.INSERT:
			// NOTE: When attaching to a right-red, will not be caught during normal fixup
			return &node{key: key, value: value, red: true}, value, visitor.INSERTED
		case visitor.GET:
			return nil, nil, visitor.NOT_FOUND
		default:
			panic(fmt.Sprintf("illegal action '%s' when visiting non-found key", action))
		}
	}
	c := fcmp(key, h.key)
	if c == cmp.LT {
		if h.left != nil && !isRed(h.left) && !isRed(h.left.left) {
			h = moveRedLeft(h)
		}
		h.left, value, result = visit(h.left, key, fcmp, f)
	} else {
		if isRed(h.left) {
			h = rotateRight(h)
			c = fcmp(key, h.key)
		}
		if h.right != nil && !isRed(h.right) && !isRed(h.right.left) {
			h = moveRedRight(h)
			c = fcmp(key, h.key)
		}
		if c == cmp.EQ {
			var action visitor.Action
			value, action = f(h.value, true)
			switch action {
			case visitor.GET:
				value, result = h.value, visitor.FOUND
			case visitor.REPLACE:
				h.value, result = value, visitor.REPLACED
			case visitor.REMOVE:
				// If nothing to replace with, then return without fixup
				if h.right == nil {
					return nil, h.value, visitor.REMOVED
				} else {
					value, result = h.value, visitor.REMOVED
					// Replace current node key, value with successor and delete successor
					var m *node
					h.right, m = removeMin(h.right)
					h.key, h.value = m.key, m.value
				}
			default:
				panic(fmt.Sprintf("illegal action '%s' when visiting found key", action))
			}
		} else {
			h.right, value, result = visit(h.right, key, fcmp, f)
		}
	}
	// Inserting a node onto a right-red will create a double-red combo that the normal fixup
	// will not catch.  We catch it here, before calling fixup
	if result == visitor.INSERTED && isRed(h.right) && isRed(h.right.left) {
		h.right = rotateRight(h.right)
		h = rotateLeft(h)
	}
	return fixUp(h), value, result
}
