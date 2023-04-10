package cache

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path"
	"testing"

	"github.com/DAtek/gotils"
	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/assert"
)

func TestFileCacheGet(t *testing.T) {
	t.Run("Returns error if key not cached", func(t *testing.T) {
		defer deleteCacheDir()
		c := newCache()
		val, err := c.Get("1")

		assert.Equal(t, "", val)
		assert.Error(t, err)
	})

	t.Run("Returns value if found", func(t *testing.T) {
		defer deleteCacheDir()
		c := newCache()
		key := "meaning of life"
		value := "42"
		hash := md5.Sum([]byte(key))
		filename := hex.EncodeToString(hash[:])
		filePath := path.Join(cacheDir, filename)

		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			panic(err)
		}
		compressor := brotli.NewWriterLevel(file, brotli.BestCompression)
		io.WriteString(compressor, value)
		compressor.Flush()
		file.Close()

		res, err := c.Get(key)

		assert.Nil(t, err)
		assert.Equal(t, value, res)
	})

	t.Run("Returns error if cached data is not compressed", func(t *testing.T) {
		defer deleteCacheDir()
		c := newCache()
		key := "meaning of life"
		value := "42"
		hash := md5.Sum([]byte(key))
		filename := hex.EncodeToString(hash[:])
		filePath := path.Join(cacheDir, filename)

		file := gotils.ResultOrPanic(os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755))
		file.WriteString(value)
		file.Close()

		_, err := c.Get(key)

		assert.Error(t, err)
	})
}

func TestFileCacheHas(t *testing.T) {
	t.Run("False if file not exists", func(t *testing.T) {
		c := newCache()

		assert.False(t, c.Has("something"))
	})

	t.Run("False if file not exists", func(t *testing.T) {
		c := newCache()
		defer deleteCacheDir()
		key := "meaning of life"
		hash := md5.Sum([]byte(key))
		filename := hex.EncodeToString(hash[:])
		filePath := path.Join(cacheDir, filename)

		f, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		f.Close()

		assert.True(t, c.Has(key))
	})
}

func TestFileCacheDelete(t *testing.T) {
	t.Run("Returns error if key not cached", func(t *testing.T) {
		defer deleteCacheDir()
		c := newCache()

		err := c.Delete("1")

		assert.Error(t, err)
	})

	t.Run("Deletes file", func(t *testing.T) {
		defer deleteCacheDir()
		c := newCache()
		key := "meaning of life"
		hash := md5.Sum([]byte(key))
		filename := hex.EncodeToString(hash[:])
		filePath := path.Join(cacheDir, filename)

		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			panic(err)
		}
		file.Close()

		c.Delete(key)

		_, err = os.Stat(filePath)

		assert.True(t, os.IsNotExist(err))
	})
}

func TestFileCacheSet(t *testing.T) {
	t.Run("Item is being cached", func(t *testing.T) {
		defer deleteCacheDir()
		c := newCache()
		key := "beer"
		value := "lager"

		gotils.NilOrPanic(c.Set(key, value))

		hash := md5.Sum([]byte(key))
		filename := hex.EncodeToString(hash[:])
		filePath := path.Join(cacheDir, filename)
		file, _ := os.Open(filePath)
		reader := brotli.NewReader(file)
		decompressedBuf := &bytes.Buffer{}
		io.Copy(decompressedBuf, reader)

		assert.Equal(t, value, decompressedBuf.String())
	})

	t.Run("Returns error if can't open file", func(t *testing.T) {
		c := NewFileCache("/var/this_directory_does_not_exists")

		assert.Error(t, c.Set("key", "value"))
	})
}

var tmpDir = os.Getenv("TMP_DIR")

var cacheDir = path.Join(tmpDir, "cache")

func newCache() ICache {
	os.Mkdir(cacheDir, 0755)
	return NewFileCache(cacheDir)
}

func deleteCacheDir() {
	gotils.NilOrPanic(os.RemoveAll(cacheDir))
}
