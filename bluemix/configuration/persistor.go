package configuration

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/file_helpers"
	"github.com/gofrs/flock"
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
	filePath      string
	fileLock      *flock.Flock
	parentContext context.Context
}

func NewDiskPersistor(path string) DiskPersistor {
	return DiskPersistor{
		filePath:      path,
		fileLock:      flock.New(path),
		parentContext: context.Background(),
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

	if err != nil { // strange: requiring a reading error (to allow write attempt to continue), as long as it is not a permission error
		err = dp.lockedWrite(data)
	}
	return err
}

func (dp DiskPersistor) lockedWrite(data DataInterface) error {
	lockCtx, cancelLockCtx := context.WithTimeout(dp.parentContext, 3*time.Second) /* allotting a 3-second timeout means there can be a maximum of 5 retrials (each up to 500 ms, as
	specified after the deferred call to cancelLockCtx) */
	defer cancelLockCtx()
	_, lockErr := dp.fileLock.TryLockContext(lockCtx, 500*time.Millisecond) /* provide a file lock, in addition to the RW mutex (in calling functions), just while dp.write is called
	The boolean (first return value) can be wild-carded because lockErr must be non-nil when the lock-acquiring fails */
	if lockErr != nil {
		return lockErr
	}

	writeErr := dp.write(data)
	if writeErr != nil {
		return writeErr
	}
	return dp.fileLock.Unlock()
}

func (dp DiskPersistor) Save(data DataInterface) error {
	return dp.lockedWrite(data)
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
