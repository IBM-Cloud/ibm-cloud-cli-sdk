// Package file_helpers provides file operation helpers.
package file_helpers

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileExists checks if the file exist or not
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// RemoveFile removes the file.
// If the file does not exist, it do nothing and return nil.
func RemoveFile(path string) error {
	if FileExists(path) {
		return os.Remove(path)
	}
	return nil
}

// CopyFile copies file contents of src to dest. Both of stc and dest must be a path name.
func CopyFile(src string, dest string) (err error) {
	srcFile, err := os.Open(filepath.Clean(src))
	if err != nil {
		return
	}
	/* #nosec G307 */
	defer func() {
		if err := srcFile.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()

	srcStat, err := srcFile.Stat()
	if err != nil {
		return
	}

	if !srcStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file.", src)
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return
	}
	/* #nosec G307 */
	defer func() {
		if err := destFile.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()

	err = os.Chmod(dest, srcStat.Mode())
	if err != nil {
		return
	}

	_, err = io.Copy(destFile, srcFile)
	return
}

// CopyDir copies src directory recursively to dest.
// Both of src and dest must be a path name.
// It returns error if dest already exist.
func CopyDir(src string, dest string) (err error) {
	srcStat, err := os.Stat(src)
	if err != nil {
		return
	}

	if !srcStat.Mode().IsDir() {
		return fmt.Errorf("%s is not a directory.", src)
	}

	err = os.MkdirAll(dest, srcStat.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.Mode().IsDir() {
			err = CopyDir(srcPath, destPath)
		} else {
			err = CopyFile(srcPath, destPath)
		}
		if err != nil {
			return
		}
	}

	return
}
