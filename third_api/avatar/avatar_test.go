package avatar

import (
	"testing"

	"github.com/fzf-labs/fpkg/util/validutil"
	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {
	assert.True(t, validutil.IsURL(URL()))
}
