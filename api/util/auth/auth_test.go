package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_computeSecretHash(t *testing.T) {
	secretHash := computeSecretHash("john.smith", "123", "456")
	assert.NotEmpty(t, secretHash)
}
