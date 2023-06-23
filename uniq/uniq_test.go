package uniq

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"strings"
	"testing"
)

func TestUniq(t *testing.T) {
	t.Run("regular", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		err := Uniq(strings.NewReader(`asd
asd
end
end
is
sad
test
the
this
this
`), buffer)
		require.NoError(t, err)
		assert.Equal(t, `asd 2
end 2
is 1
sad 1
test 1
the 1
this 2
`, buffer.String())

	})
	t.Run("with empty", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		err := Uniq(strings.NewReader(`


asd
asd
end
end
is
sad
test
the
this
this
`), buffer)
		require.NoError(t, err)
		assert.Equal(t, ` 3
asd 2
end 2
is 1
sad 1
test 1
the 1
this 2
`, buffer.String())

	})
	t.Run("one", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		err := Uniq(strings.NewReader(`asd`), buffer)
		require.NoError(t, err)
		assert.Equal(t, `asd 1
`, buffer.String())

	})
	t.Run("one many times", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		err := Uniq(strings.NewReader(`asd
asd
asd
asd
asd
`), buffer)
		require.NoError(t, err)
		assert.Equal(t, `asd 5
`, buffer.String())

	})
	t.Run("empty", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		err := Uniq(strings.NewReader(``), buffer)
		require.NoError(t, err)
		assert.Equal(t, ``, buffer.String())

	})
}

func Test(t *testing.T) {
	open, err := os.Open("/Users/gdanichev/GolandProjects/github.com/yonesko/k-way-merge.git/input_sorted")
	require.NoError(t, err)
	output, err := os.Create(fmt.Sprintf("%s_uniq", open.Name()))
	if err != nil {
		log.Fatal(err)
	}
	err = Uniq(open, output)
	if err != nil {
		log.Fatal(err)
	}
	err = output.Close()
	if err != nil {
		log.Fatal(err)
	}
}
