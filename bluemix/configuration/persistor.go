package configuration

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
	runtimeGOOS   string
}

func NewDiskPersistor(path string) DiskPersistor {
	return DiskPersistor{
		filePath:      path,
		fileLock:      flock.New(path),
		parentContext: context.Background(),
		runtimeGOOS:   runtime.GOOS,
	}
}

func (dp DiskPersistor) Exists() bool {
	return file_helpers.FileExists(dp.filePath)
}

func (dp *DiskPersistor) windowsLockedRead(data DataInterface) error {
	// TO DO: exclusive file-locking for the reading NOT yet implemented
	return dp.read(data)
}

func isBlockingLockError(err error) bool {
	return err != nil && !strings.Contains(err.Error(), "no such file or directory")
}

func (dp *DiskPersistor) lockedRead(data DataInterface) error {
	lockCtx, cancelLockCtx := context.WithTimeout(dp.parentContext, 30*time.Second) /* allotting a 30-second timeout means there can be a maximum of 298 failed retrials (each up to 500 ms, as
	specified after the deferred call to cancelLockCtx). 30 appears to be a conventional value for a parent context passed to TryLockContext, as per docs */
	defer cancelLockCtx()
	_, lockErr := dp.fileLock.TryLockContext(lockCtx, 100*time.Millisecond) /* provide a file lock just while dp.read is called, because it calls an unmarshaling function
	The boolean (first return value) can be wild-carded because lockErr must be non-nil when the lock-acquiring fails (whereby the boolean will be false) */
	defer dp.fileLock.Unlock()
	if isBlockingLockError(lockErr) {
		return lockErr
	}
	readErr := dp.read(data)
	if readErr != nil {
		return readErr
	}
	return nil
}

func (dp DiskPersistor) readWithFileLock(data DataInterface) error {
	switch dp.runtimeGOOS {
	case "windows":
		return dp.windowsLockedRead(data)
	default:
		return dp.lockedRead(data)
	}
}

func (dp DiskPersistor) writeWithFileLock(data DataInterface) error {
	switch dp.runtimeGOOS {
	case "windows":
		return dp.windowsLockedWrite(data)
	default:
		return dp.lockedWrite(data)
	}
}

func (dp DiskPersistor) Load(data DataInterface) error {
	err := dp.readWithFileLock(data)
	if os.IsPermission(err) {
		return err
	}

	if err != nil { /* would happen if there was nothing to read (EOF) */
		err = dp.writeWithFileLock(data)
	}
	return err
}

func (dp *DiskPersistor) windowsLockedWrite(data DataInterface) error {
	// TO DO: exclusive file-locking for the writing NOT yet implemented
	return dp.write(data)
}

func (dp DiskPersistor) lockedWrite(data DataInterface) error {
	lockCtx, cancelLockCtx := context.WithTimeout(dp.parentContext, 30*time.Second) /* allotting a 30-second timeout means there can be a maximum of 298 failed retrials (each up to 500 ms, as
	specified after the deferred call to cancelLockCtx). 30 appears to be a conventional value for a parent context passed to TryLockContext, as per docs */
	defer cancelLockCtx()
	_, lockErr := dp.fileLock.TryLockContext(lockCtx, 100*time.Millisecond) /* provide a file lock just while dp.read is called, because it calls an unmarshaling function
	The boolean (first return value) can be wild-carded because lockErr must be non-nil when the lock-acquiring fails (whereby the boolean will be false) */
	defer dp.fileLock.Unlock()
	if isBlockingLockError(lockErr) {
		return lockErr
	}
	writeErr := dp.write(data)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func (dp DiskPersistor) Save(data DataInterface) error {
	return dp.writeWithFileLock(data)
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
