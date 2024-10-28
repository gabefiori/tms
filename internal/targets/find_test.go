package targets

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	tempDir := t.TempDir()

	dir1 := filepath.Join(tempDir, "dir1")
	dir2 := filepath.Join(tempDir, "dir2")
	dir3 := filepath.Join(tempDir, "dir3")
	subDir1 := filepath.Join(dir1, "subdir1")
	subDir2 := filepath.Join(dir1, "subdir2")

	assert.NoError(t, os.MkdirAll(dir1, 0755))
	assert.NoError(t, os.MkdirAll(dir2, 0755))
	assert.NoError(t, os.MkdirAll(dir3, 0755))
	assert.NoError(t, os.MkdirAll(subDir1, 0755))
	assert.NoError(t, os.MkdirAll(subDir2, 0755))

	symlinkPath := filepath.Join(tempDir, "symlink_to_dir2")
	assert.NoError(t, os.Symlink(dir2, symlinkPath))

	tests := []struct {
		name     string
		rootDir  string
		maxDepth uint8
		expected []string
	}{
		{
			name:     "Depth 0",
			rootDir:  dir1,
			maxDepth: 0,
			expected: []string{dir1},
		},
		{
			name:     "Depth 1",
			rootDir:  dir1,
			maxDepth: 1,
			expected: []string{
				subDir1,
				subDir2,
				dir1,
			},
		},
		{
			name:     "Depth 2",
			rootDir:  tempDir,
			maxDepth: 2,
			expected: []string{
				tempDir,
				dir1,
				subDir1,
				subDir2,
				dir2,
				dir3,
				symlinkPath,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Find(tt.rootDir, tt.maxDepth)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

func BenchmarkFind(b *testing.B) {
	tempDir := b.TempDir()
	tempDir2 := b.TempDir()

	dirs := make(map[string]struct{})

	for i := 0; i < 100; i++ {
		dirPath := filepath.Join(tempDir, "sub_dir_1", fmt.Sprintf("sub_dir_%d", i))
		dirs[dirPath] = struct{}{}
	}

	for i := 0; i < 100; i++ {
		dirPath := filepath.Join(tempDir2, "sub_dir_3", fmt.Sprintf("sub_dir_%d", i))
		dirs[dirPath] = struct{}{}
	}

	for dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			b.Fatal(err)
		}
	}

	symlinkPath := filepath.Join(tempDir, "symlink")
	err := os.Symlink(filepath.Join(tempDir2, "sub_dir_3"), symlinkPath)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("Depth 0", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := Find(tempDir, 0)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Depth 1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := Find(tempDir, 1)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Full Depth", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := Find(tempDir, 3)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
