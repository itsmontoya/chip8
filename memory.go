package main

type memory [4096]byte

func (m *memory) clear() {
	for i := 0; i < 4096; i++ {
		m[i] = 0
	}
}
