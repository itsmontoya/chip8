package main

type graphics [64 * 32]byte

func (g graphics) forEachDelta(in graphics, fn func(index int, val byte)) {
	for i, val := range g {
		if val == in[i] {
			// Values are the same, no drawing needed
			continue
		}

		fn(i, val)
	}
}
