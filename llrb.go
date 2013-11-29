package llrb

//import . "github.com/iNamik/go_pkg/debug/assert"

import "sync"

import (
	"github.com/iNamik/go_bst"
	"github.com/iNamik/go_bst/finder"
	"github.com/iNamik/go_bst/visitor"
	"github.com/iNamik/go_bst/walker"
	"github.com/iNamik/go_cmp"
)

/**********************************************************************
 ** Types & Interfaces
 **********************************************************************/

// T
type T interface {
	bst.T
	finder.I
	visitor.I
	walker.I
	bst.I_Size
	finder.I_Min
	finder.I_Max
}

// node
type node struct {
	key   interface{}
	value interface{}
	left  *node
	right *node
	red   bool
}

// tree
type tree struct {
	mutex *sync.Mutex
	root  *node
	fcmp  cmp.F
	size  int
}

/**********************************************************************
 ** Public Functions
 **********************************************************************/

// New
func New(fcmp cmp.F) T {
	return &tree{mutex: &sync.Mutex{}, root: nil, fcmp: fcmp, size: 0}
}

// Empty
func (t *tree) Empty() bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	//Assert((t.root != nil && t.size > 0) || (t.root == nil && t.size == 0))
	return t.root == nil
}

// Size
func (t *tree) Size() int {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	//Assert((t.root != nil && t.size > 0) || (t.root == nil && t.size == 0))
	return t.size
}

// tree:ReplaceOrInsert
func (t *tree) ReplaceOrInsert(key interface{}, value interface{}) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	var replaced bool
	t.root, replaced = replaceOrInsert(t.root, key, value, t.fcmp)
	//Assert(t.root != nil)
	t.root.red = false
	if replaced == false { // !replaced == inserted
		t.size++
	}
	//Assert((t.root != nil && t.size > 0) || (t.root == nil && t.size == 0))
	return replaced
}

// tree:Get
func (t *tree) Get(key interface{}) (interface{}, bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	h := get(t.root, key, t.fcmp)
	if h != nil {
		return h.value, true
	}
	return nil, false
}

// tree:Remove
func (t *tree) Remove(key interface{}) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	var removed bool
	t.root, removed = remove(t.root, key, t.fcmp)
	// RS did not check to see if root existed before coloring it
	if t.root != nil {
		t.root.red = false
	}
	if removed {
		t.size--
	}
	//Assert((t.root != nil && t.size > 0) || (t.root == nil && t.size == 0))
	return removed
}

/**********************************************************************
 ** Private Functions
 **********************************************************************/

// replaceOrInsert uses standard BST technique on the way down,
// adjusting the tree on the way up
func replaceOrInsert(h *node, key interface{}, value interface{}, fcmp cmp.F) (*node, bool) {
	if h == nil {
		return &node{key: key, value: value, red: true}, false
	}
	replaced := true
	switch fcmp(key, h.key) {
	case cmp.LT:
		h.left, replaced = replaceOrInsert(h.left, key, value, fcmp)
	case cmp.GT:
		h.right, replaced = replaceOrInsert(h.right, key, value, fcmp)
	default:
		h.value = value
	}
	return fixUp(h), replaced
}

// get uses standard BST technique
func get(h *node, key interface{}, fcmp cmp.F) *node {
	for h != nil {
		switch fcmp(key, h.key) {
		case cmp.LT:
			h = h.left
		case cmp.GT:
			h = h.right
		default:
			return h
		}
	}
	return nil
}

