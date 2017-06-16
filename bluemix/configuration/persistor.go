package configuration

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/common/file_helpers"
)

const (
	filePermissions = 0600
	dirPermissions  = 0700
)

type DataInterface interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

type Persistor interface {
	Exists() bool
	Load(DataInterface) error
	Save(DataInterface) error
}

type DiskPersistor struct {
	filePath string
}

func NewDiskPersistor(path string) DiskPersistor {
	return DiskPersistor{
		filePath: path,
	}
}

func (dp DiskPersistor) Exists() bool {
	return file_helpers.FileExists(dp.filePath)
}

func (dp DiskPersistor) Load(data DataInterface) error {
	err := dp.read(data)
	if os.IsPermission(err) {
		return err
	}

	if err != nil {
		err = dp.write(data)
	}
	return err
}

func (dp DiskPersistor) Save(data DataInterface) error {
	return dp.write(data)
}

func (dp DiskPersistor) read(data DataInterface) error {
	err := os.MkdirAll(filepath.Dir(dp.filePath), dirPermissions)
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadFile(dp.filePath)
	if err != nil {
		return err
	}

	err = data.Unmarshal(bytes)
	return err
}

func (dp DiskPersistor) write(data DataInterface) error {
	bytes, err := data.Marshal()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dp.filePath, bytes, filePermissions)
	return err
}
