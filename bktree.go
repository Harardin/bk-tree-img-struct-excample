package main

type BKTree struct {
	root *node
}

// We Use static struct to save memory
type Image struct { // Struct can be modify to any struct
	Phash  uint64
	Models []string
}

type node struct {
	entry    Image
	children []struct {
		distance int
		node     *node
	}
}

func (bk *BKTree) Add(entry Image) {
	if bk.root == nil {
		bk.root = &node{
			entry: entry,
		}
		return
	}
	bk.root.addChild(entry)
}

func (n *node) addChild(e Image) {
	newnode := &node{entry: e}
loop:
	d := distance(n.entry.Phash, e.Phash)
	for _, c := range n.children {
		if c.distance == d {
			n = c.node
			goto loop
		}
	}
	n.children = append(n.children, struct {
		distance int
		node     *node
	}{d, newnode})
}

type Result struct {
	Distance int
	Entry    Image
}

func (bk *BKTree) Search(needle Image, tolerance int) []*Result {
	results := make([]*Result, 0)
	if bk.root == nil {
		return results
	}
	candidates := []*node{bk.root}
	for len(candidates) != 0 {
		c := candidates[len(candidates)-1]
		candidates = candidates[:len(candidates)-1]
		d := distance(c.entry.Phash, needle.Phash)
		if d <= tolerance {
			results = append(results, &Result{
				Distance: d,
				Entry:    c.entry,
			})
		}

		low, high := d-tolerance, d+tolerance
		for _, c := range c.children {
			if low <= c.distance && c.distance <= high {
				candidates = append(candidates, c.node)
			}
		}
	}
	return results
}

const (
	m1  = 0x5555555555555555 //binary: 0101...
	m2  = 0x3333333333333333 //binary: 00110011..
	m4  = 0x0f0f0f0f0f0f0f0f //binary:  4 zeros,  4 ones ...
	m8  = 0x00ff00ff00ff00ff //binary:  8 zeros,  8 ones ...
	m16 = 0x0000ffff0000ffff //binary: 16 zeros, 16 ones ...
	m32 = 0x00000000ffffffff //binary: 32 zeros, 32 ones
	hff = 0xffffffffffffffff //binary: all ones
	h01 = 0x0101010101010101 //the sum of 256 to the power of 0,1,2,3...
)

// hammingDistance calculates the hamming distance between two 64-bit values.
// The implementation is based on the code found on:
// http://en.wikipedia.org/wiki/Hamming_weight#Efficient_implementation
func distance(left, right uint64) int {
	x := left ^ right
	x -= (x >> 1) & m1             //put count of each 2 bits into those 2 bits
	x = (x & m2) + ((x >> 2) & m2) //put count of each 4 bits into those 4 bits
	x = (x + (x >> 4)) & m4        //put count of each 8 bits into those 8 bits
	return int((x * h01) >> 56)    //returns left 8 bits of x + (x<<8) + (x<<16) + (x<<24) + ...
}
