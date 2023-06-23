package sort

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	t.Run("many chunks", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		err := Sort(strings.NewReader(`this
test
asd
the
end
sad
this
is
asd
end`), 3, buffer)
		require.NoError(t, err)
		assert.Equal(t, `asd
asd
end
end
is
sad
test
the
this
this
`, buffer.String())
	})
	t.Run("one chunk", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		err := Sort(strings.NewReader(`this
test
asd
the
end
sad
this
is
asd
end`), 1e6, buffer)
		require.NoError(t, err)
		assert.Equal(t, `asd
asd
end
end
is
sad
test
the
this
this
`, buffer.String())
	})
	t.Run("one row", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		err := Sort(strings.NewReader(`this
test
asd
the
end
sad
this
is
asd
end`), 1, buffer)
		require.NoError(t, err)
		assert.Equal(t, `asd
asd
end
end
is
sad
test
the
this
this
`, buffer.String())
	})
}
