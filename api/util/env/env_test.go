package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Env_GetEnv(t *testing.T) {
	type state struct {
		key string
	}
	type want struct {
		result string
	}
	type test struct {
		name string
		state
		want
	}

	os.Setenv("BLAH", "Blah!")
	parseEnvVars()

	cases := []test{
		{"env var exists", state{"BLAH"}, want{"Blah!"}},
		{"env var does not exist", state{"BLOOP"}, want{""}},
	}

	for _, c := range cases {
		result := GetEnv(c.state.key)
		assert.Equal(t, result, c.want.result)
	}
}
