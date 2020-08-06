package vm

// Graphics represents the system graphics
type Graphics [64 * 32]byte

// ForEachDelta will iterate over all the pixels which changed since the last frame
func (g *Graphics) ForEachDelta(in Graphics, fn func(index int, val byte)) {
	for i, val := range g {
		if val == in[i] {
			// Values are the same, no drawing needed
			continue
		}

		fn(i, val)
	}
}

func (g *Graphics) setAllTo(val byte) {
	for i := range g {
		g[i] = val
	}
}

func (g *Graphics) clear() {
	g.setAllTo(0)
}