// remove
// NOTE: This may modify the tree state, even if it does not find @key
func remove(h *node, key interface{}, fcmp cmp.F) (*node, bool) {
	if h == nil {
		return nil, false
	}
	removed := false
	c := fcmp(key, h.key)
	// Less
	if c == cmp.LT {
		// If left is nil, then right is also nil, so nothing to do
		if h.left == nil {
			//Assert(h.right == nil)
			// We didn't find it, but we're done, so propagate h back up
			return h, false // No child nodes, so no need for fixup
		}
		//Assert(h.left != nil)
		// If we have a 2-node, then borrow from right to make 2-3 or 2-3-4
		if !isRed(h.left) && !isRed(h.left.left) {
			h = moveRedLeft(h)
			//Assert(h != nil)
		}
		// Move down (left)
		h.left, removed = remove(h.left, key, fcmp)

		// Greater or Equal
	} else {
		if isRed(h.left) {
			h = rotateRight(h) // Guarantees that key > h.key
			//Assert(h != nil)
			//Assert(h.red == false)
			//Assert(h.right != nil)
			//Assert(h.right.red == true)
			//Assert(h.right.left == nil || h.right.left.red == false)
			// Move down (right)
			h.right, removed = remove(h.right, key, fcmp)

			// If right is nil, then left is also nil, so nothing else to check
		} else if h.right == nil {
			//Assert(h.left == nil)
			// If we found the node, then there's nothing to replace it with
			if c == cmp.EQ {
				return nil, true
			}
			//Assert(c == cmp.GT)
			// We didn't find it, but we're done, so propagate h back up
			return h, false // No child nodes, so no need for fixup
		} else {
			//Assert(h.right != nil)
			// If we have a 2-node, then borrow from left to make 2-3 or 2-3-4
			if !isRed(h.right) && !isRed(h.right.left) {
				h = moveRedRight(h)
				//Assert(h != nil)
				// RS forgot to re-compare here
				c = fcmp(key, h.key)
				//Assert(c == cmp.EQ || c == cmp.GT)
			}
			// Found it!
			if c == cmp.EQ {
				// Replace current node key, value with successor and delete successor
				var m *node
				h.right, m = removeMin(h.right)
				h.key, h.value = m.key, m.value
				removed = true
			} else {
				//Assert(c == cmp.GT)
				// Move down (right)
				h.right, removed = remove(h.right, key, fcmp)
			}
		}
	}
	return fixUp(h), removed
}

// removeMin
func removeMin(h *node) (*node, *node) {
	//Assert(h != nil)
	// Remove node at bottom level
	// (h must be red by invariant)
	if h.left == nil {
		//Assert(h.red == true)
		return nil, h
	}
	// Push red link down if necessary
	if !isRed(h.left) && !isRed(h.left.left) {
		h = moveRedLeft(h)
	}
	var removed *node
	// Move down one level
	h.left, removed = removeMin(h.left)
	return fixUp(h), removed
}

// isRed
func isRed(h *node) bool {
	return h != nil && h.red
}

// colorFlip
func colorFlip(h *node) {
	//Assert(h != nil)
	//Assert(h.left != nil)
	//Assert(h.right != nil)
	h.red = !h.red
	h.left.red = !h.left.red
	h.right.red = !h.right.red
}

// rotateLeft is a standard tree function
func rotateLeft(h *node) *node {
	//Assert(h != nil)
	//Assert(h.right != nil)
	//Assert(h.right.red == true) // Always red per invariance
	var x *node
	x, h.right, h.right.left = h.right, h.right.left, h
	x.red, h.red = h.red, true
	return x
}

// rotateRight is a standard tree function
func rotateRight(h *node) *node {
	//Assert(h != nil)
	//Assert(h.left != nil)
	//Assert(h.left.red == true) // Always red per invariance
	var x *node
	x, h.left, h.left.right = h.left, h.left.right, h
	x.red, h.red = h.red, true
	return x
}

// moveRedLeft is called during a delete when preparing to move down the tree
// to the left. @h is the node we are moving from, and we want to ensure the new
// left node is not a 2-node, as deleting a 2-node would break black-link balance.
// The caller is expected to ensure that h.left is black and h.left.left is Black
// before calling this function.
//
// This function enforces the invariant: h.left is red or h.left.left is red
func moveRedLeft(h *node) *node {
	//Assert(h != nil)
	colorFlip(h)
	//Assert(h.right != nil)
	if isRed(h.right.left) {
		h.right = rotateRight(h.right)
		h = rotateLeft(h)
		colorFlip(h)
	}
	return h
}

// moveRedRight is called during a delete when preparing to move down the tree
// to the right. @h is the node we are moving from, and we want to ensure the new
// right node is not a 2-node, as deleting a 2-node would break black-link balance.
// The caller is expected to ensure that h.right is black and h.right.right is Black
// before calling this function.
//
// This function enforces the invariant: h.right is red or h.right.right is red
func moveRedRight(h *node) *node {
	//Assert(h != nil)
	colorFlip(h)
	//Assert(h.left != nil)
	if isRed(h.left.left) {
		h = rotateRight(h)
		colorFlip(h)
	}
	return h
}

// fixUp fixes right-leaning reds and consecutive left-leaning reds, and splits 4-nodes
func fixUp(h *node) *node {
	//Assert(h != nil)
	// Fix right-leaning reds
	if isRed(h.right) {
		h = rotateLeft(h)
	}
	// Fix two reds in a row - left and left.left
	if isRed(h.left) && isRed(h.left.left) {
		h = rotateRight(h)
	}
	// Split 4-nodes
	if isRed(h.left) && isRed(h.right) {
		colorFlip(h)
	}
	//Assert(h.right == nil || h.right.red == false)
	//Assert(h.left == nil || h.left.red == false || h.left.left == nil || h.left.left.red == false)
	return h
}
