package cache

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path"
	"sync"

	"github.com/andybalholm/brotli"
)

type fileCache struct {
	workdir string
	mutex   *sync.Mutex
}

func NewFileCache(workdir string) ICache {
	return &fileCache{
		workdir: workdir,
		mutex:   &sync.Mutex{},
	}
}

func (c *fileCache) Get(key string) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	file, err := os.Open(c.getFilePath(key))

	if err != nil {
		return "", err
	}

	defer file.Close()
	reader := brotli.NewReader(file)
	decompressedBuf := &bytes.Buffer{}
	_, copyErr := io.Copy(decompressedBuf, reader)

	if copyErr != nil {
		return "", copyErr
	}

	return decompressedBuf.String(), nil
}

func (c *fileCache) Set(key string, val string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	filePath := c.getFilePath(key)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755)

	if err != nil {
		return err
	}

	compressor := brotli.NewWriterLevel(file, brotli.BestCompression)
	io.WriteString(compressor, val)

	if flushErr := compressor.Flush(); flushErr != nil {
		return err
	}

	if closeErr := file.Close(); closeErr != nil {
		return closeErr
	}

	return nil
}

func (c *fileCache) Delete(key string) error {
	filepath := c.getFilePath(key)
	_, err := os.Stat(filepath)

	if err != nil {
		return err
	}

	return os.Remove(filepath)
}

func (c *fileCache) Has(key string) bool {
	_, err := os.Stat(c.getFilePath(key))
	return err == nil
}

func (c *fileCache) getFilePath(key string) string {
	hash := md5.Sum([]byte(key))
	filename := hex.EncodeToString(hash[:])
	return path.Join(c.workdir, filename)
}
