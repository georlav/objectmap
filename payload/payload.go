// Package payload interface
package payload

// Payloads slice contain payload objects
type Payloads []Payload

// Payload interface
type Payload interface {
	Input() [][]byte
	Output() []byte
	Description() string
}
