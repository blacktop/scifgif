// Code generated by go-bindata.
// sources:
// config/config.toml
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

var _configConfigToml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x53\x4b\x6f\xdb\x30\x0c\xbe\xeb\x57\x10\xce\x65\x03\x56\xc7\x8f\x66\xe9\x0c\xe4\x50\x14\x3d\x74\x58\x37\xa0\x3d\x06\xc5\x40\xdb\x8c\x2d\x44\x2f\x48\x72\xfa\xf8\xf5\x03\x95\xf4\x11\xac\x87\x0d\xa8\x72\x50\xc4\x8f\xe4\xf7\xf1\x63\x32\x83\x0b\xeb\x1e\xbd\x1c\xc6\x08\x9f\xba\xcf\x50\x15\x65\x0d\x27\x7c\x2d\xa1\x55\xd8\x6d\xa3\x75\xf0\xdd\x86\x71\x42\xb8\x46\x69\xe8\x0b\x9c\x2b\x05\x37\x5c\x10\xe0\x86\x02\xf9\x1d\xf5\xb9\x98\xc1\x2d\x11\xfc\xb8\xba\xb8\xfc\x79\x7b\x09\x1b\xeb\x41\xc9\x8e\x4c\x20\x90\x66\x63\xbd\xc6\x28\xad\xc9\x85\x98\x7d\xcc\x11\x33\xb8\x3e\x67\x36\xb8\xb0\x66\x23\x87\xc9\x27\x02\xf8\xff\x3e\x1f\xa4\x47\x44\x19\x15\xc1\x0a\xb2\x6b\xe4\xc9\xe1\x66\x32\x51\x6a\x3a\xd6\x97\x89\x1d\xf9\xc0\x42\x57\x90\xed\x8a\xbc\xce\xcb\x32\x13\x62\x8d\x53\x1c\xad\xbf\x13\x00\x06\x75\xea\xf2\xec\x7d\x26\x00\xac\x1f\xd0\xc8\xa7\xfd\x84\x2f\x0c\x57\xbf\xb8\xf2\x9e\x5a\x2e\x9b\xbc\x62\xa4\xc8\xd3\xa7\x39\x2b\xb8\x0e\x7b\x2d\xcd\xef\x03\x54\x56\xcb\x04\x96\x4d\x5d\xd7\x35\x97\x92\x46\xa9\xb8\x78\xb4\x21\x72\x4a\xd0\xd1\xe5\xf4\x80\xda\x29\xca\x3b\xab\xb9\x87\xb3\x9e\xb1\x6a\xc1\x24\x81\x3c\xe7\xf1\xcd\x3a\x13\x8e\x21\x70\x8c\xef\x7b\xeb\x7b\x6e\xdc\x63\xc4\x16\x03\xbd\x9d\x47\x27\xcd\x27\xa4\x30\x44\xd9\x71\xa5\xd4\x38\xbc\x81\xe6\x07\x28\x10\xfa\x6e\x6c\x16\xf9\x82\x93\xd2\xcf\x2b\x91\x2a\xdb\xa1\x62\xa5\xcf\xaa\x98\x76\xfd\xad\x2a\x0a\xa6\x61\xab\xed\x94\x94\x16\x02\x80\x0c\xb6\x8a\x7a\x58\x41\xf4\x13\x09\xb1\x9e\xe4\x3b\x62\xb6\xb2\x45\x83\xef\x69\xd9\x23\xff\x2a\xe2\xf4\xb4\xbe\x7b\x8f\x94\xcc\x4e\x7a\x6b\x34\x99\xc8\xb8\x9f\xd2\xf6\x7a\xda\x91\xb2\x8e\xa3\xc9\x2c\xdb\x6d\x29\xad\x5e\x63\x37\x4a\x43\x27\xc7\x2a\xb3\xd4\xb9\x77\x56\x9a\xb4\xa4\xd8\xb9\x66\x3e\x7f\x11\xd2\x54\xf5\xf2\x6b\x76\xe4\x40\x99\x2c\x68\xa5\xe9\xc3\x6b\x9b\x66\xae\x51\xdd\xa3\xa7\xc6\x5b\x4e\x57\xd2\x6c\xc3\xdf\x8b\x69\x8e\xb6\xc0\x89\x9d\x9b\x60\x05\x8b\xe2\x70\x58\x27\x69\xeb\x1f\x39\x58\x9d\x56\x67\x67\x1c\x14\x6b\x65\x87\x61\x3f\xc6\x46\x2a\x3a\x1e\x21\x57\x76\xc8\xd2\x80\x0f\x41\x3e\x31\x50\x16\xfb\xe7\xde\xf5\xfa\xf0\x6a\xb1\xdb\x4e\x8e\x55\x2d\x59\x21\x8f\x98\xfe\x42\x2b\xd8\xa0\x0a\xec\xa8\xf3\xf6\xe1\xf1\xd5\xeb\x17\x04\x60\x8c\xd1\x31\x63\x76\xf8\x1e\xf6\x8f\x3f\x01\x00\x00\xff\xff\x6c\x61\x71\xdd\xdf\x04\x00\x00")

func configConfigTomlBytes() ([]byte, error) {
	return bindataRead(
		_configConfigToml,
		"config/config.toml",
	)
}

func configConfigToml() (*asset, error) {
	bytes, err := configConfigTomlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "config/config.toml", size: 1247, mode: os.FileMode(420), modTime: time.Unix(1500767390, 0)}
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
	"config/config.toml": configConfigToml,
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
	"config": &bintree{nil, map[string]*bintree{
		"config.toml": &bintree{configConfigToml, map[string]*bintree{}},
	}},
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

