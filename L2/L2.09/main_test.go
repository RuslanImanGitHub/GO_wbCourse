package main

import "testing"

func TestUnpackSequenceV2_TabelDriven(t *testing.T) {
	testCases := []struct {
		name          string
		data          string
		expected      string
		expectedError bool
	}{
		{"No unpacking needed", "abc", "abc", false},
		{"Unpack one digit", "a4", "aaaa", false},
		{"Unpack multiple digits", "a11", "aaaaaaaaaaa", false},
		{"Error all digits", "45", "", true},
		{"Empty input", "", "", false},
		{"Shielded digits", `qwe\4\5`, "qwe45", false},
		{"Unpack shielded digit", `qwe\45`, "qwe44444", false},
		{"Input start with digit", `1qwe\45`, "", true},
		{"Digit one", `q1we`, "qwe", false},
		{"Shielding chars", `\q\w\e`, "qwe", false},
		{"Shielding shielding", `\\`, `\`, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := UnpackSequenceV2(tc.data)
			if result != tc.expected {
				t.Errorf("UnpackSequence(%s) = %s; want %s",
					tc.data, result, tc.expected)
			} else if err != nil && !tc.expectedError {
				t.Errorf("UnpackSequence(%s) returned error %v; error expected %t",
					tc.data, err, tc.expectedError)
			} else if tc.expectedError && err == nil {
				t.Errorf("UnpackSequence(%s) didn't returned error; error expected %t",
					tc.data, tc.expectedError)
			}
		})
	}
}
