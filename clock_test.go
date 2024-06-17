package clock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
	assert.InEpsilon(t, 0.02083, Sec(120.0), 0.001)
}
