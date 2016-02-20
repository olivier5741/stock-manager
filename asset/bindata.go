// Code generated by go-bindata.
// sources:
// bindata.go
// en-us.all.yaml
// fr-be.all.yaml
// translation.go
// DO NOT EDIT!

package asset

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

var _bindataGo = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func bindataGoBytes() ([]byte, error) {
	return bindataRead(
		_bindataGo,
		"bindata.go",
	)
}

func bindataGo() (*asset, error) {
	bytes, err := bindataGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "bindata.go", size: 0, mode: os.FileMode(436), modTime: time.Unix(1455825040, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _enUsAllYaml = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x7c\x94\x4d\x4e\xc3\x30\x10\x85\xf7\x3d\xc5\x5c\x80\x1e\xa0\xdb\x4a\xec\x10\x48\x08\xb6\x96\x49\x9c\xd6\x22\x99\x09\xce\x24\x50\x55\xbd\x3b\x63\x3b\xe5\xc7\x3f\x5d\x81\xde\x7c\xef\x79\xfc\xe4\xf4\x0e\x6c\xbb\x83\x66\x5a\xd4\xd1\xe8\xd6\x38\xd5\xe9\x86\xc9\xa9\x89\xa9\x79\x57\x33\x5a\x56\xe4\xbc\xee\xff\xdd\x00\xb0\xd3\x38\xf5\x9a\x2d\xe1\x0e\xee\x03\xbb\xc9\x32\x2c\x9b\x41\x8d\x8e\xda\xb9\xc9\x3c\x4f\x22\x57\x1c\x8b\xee\x67\x93\xf2\xaf\xba\x87\xf3\x79\xbb\xa7\x19\xf9\x72\xa9\x38\x4b\xcb\xbd\x88\x76\xd3\x59\xbf\xd7\xa3\x9f\x84\x80\xdc\x15\x8b\x19\x2c\xa6\xa6\x07\x8b\xdb\x1a\x5e\x3a\xe4\xd9\x4f\xfe\x1e\x62\x9c\x93\xe2\x5b\xeb\x94\x13\x6f\x8a\xef\x35\x22\x31\xf8\x11\x74\x8e\x06\x10\xd0\xf8\xfa\x4f\xff\xec\x9d\xed\x8d\x6a\x84\xe2\xac\xc9\x35\x21\x0e\xc1\x83\x05\xa7\x6c\x3e\xe3\xa0\xdd\x74\xd4\x7d\x25\xe0\x67\x7e\x23\x63\xb1\xe6\x53\x31\x55\x12\xae\x7e\x3e\x1a\xf0\x24\x30\xf9\xce\xf2\x28\x1a\x4d\xd6\xf3\x9a\xe1\x47\xe5\x05\x50\x0f\x46\x21\x29\x79\x9c\x62\x50\xbe\xac\x4a\x46\x67\xb1\x85\x88\xc9\x9f\x46\x1e\x84\xc5\x83\x5f\xe6\x1a\x53\xce\x1e\x65\xfd\x5a\xbb\x61\x56\xf6\x7f\xcc\xc6\x9d\xe2\x83\xa8\x98\x03\x01\x81\x80\x56\xb3\x7e\xd3\xd3\x35\x22\xd4\x11\x4e\x6f\x08\x3b\x7b\x48\x13\x56\x35\x85\x2d\x2e\x06\xc3\x23\x49\xf8\xdf\x41\x6a\xe9\x29\x0b\x0f\xfb\x4f\x19\x19\x3e\xa0\x94\x8d\x62\x8a\x56\x7e\x0c\x56\x39\x8f\x2e\xb6\x14\xc5\x22\xaa\xf2\xef\x51\x94\x32\x4a\x73\xb6\x87\x97\x52\x98\xa9\x7c\x41\x79\x1e\xd9\x1d\x47\x67\x3a\xfb\xa5\x5a\xa7\xbb\x2c\x3b\x8a\xdf\x01\x00\x00\xff\xff\x92\x05\xdf\x75\x68\x05\x00\x00")

func enUsAllYamlBytes() ([]byte, error) {
	return bindataRead(
		_enUsAllYaml,
		"en-us.all.yaml",
	)
}

func enUsAllYaml() (*asset, error) {
	bytes, err := enUsAllYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "en-us.all.yaml", size: 1384, mode: os.FileMode(436), modTime: time.Unix(1454155656, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _frBeAllYaml = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x7c\x94\x41\x6e\xdc\x3a\x0c\x86\xf7\x39\x05\x2f\xf0\x72\x80\x6c\x67\xfb\x0a\x14\x28\xda\xad\xc0\x5a\x74\x47\xa8\x2d\x3a\x94\xe4\x36\x08\x72\x97\x6e\xdd\x6b\xf8\x62\xa5\xc6\xd3\x26\x23\xda\xd9\x19\x32\xff\x8f\xfc\x49\x89\xff\x41\xf0\x0f\xd0\xa5\xd9\x9d\x09\x3d\x89\xeb\xb1\xcb\x2c\x2e\x65\xee\xbe\xbb\x12\x43\x76\x2c\xf5\xbc\x7e\xde\x01\x64\xc1\x98\x06\xcc\x81\xe3\x03\x9c\x98\xfa\xfe\xfe\xce\x30\x42\xa6\xd1\x4d\xc2\xbe\x74\x46\xf3\xb1\x1e\x2b\x6a\x5f\x34\xe3\x50\xa8\x95\x7c\xc1\x81\x8a\xc0\xf3\xf3\xfd\x89\x4b\xcc\x2f\x2f\x07\xe2\xbd\x12\x3f\xeb\x19\xbd\x2b\x3d\xb6\x57\xb5\xeb\xa2\x2e\xc7\x11\xa3\x27\x2b\xdd\x9a\x34\x86\xd8\x2a\x3f\x84\xb8\xd3\x96\xd7\x9e\x1e\x64\xfa\x54\x03\xae\x3a\x12\xd1\x31\xf8\x20\x4e\x54\xdd\x0a\xfe\x27\xf0\x9c\x52\x20\x81\x48\x30\x51\xc9\xb0\xfe\xce\x42\x30\x94\x1b\x7d\x1f\x06\x72\x9d\x12\xb2\xe9\xaa\x22\xfa\xd0\x9d\x2d\xa2\x93\x75\x59\x97\x1d\x8c\x3a\x29\x71\x44\x49\x67\x1c\xde\xa1\x69\x5c\x4b\xe4\x38\x93\xe4\xb0\xcf\x9c\x03\xfd\x70\x99\x0d\x11\x61\x2e\x74\x40\x22\xa0\x58\x13\x59\x22\x4f\x64\xc6\x71\x68\x95\x4b\xa5\x19\x48\xc4\x91\x5c\x64\xa7\x4f\x41\xf5\xae\x17\x1e\xed\xc4\x08\xb6\xdf\x0d\x32\x8b\x42\xd7\x85\x60\xfd\x05\x13\x6a\xa9\x02\x5e\xeb\x56\x23\x3c\xd6\xaf\x6b\x21\xfb\x39\x55\x90\xcc\xa0\x4e\xad\xd6\x64\xc4\x7a\x7b\x6e\x88\x8f\x85\xe4\x69\xbb\x70\x3b\x7d\x15\x7a\x2c\xaa\x25\x50\xf7\x09\x06\x84\xaf\x98\xa8\x26\xf0\x1c\x63\xad\xdd\x53\x82\x8b\x38\x01\xc2\xba\x74\x67\x2e\xff\x12\x5c\xba\x7c\xa9\x56\x87\xd1\x87\x6f\x2d\xff\x7a\xda\x06\x07\x1d\x5c\xd4\xcd\xf2\xd4\xc6\x6f\x3f\x30\x08\x19\xcd\xc0\x86\xae\xf6\x74\x17\x24\x13\x7a\x79\xc5\xb6\x94\x9b\x97\xfb\x1a\x7d\xb0\x9b\xa6\x6d\x37\x59\xfa\x6e\x23\xd3\x9b\xd7\xda\x84\x3a\xbb\x12\xd4\xa3\x90\xad\x64\x0b\xe7\x62\x6a\x49\x5c\xaf\xb9\x89\xcf\xbc\x6f\x55\x6f\xdb\x5f\xb7\xf2\x56\x34\x09\xf5\xe1\xa7\xf3\x82\xbd\x49\xa1\x4f\x08\x73\xd6\xc2\xe8\x4f\x00\x00\x00\xff\xff\x1d\x6d\x10\x0a\x06\x06\x00\x00")

func frBeAllYamlBytes() ([]byte, error) {
	return bindataRead(
		_frBeAllYaml,
		"fr-be.all.yaml",
	)
}

func frBeAllYaml() (*asset, error) {
	bytes, err := frBeAllYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "fr-be.all.yaml", size: 1542, mode: os.FileMode(436), modTime: time.Unix(1455664008, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _translationGo = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x90\xb1\x6a\xc3\x30\x10\x86\x67\xeb\x29\x0e\x4f\x36\xd8\x16\xc9\x54\x02\x19\x5a\x4a\xd6\x74\xf0\x5e\xce\xf6\xd9\x11\x91\x4f\x41\x92\x1b\x44\xe9\xbb\x57\x4a\xd3\xc1\xed\x92\x45\x02\xe9\xbb\xfb\x3f\xfe\x0b\xf6\x67\x9c\x08\xd0\x39\xf2\x42\xa8\xf9\x62\xac\x2f\x44\x96\x4f\xca\x9f\x96\xae\xe9\xcd\x2c\x59\xf5\x67\xc7\x61\x20\x2b\x27\x53\xab\xcd\x13\xcb\x74\xe4\xa2\x14\xe2\x03\x6d\xa4\x5b\x0b\xe9\xa5\x69\x2d\xb2\xd3\xe8\xe9\xb0\x70\x9f\xbe\xc7\x78\x83\x62\xe5\x8b\x12\x3e\x45\x26\x25\xb4\x47\x78\x3d\xc2\x0e\x46\xc5\x03\x20\x1b\x7f\x22\x0b\x57\x0c\x30\x63\xe8\x08\xae\x31\xf6\x47\xc6\x55\x30\x19\xe8\x22\x56\x0f\xe8\x51\x64\x66\xf1\x9b\x0a\xde\x61\xb7\x87\xe7\x04\x14\x39\x71\xbd\xb8\x06\xb5\x6e\x02\xce\x3a\x2f\x6f\xcc\x76\xcd\x8c\xb6\xee\x68\xc5\xdc\x4c\xdf\xd0\x3a\xfa\xd5\x55\x86\x0f\x4a\xd3\x4b\xf0\xe4\xfe\x6d\xad\x52\xee\x03\x63\x7f\x82\xd2\xd8\xb6\x4c\xd5\x24\x9f\xfd\xbd\x9f\xd4\xc7\x1d\x8d\x2a\x5f\xdf\x01\x00\x00\xff\xff\xce\xec\xc1\xd0\x7e\x01\x00\x00")

func translationGoBytes() ([]byte, error) {
	return bindataRead(
		_translationGo,
		"translation.go",
	)
}

func translationGo() (*asset, error) {
	bytes, err := translationGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "translation.go", size: 382, mode: os.FileMode(436), modTime: time.Unix(1455823482, 0)}
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
	"bindata.go": bindataGo,
	"en-us.all.yaml": enUsAllYaml,
	"fr-be.all.yaml": frBeAllYaml,
	"translation.go": translationGo,
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
	"bindata.go": &bintree{bindataGo, map[string]*bintree{}},
	"en-us.all.yaml": &bintree{enUsAllYaml, map[string]*bintree{}},
	"fr-be.all.yaml": &bintree{frBeAllYaml, map[string]*bintree{}},
	"translation.go": &bintree{translationGo, map[string]*bintree{}},
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

