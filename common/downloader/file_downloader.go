// Package downloader provides a simple file downlaoder
package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ProxyReader is an interface to proxy read bytes
type ProxyReader interface {
	Proxy(size int64, reader io.Reader) io.Reader
	Finish()
}

// FileDownloader is a file downloader
type FileDownloader struct {
	SaveDir       string       // path of the directory to save the downloaded file
	DefaultHeader http.Header  // Default header to applied to the download request
	Client        *http.Client // HTTP client to use, default is http DefaultClient
	ProxyReader   ProxyReader
}

// New creates a file downloader
func New(saveDir string) *FileDownloader {
	return &FileDownloader{
		SaveDir:       saveDir,
		Client:        http.DefaultClient,
		DefaultHeader: make(http.Header),
	}
}

// Download downloads a file from a URL and returns the path and size of the
// downloaded file
func (d *FileDownloader) Download(url string) (dest string, size int64, err error) {
	return d.DownloadTo(url, "")
}

// DownloadTo downloads a file from a URL to the file with the specified name in
// the download directory
func (d *FileDownloader) DownloadTo(url string, outputName string) (dest string, size int64, err error) {
	req, err := d.createRequest(url)
	if err != nil {
		return "", 0, fmt.Errorf("download request error: %v", err)
	}

	client := d.Client
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", 0, fmt.Errorf("Unexpected response code %d", resp.StatusCode)
	}

	if outputName == "" {
		outputName = d.determinOutputName(resp)
	}
	dest = filepath.Join(d.SaveDir, outputName)

	f, err := os.OpenFile(filepath.Clean(dest), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return dest, 0, err
	}
	/* #nosec G307 */
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()

	var r io.Reader = resp.Body
	if d.ProxyReader != nil {
		defer d.ProxyReader.Finish()
		r = d.ProxyReader.Proxy(resp.ContentLength, r)
	}

	size, err = io.Copy(f, r)
	if err != nil {
		return dest, size, err
	}

	return dest, size, nil
}

// RemoveDir removes the download directory
func (d *FileDownloader) RemoveDir() error {
	return os.RemoveAll(d.SaveDir)
}

func (d *FileDownloader) createRequest(url string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if d.DefaultHeader != nil {
		r.Header = d.DefaultHeader
	}

	if r.Header.Get("User-Agent") == "" {
		r.Header.Set("User-Agent", "bluemix-cli")
	}

	return r, nil
}

func (d *FileDownloader) determinOutputName(resp *http.Response) string {
	n := getFileNameFromHeader(resp.Header.Get("Content-Disposition"))

	if n == "" {
		n = getFileNameFromUrl(resp.Request.URL)
	}

	if n == "" {
		n = "index.html"
	}

	return n
}

func getFileNameFromUrl(url *url.URL) string {
	path := path.Clean(url.Path)
	if path == "." {
		return ""
	}

	fields := strings.Split(path, "/")
	if len(fields) == 0 {
		return ""
	}

	return fields[len(fields)-1]
}

func getFileNameFromHeader(header string) string {
	if header == "" {
		return ""
	}

	for _, field := range strings.Split(header, ";") {
		field = strings.TrimSpace(field)

		if strings.HasPrefix(field, "filename=") {
			name := strings.TrimLeft(field, "filename=")
			return strings.Trim(name, `"`)
		}
	}

	return ""
}
