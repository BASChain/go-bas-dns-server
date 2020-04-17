// Code generated by go-bindata.
// sources:
// LICENSE
// DO NOT EDIT!

package config

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _license = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xdc\x5a\x5f\x73\xdc\x46\x72\x7f\xd7\xa7\xe8\x6c\x55\x2a\x64\x15\xb4\xd2\x39\x77\x49\xce\x7e\xa2\x45\xea\xbc\x89\xbc\x54\x91\xab\x28\x2e\x97\x1f\x66\x81\xc6\x62\xa2\xc1\x0c\x3c\x33\xe0\x12\xf9\xf4\xa9\xee\xf9\x83\xc1\xee\x52\x56\x2a\x6f\xe7\x07\x97\x48\x02\x3d\x3d\xfd\xe7\xd7\xbf\xee\xc6\x2b\xf8\xa3\xff\x6e\x06\x51\x77\x08\x1f\x64\x8d\xda\xe1\xd7\x9e\xff\x4f\xb4\x4e\x1a\x0d\xdf\xad\xdf\x56\xf0\xef\x42\x8f\xc2\x4e\xf0\xdd\xdb\xb7\x7f\x7e\xf1\xa5\xce\xfb\xe1\xfb\x37\x6f\x8e\xc7\xe3\x5a\xf0\x31\x6b\x63\x0f\x6f\x54\x38\xca\xbd\x79\x45\x2f\xee\xee\x1e\x7e\x7e\x84\x9b\xed\x2d\xbc\xbb\xdf\xde\x6e\x76\x9b\xfb\xed\x23\xbc\xbf\x7f\x80\x4f\x8f\x77\x15\x3c\xdc\x7d\x7c\xb8\xbf\xfd\xf4\x8e\x7e\x5d\xf1\x53\xb7\x9b\xc7\xdd\xc3\xe6\xc7\x4f\xf4\x1b\x16\xf0\xa7\x35\xdc\x62\x2b\xb5\xf4\xd2\x68\xb7\x7e\x15\xb5\x59\xc5\x1b\xad\xc0\x75\x42\x29\xe8\x51\x68\xf0\x1d\x82\x47\xdb\x3b\x10\xba\x81\xda\xe8\x26\xbc\x05\xad\xb1\x30\x3a\xac\xc0\xe2\x60\x4d\x33\xd6\xf4\xeb\x2a\x8a\xa2\x67\x1b\xe9\xbc\x95\xfb\x91\x7e\x0f\xc2\x41\x43\x47\x62\x03\xfb\x09\x1e\xb1\x0e\x42\xfe\x04\xbe\xb3\x66\x3c\x74\xf0\x57\x30\x2d\xf8\x4e\x3a\x68\x4c\x3d\xf6\xa8\xfd\xa9\x5e\xc6\x9e\x29\x56\x9b\x61\xb2\xf2\xd0\x79\x30\x47\x8d\x16\x8c\x05\xd4\x5e\xfa\x09\xc4\xe8\x3b\x63\xe5\xff\xf0\x79\x51\xce\xa5\x37\x7c\x27\x3c\x48\x07\x07\x2b\xb4\x97\xfa\xc0\x0f\x45\x3b\x14\x0a\xe0\x41\x28\xb8\x63\xd1\x67\x4a\x8c\x9a\x2e\xc8\xda\x23\x88\x9a\xa5\x24\x2d\x74\x03\x42\xa9\x28\xc6\xf8\x0e\xa3\x82\x12\x5d\x38\xba\x36\xda\x5b\xa3\x2a\x10\x16\xd3\x0f\x8a\x95\xae\xe8\x36\xf4\xdb\x51\x37\x68\xa1\x36\x7d\x6f\x74\x94\x14\x1f\x84\xa3\xf4\x5d\x90\x13\x0e\x5c\xc3\x7b\x63\x59\x8f\x61\xb4\x83\x71\xe8\x66\xab\x66\x87\x27\x1f\xad\xa2\x94\x15\x5f\xc5\xc1\x95\xbc\x0e\xaf\x9a\x23\xda\x0a\x1a\x69\xb1\xf6\xa4\x84\xd4\xe1\xdf\x15\x78\x03\xb5\x18\x1d\xd2\x73\x51\x4a\xf8\x13\x5b\xc0\x42\x2f\xb4\x38\x20\x39\x8f\xce\x75\x63\xdd\x45\xc5\x2a\x38\x76\xc8\xd7\xdf\x4f\x41\x7b\xc1\xb2\x4b\xcb\x1c\x25\x45\x93\xb1\x70\x25\xe5\x75\x70\x8f\xeb\xe4\x40\x92\x5a\xd9\xfa\x09\x06\xb4\x35\x89\xbe\xfa\xcb\xdb\x7f\xbc\xe6\xe3\x8c\xc5\x68\xf8\x24\x68\xf4\xce\x0b\xdd\x90\x0f\x5c\x27\x2c\xba\x24\x51\x5e\xc3\x1e\x35\xb6\xb2\x96\x42\x2d\xa5\x17\x7a\xce\x2e\xff\xc5\x8c\x2b\xb8\x32\x96\xff\x65\x57\xd7\xa5\xd7\x85\x66\x9b\x3c\xc9\x66\x24\x59\x16\xca\xf8\x88\x02\xf0\x19\x6d\x2d\x1d\x29\x32\xa0\xed\xa5\x73\x1c\xf0\x1c\x67\x21\x09\xd8\x2d\x67\xa1\xf6\x68\x46\x5b\xe3\x8a\xd2\xab\x3f\x8d\xb4\xc1\x62\x8b\xd6\x62\x13\xfe\xda\xb2\xc5\xbf\xd0\x11\xbd\x69\x64\x2b\x6b\xc1\x59\x95\x1c\x2c\x75\xad\x46\x36\xc5\x7e\xf4\xa0\x8d\x07\x25\x7b\x49\xa7\x7b\x03\xce\xb4\xfe\x48\xe1\xe5\xf8\x40\xa8\x4d\x83\x55\xce\x3d\x16\x14\xc5\x84\x07\xaa\x94\xff\xad\x3c\x8c\x96\xff\x0e\xad\x54\x58\xc0\xc7\xfd\xfe\xbf\xb1\xf6\xe7\xaa\x0b\x3d\x85\xdf\x59\x74\xa3\xe2\xfc\x68\xad\xe9\xa1\xc7\xba\x13\x5a\xd6\x22\x25\x88\xb7\x42\x3b\x7a\x52\xa4\x80\xe2\xdf\xa8\xf8\x63\x0b\x02\x82\x79\x58\x5c\xb5\xbc\x60\x94\x71\x72\xcd\xda\xf4\x83\xa4\x84\x32\xac\x5c\xbc\xe6\x01\x35\x5a\x41\x8f\x2c\x2e\x5c\xa2\x57\x6d\xf4\x53\x40\x6f\x47\x72\x42\xee\xf6\xd8\x48\x01\x7e\x1a\xca\x6b\x7f\x36\xf6\xcb\x19\x28\x1c\x8d\xfd\xc2\x1a\x33\x0e\x51\xa4\xcd\x29\x20\x75\xba\x46\x4e\x80\x60\xba\x78\xad\x5e\x34\x08\xe2\x49\x48\x25\xf6\x2a\xe5\x7f\x81\x4b\x15\xa1\x29\x05\x60\x2d\x62\x28\x89\x8c\x0b\x09\xdd\xb4\xf1\xb2\xc6\x0c\x6f\xc1\x52\xd8\xd0\xd9\x04\x2b\xde\x53\x6d\x61\x0b\x25\x6d\xa3\x88\x2b\xa1\x01\x9f\x45\x3f\x28\xa4\x17\x07\x6b\x9e\x64\x7c\x91\x9e\xbc\x19\x06\xd4\x8d\x7c\x86\x3d\x2a\x73\xbc\x9e\xad\x70\x8b\x56\x3e\x09\x2f\x9f\x10\xc8\x20\x6e\x75\x1a\x01\x74\xc6\x65\x1b\xc4\xdb\x47\x49\xc1\x06\x49\xf1\xbd\x70\xe4\x3c\xcd\xa9\xd8\xd0\x19\x14\xfd\xd6\xf4\x01\xab\xe8\x28\x76\x17\xe5\xc2\xb1\x93\x75\x57\x80\x01\x36\xd2\x1b\x4b\xe9\x6e\xf1\x49\xb2\x2b\x29\x8a\xb5\xf1\x31\x4f\x00\x95\xd8\x1b\x9b\x7e\x32\x36\xb9\xb9\xcc\xa6\x28\x8c\xaa\x1c\x3a\xd4\x9e\xad\x2f\xe0\xd8\x19\xc5\x49\x01\xc6\xca\x83\xd4\x42\x5d\xf0\xf9\x39\x1e\x27\x9c\x6a\x17\xe9\x5f\xc1\xa9\xf9\xa2\xf5\x28\x9a\xa3\xef\x58\x7c\xac\x1a\x16\x7b\x21\x73\x7e\xe2\x20\x2c\x47\x0a\xd9\x85\xaf\xd1\xa3\x45\x35\x81\x92\xfa\x0b\x1b\x6e\x2f\x35\xc7\x89\x16\x3d\x5e\x27\xa7\x4b\xed\xd1\xb6\xa2\xe6\x22\x51\x15\x35\x32\x1b\xf5\x4c\x29\xb2\x0e\x9a\x76\xf6\xfa\x3b\x82\xf2\x58\xe3\x2f\x7a\xfc\x34\x07\x72\xca\x16\xe7\x65\x03\xc6\x84\x4b\xb5\x34\xeb\x41\xc2\x16\x3e\xe1\x18\x6e\x22\x13\x49\x92\x4c\xb0\x0d\xbf\x65\xec\x8b\xca\x57\x45\x52\x78\x42\x7d\xa3\x85\x52\x09\xb6\xdd\xb8\xef\xa5\x8f\xe0\x91\x78\x07\x47\x17\x6b\xce\xea\xc5\x54\xe0\x83\x18\xc7\xcf\x68\x45\xf2\x32\x97\xbb\xaf\x56\x8b\x92\xa8\x10\x2a\xf3\xf1\x14\xef\x7b\xec\x84\x6a\xc1\xb4\x2f\x93\x97\x6f\xab\xf6\xb0\xca\x77\x5a\x45\x59\xa1\xde\x67\x58\x36\x2d\xa0\xc2\xda\x5b\xa3\x65\x5d\x91\x17\xf6\x42\x71\x1c\x1d\x2d\xbd\xa7\x99\x7c\x8c\x3a\x5a\x1f\x28\x0b\x4a\xa3\xe3\x6c\x28\xb2\x93\x77\x73\xb2\xb0\xfd\x5d\xf5\xd5\x52\x94\xb1\xab\x3c\xc3\xe8\x42\x27\xe8\x85\x54\xf4\xb2\x92\xce\xbb\xaa\x2c\x59\x99\x0a\xb9\xc9\x79\xec\x5d\x09\xe1\xd2\xb9\x11\xa9\x84\xd4\x5c\x23\xe3\x13\xc1\xfd\x54\xf9\x02\x5b\xc9\x5c\xab\x34\x7a\x55\xc0\xc8\x22\x0a\x0a\x6b\x93\xdd\x1a\xe9\xea\xd1\x71\x95\xe7\x13\x7b\xc6\xcb\x48\x23\x3f\x33\xe2\xcd\xa5\x09\x9f\x93\x11\x96\x77\x4d\xf1\x58\x1b\xed\x06\x59\x8f\x66\x74\x6a\x82\x5e\xd8\x2f\x04\x7d\x76\x66\x47\x89\x72\xa1\x93\x07\xcd\xd8\x2f\x35\xfb\x88\x0d\x7b\x31\x12\x09\xac\x56\x5b\xe3\x41\x40\x99\xab\xeb\xd5\x79\x0a\x9f\xf0\xeb\x7c\xed\x94\x81\x7f\x48\x79\x4a\x03\x12\x3e\xf6\x27\x87\x42\x27\x1c\xec\x11\x35\x58\xac\x91\x91\x7c\x3f\x2d\xce\x99\x93\xd0\xe1\xef\x23\x6a\xaf\xe8\xd8\xda\xd8\xc1\x84\x72\x4d\x84\xb7\x48\xbf\x00\x44\xdf\xad\xe1\x6f\x44\xab\xe8\xd8\x77\xf9\xfa\x89\x59\xc1\xe3\x18\x8a\x6b\x8c\xd5\x8b\xcd\x4c\x91\x66\x25\x2a\xa3\xa8\x3b\x28\x0c\x04\x04\x21\xfb\x29\xb0\x38\xe6\x05\xbf\x98\x11\x04\x31\xbc\x01\xfd\x28\x54\x0a\xbf\xa3\xb1\xaa\x39\x4a\xe2\x1a\xda\xe8\xd7\xec\x79\x27\x9f\xf8\xc7\xd7\x75\x27\xec\x81\x1a\x27\x33\x09\xe5\xa7\xd7\xad\x45\xac\x40\x5a\x8b\x4f\xa6\x26\x20\x3f\xab\xe6\xb1\xff\xa3\x03\x53\xb7\x85\x15\xd1\xc1\x81\xe2\xf8\x0c\xe9\x66\x38\x1f\xc6\xbd\x92\xb5\x9a\x28\x50\x07\x25\xa6\x6a\xfe\xcd\x80\x36\x94\x5a\xc7\xbf\x89\xc4\xa2\xec\xdb\x4a\x9a\x9f\xb1\x98\xc9\xf2\xd9\x89\x17\xca\x39\x63\x4b\x70\xd0\x3f\x17\x0e\xfa\x28\x08\x74\xff\x0e\xbc\x73\x85\xcf\x35\x0e\x9e\x12\xcc\xf9\x94\x8c\xac\xa0\x0b\x0d\xd1\x35\x0c\xe1\xae\x85\xf7\x7a\xf1\x05\x2b\xe8\xc4\x13\x32\xcb\x4b\x0a\x71\x1f\x6d\xda\x96\x78\x9e\x01\x87\x4a\x55\xf1\xff\xb2\x1f\x8c\xf5\xc1\x31\x19\x07\x22\x51\x8e\xac\x90\x61\x26\xdd\x8c\x4c\x10\x7c\x94\x4e\x15\xc3\xa0\xa8\xdd\x34\x5a\x4d\xc1\xca\x84\x5d\x51\xb5\x5a\x09\xd9\xbb\xf8\x6c\x71\xb9\xfd\x14\x84\x94\xd6\xcd\xb8\xa9\xb1\x46\xe7\x84\x95\x9c\x9d\xad\x95\xfa\x90\x3a\x1a\x94\xa9\xf6\x95\x89\x7f\xe5\xae\x41\x28\xa3\x31\x56\xc4\xda\xf4\x7b\xa9\x33\xab\xe7\xd7\x4e\x5f\x48\x17\x0a\x1d\x6e\xac\xb6\xde\x44\x92\xb7\x54\x2e\x1e\x71\x24\x57\xa4\x5a\xb7\x86\x4d\x4b\xfe\xcf\xbd\x90\xf3\xd2\x53\x4c\x67\xa7\x78\x79\x08\x2a\x88\x83\xa0\x3f\x33\xc8\xc5\xc6\xfd\x6a\x2e\x58\x99\x5b\x5b\xe3\xdc\x6b\x36\x18\x5d\xa3\x36\x23\xf1\xa7\xf0\xb3\xd4\x20\x40\x89\xa3\x1b\xa5\xa7\xab\x2a\x3c\x84\x22\x20\x7c\x56\x7e\xe6\x04\x27\xa8\xf8\x35\x80\xe3\x9a\x10\x14\x77\xb1\xd5\x9e\xe5\xd4\xb3\x73\xa6\x74\xad\xe4\x8f\x9e\x99\xaa\xef\x30\x50\xb1\x65\x24\x26\xca\x94\x9a\xd1\x98\x29\xa9\xd1\x98\x73\x2c\x96\xbc\xc4\xaa\x42\x75\xa0\x14\x25\xef\xa5\x58\x11\x2e\x11\xb6\x46\xf8\x1c\x7c\xd9\xba\xd2\x71\x9f\xd8\x04\x28\xf8\xf3\x1a\x1e\xb0\x9c\x0c\xad\xf9\xe8\x5e\x4c\x33\xb2\x9d\xa2\x50\x6d\x06\x99\xb8\xcd\x02\x8f\xbe\xc2\xf2\xd8\x25\x44\x1b\xb1\x91\x63\x5f\x85\x38\x22\x46\x23\x7d\x67\x72\x45\x5e\xb6\xcd\xa1\x84\xbf\x80\x64\xd5\xdc\x0a\xb1\x41\xe6\xd0\xea\x11\x83\x97\x5b\xa3\x94\x39\x86\xfa\x9e\xb0\xeb\xfb\x57\xb9\xaf\xba\x0e\x37\x1d\x9d\x87\x03\xe9\x4b\xea\x85\x7e\xc3\x62\x2d\x07\x89\x04\x5a\x25\xf5\xcd\xdd\x21\xfd\x77\x76\x51\xc1\xf5\xe1\xb4\x93\xf8\x81\xcb\x68\x3a\x73\x5f\x9c\x19\x06\x37\x33\x95\xa6\x3e\x8a\xfa\xf7\x30\xd4\xb1\x14\x42\xd6\xf4\x52\x53\x9c\x84\xee\xd1\x15\xc7\x13\xc4\xe5\x90\x26\x99\xd4\xba\x1f\xd8\x18\x18\xe4\x2c\x4f\xae\x8b\x93\x2d\x7a\x21\x75\x95\x78\x73\xd1\xc2\x73\x77\xa0\xa7\xb3\xcb\x15\x07\xe7\x03\xe7\x80\xa8\x28\xc3\xe6\xea\x58\xc5\xe8\xae\x08\x16\x1b\x24\xde\x54\x15\x64\x82\x43\xd4\xcf\xe9\x16\xef\x16\x46\x10\x17\xf4\x39\x85\xd4\x25\x73\x0b\xe8\x99\x64\xb0\x72\x8d\x61\x42\x3b\xa0\xa5\x6b\x92\x39\x43\xc6\x59\x3f\x17\x2e\x88\x0c\xfe\xf4\xa2\x4b\xa3\x35\xd7\x04\x5a\xd9\xff\xb1\xf1\x23\x57\xaf\xb6\xf7\xbb\xcd\xbb\xbb\x15\x78\x7c\xf6\x6c\x6f\x4a\xbb\x78\x06\x51\xee\xe2\x9c\x32\xbb\x0a\x08\xb8\x90\x29\x67\x96\x65\x7f\x15\xa2\x52\xeb\x29\xc0\xa2\x68\xb8\xc7\x9c\x83\x0e\x2f\x9a\x95\x40\x49\x48\x8d\xa5\xf9\x23\xa8\x31\x32\x84\x8b\xf0\x15\xaa\x6f\xb1\x6b\x21\xe6\xb2\x85\x2f\xda\x95\x83\x4d\x78\x50\x28\x1c\xb5\x53\xe5\x94\x3e\xbe\x32\x67\xeb\xa0\xa8\x09\xfe\x3e\xa9\x29\x92\x8e\xb3\xad\x67\x0b\x2d\xa2\xca\x7d\x55\x87\x1f\x4a\x30\x5f\x04\x59\x99\xd7\xcb\x01\x14\xc8\x76\xc6\x19\x2a\x99\x87\xb9\x02\x9e\xcb\x37\xb6\x3a\xb7\xb2\x48\x5c\xaf\x98\x72\xc5\xde\xe0\x82\x95\xda\x93\x4c\x61\x02\xf1\x84\x36\x38\xcb\x77\xd2\x36\xaf\xe9\x92\x53\xf6\x8d\x36\xb6\xa7\x86\x99\x88\x05\x0a\xbb\x86\x5d\x17\xba\x30\xc2\xaf\x73\x33\x17\xfe\x66\xf2\x10\x5a\xe9\x3c\xe4\x13\xaa\x68\x5e\x89\xa1\x2c\xd5\x89\xb9\xc5\x88\x35\x2d\x66\xf3\xb9\x6c\x88\xa6\xa1\x7f\x5b\xea\x77\xca\x88\x2c\xa4\x24\xd5\xa3\x85\xbe\x25\x13\xaa\x60\x7d\x27\x9b\x45\xe8\x70\x3f\x25\x34\x1d\x8a\xba\x19\xfb\x44\x5b\x17\x11\x93\x80\x25\xf4\x7f\xc9\x9d\xa7\x98\xc6\x06\x4e\x43\x0c\xa1\x2e\x27\x13\x4f\xab\x60\x8f\x81\x07\xd8\xf1\x34\xfe\x82\x61\x5e\xda\x5b\x5c\x34\xd1\xdc\x55\x30\x6d\xe5\x61\x7d\x20\x00\x27\x83\xaf\xc2\x15\x24\x24\xde\xa3\x54\xd9\x58\x68\x24\xb1\xd6\x05\xcb\xbd\xc0\xe0\xe7\xd1\xde\x85\x95\x51\x10\x53\xec\x8a\x4c\x7b\x41\x9b\x6a\x4e\x9b\x96\x9b\xc5\xe9\x85\x56\xa4\x9c\xce\xe5\x54\x62\x79\x74\x74\x31\xcd\x9b\x15\x38\xdb\x56\x2d\xaa\x70\x66\xdd\xb5\xe9\x03\x95\xa6\x38\x5a\x8c\x65\x72\xa7\x72\xd2\x09\x2c\x1c\xf2\x17\x6e\x76\xe2\x26\x20\xf4\xaa\x33\x0b\x74\x6b\xf8\xa4\x15\x3a\xc7\x4e\xc3\xe7\x41\xc9\x5a\x52\xfb\xcb\x12\x8b\x05\x49\x9e\x6f\x4c\xa7\x2c\xb2\x18\x66\x15\x63\xac\x17\x47\x57\x33\xd3\xa7\x13\x4f\x07\x39\x81\xea\xed\xcb\xe9\xf3\xff\xa5\x35\x8b\x34\x8b\xd5\x2c\x02\x26\x88\x08\xd4\xb5\x49\xdb\xc7\xf0\xfe\xd6\x78\x7a\x29\x6f\x6f\xb8\xbe\xec\x4d\x68\xca\x28\x6d\x0f\xdc\xde\x51\x19\x61\xd5\xdc\x38\xa0\x75\xd8\x60\x58\x04\x51\x1a\x14\x2e\x89\x07\x05\x76\x11\x06\xa4\x1e\xe7\x96\xe8\x60\x31\x04\xfe\x14\x33\x84\x3b\x32\x7c\xc6\xba\x80\x78\x06\xde\x6c\x10\x8b\x07\x61\xc3\x5e\xe9\xb4\xf7\x88\xbb\x80\x7f\x59\xc3\x2e\x11\x10\x47\xb0\x58\xf0\xe8\xc6\x30\x72\xfa\x40\xb9\x8b\x8d\x10\x19\x3e\x2e\xd4\x02\x7d\x49\x6b\x0c\xd1\xa3\x2b\x18\x8d\xa3\x86\xd0\x3e\xc9\x1a\x21\xfe\x68\x2c\xc4\x18\x0e\x0f\xa7\xa0\x4d\x1a\x57\xf3\xd4\x29\xb6\xa9\x16\x7f\x1f\x65\xdc\x1e\x51\x41\x77\x46\x73\x49\x67\x97\x8e\xce\x9b\x5e\xd8\x89\xb5\x91\x1a\x1a\x74\xb5\x95\xfb\xe8\x8a\xdc\x74\xc8\x83\x3c\x9f\xcf\xa6\x6c\x4a\x7e\x8b\xd5\xe0\x42\x09\x08\x96\xfa\xd7\x35\xdc\x4a\xc7\xad\x13\x5a\x7a\xea\xb3\xb0\x64\x97\x29\x27\x41\x56\x75\x3f\x85\x06\x96\x3b\x6f\x6a\xb1\x66\x18\x60\x2f\x72\xf3\x32\x4f\xc1\xaa\xd9\x61\x31\xf7\xdd\xac\xea\x15\xe9\x8a\xa2\xee\x4e\x5b\xd4\xf2\x69\xe9\xdd\xd2\xb9\xd7\x60\x78\xe3\xb7\xba\x79\x84\xcd\xe3\x0a\x7e\xbc\x79\xdc\x3c\x26\xe3\x7e\xde\xec\x7e\xba\xff\xb4\x83\xcf\x37\x0f\x0f\x37\xdb\xdd\xe6\xee\x11\xee\x1f\xca\xb5\xfc\xfd\x7b\xb8\xd9\xfe\x02\xff\xb1\xd9\xde\x56\x80\x32\x6c\x80\x9f\x07\x4b\x97\xcc\x37\x91\x8c\x2b\x4d\x31\x26\x9d\x33\x88\xe7\xa4\x22\xe1\xd4\x04\xc7\x60\x2a\x6e\x88\xec\x39\xc4\x9a\x16\x76\x9b\xdd\x87\xbb\x0a\xb6\xf7\xdb\xd7\x9b\xed\xfb\x87\xcd\xf6\x6f\x77\x3f\xdf\x6d\x77\x15\xfc\x7c\xf7\xf0\xee\xa7\x9b\xed\xee\xe6\xc7\xcd\x87\xcd\xee\x17\x0e\xa1\xf7\x9b\xdd\xf6\xee\x31\x7c\x3e\x70\x13\x65\x7c\xbc\x79\xd8\x6d\xde\x7d\xfa\x70\xf3\x00\x1f\x3f\x3d\x7c\xbc\x7f\xbc\x0b\xd5\x36\x6c\x0b\x15\x2a\xea\xd5\xdc\x60\xb4\x93\xbc\x75\xe0\xcd\x4c\xe8\x0a\x97\xe1\x22\x86\xc1\x9a\xc1\x4a\xa2\xe7\x7c\xe1\x16\x46\x9e\x95\x72\xfc\xcd\x88\x5b\xcc\x4b\xc3\xb4\xd1\xb9\xb1\xe7\x5e\x25\xc1\xb5\x74\x8c\xec\xce\xd4\x32\xb7\xc9\x01\xd4\xe3\x9e\x95\xa7\xb1\xe5\xa2\xf5\xbc\x99\x0d\xb1\xf7\x6f\x6b\xf8\x90\x4d\x4a\x2f\x7d\x90\x62\x2f\x15\x2f\xcf\x37\x54\x79\x01\x9f\x28\x76\x49\x8f\x20\x43\x1b\x50\x3c\xec\xf4\x1d\x1a\x3b\x15\xa3\x96\xb4\xc9\xf2\xc6\xfa\x72\x64\xa0\xf1\xa0\xe4\x01\x75\x8d\xd7\x55\xde\x76\x57\x8b\x51\x6e\x9e\xfc\xfc\x61\xbc\x5f\x05\xa2\xe0\xa0\x41\x25\xf7\x4c\xe8\x58\xb9\x83\x35\xce\xe5\xbd\x45\x3a\xd2\x83\xa8\xbd\xe3\xed\xf8\xe5\xfc\x08\xe8\xb9\x28\x1f\xc6\xc2\x3e\xb9\x4c\x49\x3e\x38\x4e\x04\xd8\xb5\xa2\x17\x87\xe5\x0c\x9f\xde\x4e\x9f\x04\xcc\x1f\x07\xb8\x01\x6b\x39\x0f\xd9\xa4\xae\x65\x43\xc4\x36\xac\x12\x88\xc0\x84\x99\xae\x14\x2a\x09\x4d\x08\x5d\x77\x82\x4c\x84\x16\x84\x0d\x3b\x73\xaa\xe2\xb9\x56\xbb\x51\xf9\xd3\x46\x97\xad\x39\x66\x8c\x19\xc3\x6f\xa4\x8e\xce\x2c\x70\xb5\x9c\x18\x5c\x7d\x75\x27\x9e\xb4\xa2\x6b\x2b\x13\x02\xf6\x60\x4c\x73\x94\xaa\x9c\x1d\x7e\x01\xe7\xcd\x30\x88\x03\x56\xcc\x09\x46\x52\xbc\x15\x52\x8d\x36\x54\x23\xa1\xda\x51\xcf\xe4\x86\x8b\xe0\x85\x2f\x41\x6a\xd3\xf7\x14\xbc\xa5\x3d\xc2\xc1\xe8\xae\x2b\x8e\x43\x22\xe8\xa7\x83\xb8\x28\x23\x0f\xd3\x45\xf3\x24\x79\x49\xda\xc6\xcf\x37\x9c\x93\xd1\x08\xe9\xe3\x86\x28\x3e\x64\xc0\x5f\xd7\x70\x53\x53\x4d\x20\x2b\x24\xe4\xa5\x93\x6f\xe6\x42\x5d\x24\xc5\xe7\x8e\xa8\xfb\x32\x5d\x4f\x97\x85\x5f\x5d\xb7\x25\x16\x5a\x77\xc6\x84\x29\x28\x4f\x3a\x17\xcb\x76\x9e\xb9\x82\x80\x16\x19\x4f\x2a\x10\xac\xa1\xd0\x35\x86\x4b\x0c\x61\x0c\x1a\xd1\x6f\xe2\xb8\xc3\x5e\x4b\x9f\xf3\x31\x6f\x6f\x55\xd2\x1d\xcc\x5e\xc5\x29\x14\xf3\x96\x37\x04\x3b\xc4\x7c\xc3\xaa\x45\x3a\x2e\x52\xb1\xbf\x92\x6e\xb1\xee\xc1\x35\xfc\x64\x8e\xd4\x09\x85\x56\x32\x1b\x8c\xed\x59\x08\x9e\xef\xc7\x5f\xb4\x68\x55\x6c\x43\x32\xe7\x8e\x6b\x11\x1e\xe2\xc6\x5f\x13\x90\xce\x30\xca\xfa\x32\xd3\x99\xb7\x28\x33\xa2\xcf\x93\xa2\x22\x0c\xe2\x4c\x98\x7a\x26\xd9\x06\x7c\xa6\x84\x0f\xf9\xce\xb6\x69\xb3\x6d\x1a\x6c\x51\x37\xe1\x8d\xce\xa8\xe6\xc2\xe8\x5c\xd8\x9e\x91\x28\x91\xeb\x6c\xc5\x39\x9d\x47\x6b\xe7\x6d\x59\x9c\x1c\x0b\xe7\xd0\x52\xfa\xc4\x21\x6a\x75\x3e\x37\xde\x4f\x91\x6c\xcc\x17\x9a\xc8\x02\xb3\x4d\x33\x99\x3f\x16\xd1\x58\xd0\xc6\xac\x4b\x08\xe0\xbb\xed\x2d\xd5\xd5\x4b\x9f\xc1\xf1\xdf\x6f\x3e\x7e\xbc\xdb\xde\x6e\xfe\xeb\x7b\x72\x21\x4f\x0b\x86\x41\x4d\xf1\xf3\x85\xf2\xd3\x3d\xfa\x1b\xab\x72\xcc\xbb\x24\x00\xd8\x7d\xe3\x0b\x55\xfc\x8c\x62\x39\x4d\x48\xb4\xda\x48\x85\x76\x50\x84\xd6\xa1\x9b\xab\xe6\x4e\xbe\x95\xa8\x1a\x07\xa8\x6b\x65\x5c\x00\xfd\xbd\x15\xf5\x17\xf4\x0e\x56\xbf\xfe\xb6\x9a\x9b\x14\x25\xea\x54\xed\xa6\x14\x4c\x8c\xaa\xb1\xeb\x2b\x3a\xe9\x35\x5c\xdd\x1a\xfd\x4f\xf9\x7b\x81\x22\x47\x93\xf0\x7f\xb8\x06\xee\xd6\xb9\x4d\x75\x9d\x19\x55\x43\x14\x3f\xeb\x11\xbb\x83\xa2\x6c\x17\xbb\x59\xca\x15\x37\x69\x2f\x9e\xf3\x22\x94\x9b\xfa\xa0\xc0\x1a\x3e\x23\x08\xe5\x0c\x58\x0c\x4f\xc7\x39\x69\x42\x71\x7e\x36\xc4\x8d\x73\xcc\x58\x43\xdb\xc5\x34\x73\x48\xc5\x38\xad\x56\xf7\x38\x7f\xb2\xc2\x1b\xd2\xa4\x89\xa3\x17\x57\x83\x95\x3c\xb8\x26\x0c\x5e\x51\xad\x58\x6e\x3e\xe3\xc7\x2f\xa4\x26\x0a\x27\xf3\x3e\x3e\x5a\x2e\xed\x5d\xf3\x78\x66\x1e\x72\x08\x5b\x77\xf2\x29\x21\xe5\xbc\x4c\xfc\x75\x9a\xa6\xe9\x37\xf8\x95\xf5\x36\xed\xe9\x96\xf5\x37\x7e\x3c\x06\x49\x53\xf4\x4c\xcb\xf0\xa9\xca\x0f\x42\xe1\x8a\x1e\xc8\xdf\x5c\x5e\xff\x40\x22\x52\x3f\x42\x40\x10\xca\x57\x1c\x9f\x27\x1a\x2f\x75\x6c\x43\x19\x1a\x73\x44\x65\x8a\x53\x74\xfd\x66\xcf\xd3\x32\xb1\x18\xd9\xa5\x40\x16\x3e\x85\xfb\x1f\x7d\x72\xfa\x61\xf3\xee\x6e\xfb\x78\xf7\xfa\xbb\xf5\x5b\x7e\xe5\x5b\x18\xfa\x4b\xdc\x23\x7e\x73\xf6\xaa\x9c\x52\x2e\xec\x95\xd4\x93\x6e\xf1\xc0\x4b\x0c\xfc\xff\x49\xbf\x13\xf1\x66\xb3\x3d\x22\x2e\x54\x48\x41\xce\xb4\xa6\x95\x35\x28\xa1\x0f\xa3\x38\x20\x1c\xcc\x13\x5a\x7d\xfa\x65\x5f\x9c\x96\xcc\x7c\xdd\x9d\xdf\x6b\xfd\xea\x7f\x03\x00\x00\xff\xff\xb4\xb4\xe2\x86\x5e\x2c\x00\x00")

func licenseBytes() ([]byte, error) {
	return bindataRead(
		_license,
		"LICENSE",
	)
}

func license() (*asset, error) {
	bytes, err := licenseBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "LICENSE", size: 11358, mode: os.FileMode(420), modTime: time.Unix(1586170427, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"LICENSE": license,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"LICENSE": &bintree{license, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
