package targets

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTargets(t *testing.T) {
	t.Run("Print", func(t *testing.T) {
		tg := New()
		str := "Hello\nWorld\n"

		tg.Buf = bytes.NewBufferString(str)

		buf := new(bytes.Buffer)
		err := tg.Print(buf)

		assert.NoError(t, err)
		assert.Equal(t, str, buf.String())
	})

	t.Run("Filter", func(t *testing.T) {
		tg := New()
		tg.Buf = bytes.NewBufferString("Hello\nWorld\nTest\n")

		tg.Filter("World")
		assert.Equal(t, "World\n", tg.Buf.String())
	})

	t.Run("Collect", func(t *testing.T) {
		tempDir := t.TempDir()

		subDir1 := filepath.Join(tempDir, "subdir1")
		subDir2 := filepath.Join(tempDir, "subdir2")
		subDir3 := filepath.Join(tempDir, "subdir3")

		assert.NoError(t, os.MkdirAll(subDir1, 0755))
		assert.NoError(t, os.MkdirAll(subDir2, 0755))
		assert.NoError(t, os.MkdirAll(subDir3, 0755))

		input := []InputTarget{
			{Path: tempDir, Depth: 1},
			{Path: tempDir, Depth: 2},
		}

		tg := New()

		err := tg.Collect(input)
		assert.NoError(t, err)

		expectedOutput := fmt.Sprintf("%s\n%s\n%s\n%s", subDir3, subDir2, subDir1, tempDir)
		assert.Equal(t, expectedOutput, tg.Buf.String())
	})
}
