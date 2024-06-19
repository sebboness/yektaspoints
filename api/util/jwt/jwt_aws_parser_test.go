package jwt

import (
	"errors"
	"fmt"
	"testing"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
)

var errFail = errors.New("fail")

func Test_AwsJwtParser_GetJwtClaims(t *testing.T) {
	type state struct {
		errParse error
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
		{"fail parser", state{errParse: errFail}, want{"error parsing jwt token"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			jwtParser := &AwsJwtParser{
				KeyFunc: func(token *jwtv5.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(""), c.state.errParse
				},
			}

			mockJwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOi" +
				"IxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5M" +
				"DIyfQ.he0ErCNloe4J7Id0Ry2SEDg09lKkZkfsRiGsdX_vgEg"

			claims, err := jwtParser.GetJwtClaims(mockJwtToken)
			tests.AssertError(t, err, c.want.err)

			if err == nil {
				assert.Greater(t, len(claims), 0)
			}
		})
	}
}
