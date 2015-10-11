package helen

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// StaticFile represent the static asset file.
type StaticFile struct {
	rsc  io.ReadSeeker // the underlying ReadSeeker
	info os.FileInfo   // information about the file.
}

// FileChain is an interface for manipulating StaticFile
type FileChain func(*StaticFile) *StaticFile

// Name returns the name of the file.
func (s *StaticFile) Name() string {
	return s.info.Name()
}

// Size returns the size of thefile.
func (s *StaticFile) Size() int64 {
	return s.info.Size()
}

// Mode returns the file mode.
func (s *StaticFile) Mode() os.FileMode {
	return s.info.Mode()
}

// ModTime returns last modification time
func (s *StaticFile) ModTime() time.Time {
	return s.info.ModTime()
}

// IsDir returns true if the file is directory.
func (s *StaticFile) IsDir() bool {
	return s.info.IsDir()
}

// Sys is system stuffs
func (s *StaticFile) Sys() interface{} {
	return s.info.Sys()
}

// Read files p with content from the file
func (s StaticFile) Read(p []byte) (int, error) {
	return s.rsc.Read(p)
}

// Seek seeks some stuffs( I am kidding, help write a better comment)
func (s *StaticFile) Seek(offset int64, whence int) (int64, error) {
	return s.rsc.Seek(offset, whence)
}

// SetReadSeeker changes the underlying readseeker
func (s *StaticFile) SetReadSeeker(rs io.ReadSeeker) {
	s.rsc = rs
}

// NewFromHTTPFile returns  new *StaticFile from f
func NewFromHTTPFile(f http.File) (*StaticFile, error) {
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	b, _ := ioutil.ReadAll(f)
	return &StaticFile{rsc: bytes.NewReader(b), info: info}, nil
}
