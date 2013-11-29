go_llrb
=======

**Extensible Left Leaning Red Black Tree (LLRB) Implementation in Go**


About
-----

Package `llrb` implements a Left-Leaning Red-Black Tree (2-3 Variant),
as described by Robert Sedgewick:

 * http://www.cs.princeton.edu/~rs/talks/LLRB/RedBlack.pdf

 * http://en.wikipedia.org/wiki/Left-leaning_red%E2%80%93black_tree

The goal of this implementation is to present an LLRB as an extensible
Binary Search Tree as defined in the iNamik/go_bst package and sub-packages:

 * https://github.com/iNamik/go_bst


Base Implementation
-------------------

If you are interested in playing around with LLRB's, please see my
base LLRB implementation with a one-to-one correlation to the ideas
and code mentioned in Sedgewick's PDF:

 * https://gist.github.com/iNamik/5844150


Standard BST Methods
--------------------

All of the standard BST methods required to satisfy the
`bst.T` interface have been implemented:

 * Empty
 * ReplaceOrInsert
 * Get
 * Remove


Extensible BST Methods
----------------------

All of the extensible interfaces have been implemented:

 * Find  (see `finder.T`)
 * Visit (see `visitor.T`)
 * Walk  (see `walker.T`)


Additional BST Methods
----------------------

The following additional BST methods have been implemented:

 * Size (see `bst.I_Size`)
 * Min  (see `finder.I_Min`)
 * Max  (see `finder.I_Max`)


Effeciency
----------

Instead of storing parent and sibling information on each node,
this implementation uses recursion to accomplish various tree
navigation functions.  Specifically, the following functions
use recursion:

 * ReplaceOrInsert
 * Remove
 * Visit
 * Walk

The remaining functions do not use recursion and can be
considered efficient implementations.


License
-------

This package is released under the MIT License.
See included file 'LICENSE' for more details.


Contributors
------------

David Farell <DavidPFarrell@yahoo.com>
