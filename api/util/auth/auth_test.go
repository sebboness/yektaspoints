package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_computeSecretHash(t *testing.T) {
	secretHash := computeSecretHash("john.smith", "123", "456")
	assert.NotEmpty(t, secretHash)
}

func Test_ValidatePassword(t *testing.T) {
	type state struct {
		pw string
	}
	type want struct {
		missingLength  bool
		missingNumber  bool
		missingLower   bool
		missingUpper   bool
		missingSpecial bool
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"fail - min length", state{"Test!23"}, want{missingLength: true}},
		{"fail - number", state{"TTTT!test"}, want{missingNumber: true}},
		{"fail - lower", state{"TEST123!"}, want{missingLower: true}},
		{"fail - upper", state{"test123!"}, want{missingUpper: true}},
		{"fail - special", state{"Test1231"}, want{missingSpecial: true}},
		{"fail - max length", state{`Test!6789ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF
1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF
1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF
1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1`}, want{missingLength: true}}, // 257 chars here
		{"happy path 1", state{"Test123!"}, want{}},
		{"happy path 2", state{"Test123456$"}, want{}},
		{"happy path 3", state{"Test123456."}, want{}},
		{"happy path 4", state{"TTTt!123"}, want{}},
		{"happy path 5", state{"Tes........12"}, want{}},
		{"happy path 6", state{"Te()?\"!@#%&/\\3"}, want{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			r := ValidatePassword(c.state.pw)

			assert.Equal(t, !c.want.missingLength, r.WithinLength, "minimum length")
			assert.Equal(t, !c.want.missingLower, r.Lower, "lower case character")
			assert.Equal(t, !c.want.missingNumber, r.Number, "digit")
			assert.Equal(t, !c.want.missingSpecial, r.Special, "special character")
			assert.Equal(t, !c.want.missingUpper, r.Upper, "upper case character")
		})
	}
}
