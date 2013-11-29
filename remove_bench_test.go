package llrb

import (
	"math/rand"
	"testing"
)

import (
	"github.com/iNamik/go_cmp"
)

type f_bm_remove func(*node, interface{}, cmp.F) (*node, bool)

func (t *tree) BM_Remove(key interface{}, f f_bm_remove) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	var removed bool
	t.root, removed = f(t.root, key, t.fcmp)
	if t.root != nil {
		t.root.red = false
	}
	if removed {
		t.size--
	}
	if t.root != nil && t.size <= 0 {
		panic(t.size)
	} else if t.root == nil && t.size != 0 {
		panic(t.size)
	}
	//	Assert((t.root != nil && t.size > 0) || (t.root == nil && t.size == 0))
	return removed
}

func bm_run(size int, f f_bm_remove, b *testing.B) {
	data := rand.Perm(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r := New(cmp.F_int).(*tree)
		for _, n := range data {
			r.ReplaceOrInsert(n, n)
		}
		b.StartTimer()
		for _, n := range data {
			r.BM_Remove(n, f)
		}
	}
}

func bm_rs_get(h *node, key interface{}, fcmp cmp.F) *node {
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

func bm_rs_min(h *node) *node {
	if h != nil {
		for h.left != nil {
			h = h.left
		}
	}
	return h
}

func bm_rs_deleteMin(h *node) *node {
	if h.left == nil {
		return nil
	}
	if !isRed(h.left) && !isRed(h.left.left) {
		h = moveRedLeft(h)
	}
	h.left = bm_rs_deleteMin(h.left)
	return fixUp(h)
}

/**********************************************************************
 ** Remove Method - RS 1
 **
 ** Method : RS
 ** RemoveNode: RS
 ** Early EQ/R Fail : Present
 **********************************************************************/

func Benchmark_Remove_RS1_1(b *testing.B) {
	bm_run(1000000, bm_remove_rs1, b) //!!!!!!!!!!!!!!!!!!!!
}

func Benchmark_Remove_RS1_2(b *testing.B) {
	bm_run(1000, bm_remove_rs1, b) //!!!!!!!!!!!!!!!!!!!!
}

func bm_remove_rs1(h *node, key interface{}, fcmp cmp.F) (*node, bool) {
	if h == nil {
		return nil, false
	}
	removed := false
	c := fcmp(key, h.key)
	if c == cmp.LT {
		if h.left != nil && !isRed(h.left) && !isRed(h.left.left) {
			h = moveRedLeft(h)
		}
		h.left, removed = bm_remove_rs1(h.left, key, fcmp) //!!!!!!!!!!!!!!!!!!!!
	} else {
		if isRed(h.left) {
			h = rotateRight(h)
			c = fcmp(key, h.key)
		}
		if c == cmp.EQ && h.right == nil {
			return nil, true
		}
		if h.right != nil && !isRed(h.right) && !isRed(h.right.left) {
			h = moveRedRight(h)
			c = fcmp(key, h.key)
		}
		if c == cmp.EQ {
			//if h.right == nil {
			//	return nil, true
			//}

			h.key = bm_rs_min(h.right).key
			h.value = bm_rs_get(h.right, h.key, fcmp)
			h.right = bm_rs_deleteMin(h.right)
			removed = true

			//var m *node
			//h.right, m = removeMin(h.right)
			//h.key, h.value, removed = m.key, m.value, true
		} else {
			h.right, removed = bm_remove_rs1(h.right, key, fcmp) //!!!!!!!!!!!!!!!!!!!!
		}
	}
	return fixUp(h), removed
}

/**********************************************************************
 ** Remove Method - RS 2
 **
 ** Method : RS
 ** RemoveNode: DF
 ** Early EQ/R Fail : Present
 **********************************************************************/

func Benchmark_Remove_RS2_1(b *testing.B) {
	bm_run(1000000, bm_remove_rs2, b) //!!!!!!!!!!!!!!!!!!!!
}

func Benchmark_Remove_RS2_2(b *testing.B) {
	bm_run(1000, bm_remove_rs2, b) //!!!!!!!!!!!!!!!!!!!!
}

func bm_remove_rs2(h *node, key interface{}, fcmp cmp.F) (*node, bool) {
	if h == nil {
		return nil, false
	}
	removed := false
	c := fcmp(key, h.key)
	if c == cmp.LT {
		if h.left != nil && !isRed(h.left) && !isRed(h.left.left) {
			h = moveRedLeft(h)
		}
		h.left, removed = bm_remove_rs2(h.left, key, fcmp) //!!!!!!!!!!!!!!!!!!!!
	} else {
		if isRed(h.left) {
			h = rotateRight(h)
			c = fcmp(key, h.key)
		}
		if c == cmp.EQ && h.right == nil {
			return nil, true
		}
		if h.right != nil && !isRed(h.right) && !isRed(h.right.left) {
			h = moveRedRight(h)
			c = fcmp(key, h.key)
		}
		if c == cmp.EQ {
			//if h.right == nil {
			//	return nil, true
			//}

			//h.key = bm_rs_min(h.right).key
			//h.value = bm_rs_get(h.right, h.key, fcmp)
			//h.right = bm_rs_deleteMin(h.right)
			//removed = true

			var m *node
			h.right, m = removeMin(h.right)
			h.key, h.value, removed = m.key, m.value, true
		} else {
			h.right, removed = bm_remove_rs2(h.right, key, fcmp) //!!!!!!!!!!!!!!!!!!!!
		}
	}
	return fixUp(h), removed
}

/**********************************************************************
 ** Remove Method - RS 3
 **
 ** Method : RS
 ** RemoveNode: RS
 ** Early EQ/R Fail : Removed
 **********************************************************************/

func Benchmark_Remove_RS3_1(b *testing.B) {
	bm_run(1000000, bm_remove_rs3, b) //!!!!!!!!!!!!!!!!!!!!
}

func Benchmark_Remove_RS3_2(b *testing.B) {
	bm_run(1000, bm_remove_rs3, b) //!!!!!!!!!!!!!!!!!!!!
}

func bm_remove_rs3(h *node, key interface{}, fcmp cmp.F) (*node, bool) {
	if h == nil {
		return nil, false
	}
	removed := false
	c := fcmp(key, h.key)
	if c == cmp.LT {
		if h.left != nil && !isRed(h.left) && !isRed(h.left.left) {
			h = moveRedLeft(h)
		}
		h.left, removed = bm_remove_rs3(h.left, key, fcmp) //!!!!!!!!!!!!!!!!!!!!
	} else {
		if isRed(h.left) {
			h = rotateRight(h)
			c = fcmp(key, h.key)
		}
		//if c == cmp.EQ && h.right == nil {
		//	return nil, true
		//}
		if h.right != nil && !isRed(h.right) && !isRed(h.right.left) {
			h = moveRedRight(h)
			c = fcmp(key, h.key)
		}
		if c == cmp.EQ {
			if h.right == nil {
				return nil, true
			}

			h.key = bm_rs_min(h.right).key
			h.value = bm_rs_get(h.right, h.key, fcmp)
			h.right = bm_rs_deleteMin(h.right)
			removed = true

			//var m *node
			//h.right, m = removeMin(h.right)
			//h.key, h.value, removed = m.key, m.value, true
		} else {
			h.right, removed = bm_remove_rs3(h.right, key, fcmp) //!!!!!!!!!!!!!!!!!!!!
		}
	}
	return fixUp(h), removed
}

/**********************************************************************
 ** Remove Method - RS 4
 **
 ** Method : RS
 ** RemoveNode: DF
 ** Early EQ/R Fail : Removed
 **********************************************************************/

func Benchmark_Remove_RS4_1(b *testing.B) {
	bm_run(1000000, bm_remove_rs4, b) //!!!!!!!!!!!!!!!!!!!!
}

func Benchmark_Remove_RS4_2(b *testing.B) {
	bm_run(1000, bm_remove_rs4, b) //!!!!!!!!!!!!!!!!!!!!
}

func bm_remove_rs4(h *node, key interface{}, fcmp cmp.F) (*node, bool) {
	if h == nil {
		return nil, false
	}
	removed := false
	c := fcmp(key, h.key)
	if c == cmp.LT {
		if h.left != nil && !isRed(h.left) && !isRed(h.left.left) {
			h = moveRedLeft(h)
		}
		h.left, removed = bm_remove_rs4(h.left, key, fcmp) //!!!!!!!!!!!!!!!!!!!!
	} else {
		if isRed(h.left) {
			h = rotateRight(h)
			c = fcmp(key, h.key)
		}
		//if c == cmp.EQ && h.right == nil {
		//	return nil, true
		//}
		if h.right != nil && !isRed(h.right) && !isRed(h.right.left) {
			h = moveRedRight(h)
			c = fcmp(key, h.key)
		}
		if c == cmp.EQ {
			if h.right == nil {
				return nil, true
			}

			//h.key = bm_rs_min(h.right).key
			//h.value = bm_rs_get(h.right, h.key, fcmp)
			//h.right = bm_rs_deleteMin(h.right)
			//removed = true

			var m *node
			h.right, m = removeMin(h.right)
			h.key, h.value, removed = m.key, m.value, true
		} else {
			h.right, removed = bm_remove_rs4(h.right, key, fcmp) //!!!!!!!!!!!!!!!!!!!!
		}
	}
	return fixUp(h), removed
}

/**********************************************************************
 ** Remove Method - DF 1
 **
 ** Method : DF
 ** RemoveNode: DF
 ** Early EQ/R Fail : Present
 **********************************************************************/

func Benchmark_Remove_DF1_1(b *testing.B) {
	bm_run(1000000, bm_remove_df1, b) //!!!!!!!!!!!!!!!!!!!!
}

func Benchmark_Remove_DF1_2(b *testing.B) {
	bm_run(1000, bm_remove_df1, b) //!!!!!!!!!!!!!!!!!!!!
}

func bm_remove_df1(h *node, key interface{}, fcmp cmp.F) (*node, bool) {
	if h == nil {
		return nil, false
	}
	removed := false
	c := fcmp(key, h.key)
	if c == cmp.LT {
		if h.left == nil {
			return h, false
		}
		if !isRed(h.left) && !isRed(h.left.left) {
			h = moveRedLeft(h)
		}
		h.left, removed = bm_remove_df1(h.left, key, fcmp) //!!!!!!!!!!!!!!!!!!!!
	} else {
		if isRed(h.left) {
			h = rotateRight(h)
			h.right, removed = bm_remove_df1(h.right, key, fcmp) //!!!!!!!!!!!!!!!!!!!!

		} else if h.right == nil {
			if c == cmp.EQ {
				return nil, true
			}
			return h, false
		} else {
			if !isRed(h.right) && !isRed(h.right.left) {
				h = moveRedRight(h)
				c = fcmp(key, h.key)
			}
			if c == cmp.EQ {
				//h.key = bm_rs_min(h.right).key
				//h.value = bm_rs_get(h.right, h.key, fcmp)
				//h.right = bm_rs_deleteMin(h.right)
				//removed = true

				var m *node
				h.right, m = removeMin(h.right)
				h.key, h.value, removed = m.key, m.value, true
			} else {
				h.right, removed = bm_remove_df1(h.right, key, fcmp) //!!!!!!!!!!!!!!!!!!!!
			}
		}
	}
	return fixUp(h), removed
}

/**********************************************************************
 ** Remove Method - DF 2
 **
 ** Method : DF
 ** RemoveNode: RS
 ** Early EQ/R Fail : Present
 **********************************************************************/

func Benchmark_Remove_DF2_1(b *testing.B) {
	bm_run(1000000, bm_remove_df2, b) //!!!!!!!!!!!!!!!!!!!!
}

func Benchmark_Remove_DF2_2(b *testing.B) {
	bm_run(1000, bm_remove_df2, b) //!!!!!!!!!!!!!!!!!!!!
}

func bm_remove_df2(h *node, key interface{}, fcmp cmp.F) (*node, bool) {
	if h == nil {
		return nil, false
	}
	removed := false
	c := fcmp(key, h.key)
	if c == cmp.LT {
		if h.left == nil {
			return h, false
		}
		if !isRed(h.left) && !isRed(h.left.left) {
			h = moveRedLeft(h)
		}
		h.left, removed = bm_remove_df2(h.left, key, fcmp) //!!!!!!!!!!!!!!!!!!!!
	} else {
		if isRed(h.left) {
			h = rotateRight(h)
			h.right, removed = bm_remove_df2(h.right, key, fcmp) //!!!!!!!!!!!!!!!!!!!!

		} else if h.right == nil {
			if c == cmp.EQ {
				return nil, true
			}
			return h, false
		} else {
			if !isRed(h.right) && !isRed(h.right.left) {
				h = moveRedRight(h)
				c = fcmp(key, h.key)
			}
			if c == cmp.EQ {
				h.key = bm_rs_min(h.right).key
				h.value = bm_rs_get(h.right, h.key, fcmp)
				h.right = bm_rs_deleteMin(h.right)
				removed = true

				//var m *node
				//h.right, m = removeMin(h.right)
				//h.key, h.value, removed = m.key, m.value, true
			} else {
				h.right, removed = bm_remove_df2(h.right, key, fcmp) //!!!!!!!!!!!!!!!!!!!!
			}
		}
	}
	return fixUp(h), removed
}
