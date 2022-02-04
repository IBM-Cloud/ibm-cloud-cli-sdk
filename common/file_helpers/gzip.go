package file_helpers

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ExtractTgz extracts src archive to the dest directory. Both src and dest must be a path name.
func ExtractTgz(src string, dest string) error {
	fd, err := os.Open(filepath.Clean(src))
	if err != nil {
		return err
	}

	defer func() {
		if err := fd.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()

	gReader, err := gzip.NewReader(fd)
	if err != nil {
		return err
	}
	defer gReader.Close()

	tarReader := tar.NewReader(gReader)

	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if hdr.Name == "." {
			continue
		}

		err = extractFileInArchive(tarReader, hdr, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractFileInArchive(r io.Reader, hdr *tar.Header, dest string) error {
	fi := hdr.FileInfo()
	path := filepath.Join(filepath.Clean(dest), filepath.Clean(hdr.Name))

	if fi.IsDir() {
		return os.MkdirAll(path, fi.Mode())
	} else {
		err := os.MkdirAll(filepath.Dir(path), 0700)
		if err != nil {
			return err
		}

		//#nosec G304 -- This is a false positive
		f, err := os.Create(path)
		if err != nil {
			return err
		}

		defer func() {
			if err := f.Close(); err != nil {
				fmt.Printf("Error closing file: %s\n", err)
			}
		}()

		_, err = io.Copy(f, r)
		return err
	}
}
