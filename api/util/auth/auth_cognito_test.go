package auth

import (
	"testing"
)

func Test_CognitoController_Authenticate(t *testing.T) {
	type state struct {
	}
	type want struct {
		err string
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"happy path", state{}, want{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// ctx := context.Background()

			// c, err := NewWithClient(ctx, "", "")
			// assert.Nil(t, err)

			// res, err := c.Authenticate(ctx, username, pw)
			// assert.Nil(t, err)
			// assert.NotEmpty(t, res.Token)
			// assert.NotEmpty(t, res.Username)
			// assert.Greater(t, res.ExpiresIn, 0)

			// username = "sebboness"
			// pw = "sD97$$5L"
			// err = c.UpdatePassword(ctx, res.Session, username, pw)
			// assert.Nil(t, err)
		})
	}
}
