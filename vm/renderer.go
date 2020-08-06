package vm

// Renderer will render Graphics output
type Renderer interface {
	Draw(Graphics) error
	GetKeypad() Keypad
}
