package rest_test
import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/models"
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
	testhelpers "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/configuration"

	"github.com/stretchr/testify/assert"
)

func TestRequestQueryParam(t *testing.T) {
	assert := assert.New(t)

	req, err := GetRequest("http://www.example.com?foo=fooVal1").
		Query("foo", "fooVal2").
		Query("bar", "bar Val").
		Build()

	assert.NoError(err)
	assert.Contains(req.URL.String(), "foo=fooVal1")
	assert.Contains(req.URL.String(), "foo=fooVal2")
	assert.Contains(req.URL.String(), "bar=bar+Val")
}

func TestRequestHeader(t *testing.T) {
	assert := assert.New(t)

	req, err := GetRequest("http://www.example.com").
		Set("Accept", "application/json").
		Add("Accept-Encoding", "compress").
		Add("Accept-Encoding", "gzip").
		Build()

	assert.NoError(err)
	assert.Equal("application/json", req.Header.Get("Accept"))
	assert.Equal([]string{"compress", "gzip"}, req.Header["Accept-Encoding"])
}

func TestRequestFormText(t *testing.T) {
	assert := assert.New(t)

	req, err := PostRequest("http://www.example.com").
		Field("foo", "bar").
		Build()

	assert.NoError(err)
	err = req.ParseForm()
	assert.NoError(err)
	assert.Equal(FormUrlEncodedContentType, req.Header.Get(ContentType))
	assert.Equal("bar", req.FormValue("foo"))
}

func TestRequestFormMultipart(t *testing.T) {
	assert := assert.New(t)

	var prepareFileWithContent = func(text string) (*os.File, error) {
		f, err := os.CreateTemp("", "BluemixCliRestTest")
		if err != nil {
			return nil, err
		}
		_, err = f.WriteString(text)
		if err != nil {
			return nil, err
		}

		_, err = f.Seek(0, 0)
		if err != nil {
			return nil, err
		}

		return f, err
	}

	f, err := prepareFileWithContent("12345")
	assert.NoError(err)

	req, err := PostRequest("http://www.example.com").
		Field("foo", "bar").
		File("file1", File{Name: f.Name(), Content: f}).
		File("file2", File{Name: "file2.txt", Content: strings.NewReader("abcde"), Type: "text/plain"}).
		Build()

	assert.NoError(err)
	assert.Contains(req.Header.Get(ContentType), "multipart/form-data")

	err = req.ParseMultipartForm(int64(5000))
	assert.NoError(err)

	assert.Equal(1, len(req.MultipartForm.Value))
	assert.Equal("bar", req.MultipartForm.Value["foo"][0])

	assert.Equal(2, len(req.MultipartForm.File))

	assert.Equal(1, len(req.MultipartForm.File["file1"]))
	// As of Golang 1.17, a temp file always has a directory, so compare the string suffix that is the filename
	assert.True(strings.HasSuffix(f.Name(), req.MultipartForm.File["file1"][0].Filename))
	assert.Equal("application/octet-stream", req.MultipartForm.File["file1"][0].Header.Get("Content-Type"))

	assert.Equal(1, len(req.MultipartForm.File["file2"]))
	assert.Equal("file2.txt", req.MultipartForm.File["file2"][0].Filename)
	assert.Equal("text/plain", req.MultipartForm.File["file2"][0].Header.Get("Content-Type"))

	b1 := new(bytes.Buffer)
	f1, _ := req.MultipartForm.File["file1"][0].Open()
	io.Copy(b1, f1)
	assert.Equal("12345", string(b1.Bytes()))

	b2 := new(bytes.Buffer)
	f2, _ := req.MultipartForm.File["file2"][0].Open()
	io.Copy(b2, f2)
	assert.Equal("abcde", string(b2.Bytes()))
}

func TestRequestJSON(t *testing.T) {
	assert := assert.New(t)

	var foo = struct {
		Name string
	}{
		Name: "bar",
	}

	req, err := PostRequest("http://www.example.com").Body(&foo).Build()
	assert.NoError(err)
	body, err := io.ReadAll(req.Body)
	assert.NoError(err)
	assert.Equal("{\"Name\":\"bar\"}", string(body))
}

func TestCachedPaginationNextURL(t *testing.T) {
	testCases := []struct {
		name            string
		paginationURLs  []models.PaginationURL
		offset          int
		expectedNextURL string
	}{
		{
			name:   "return cached next URL",
			offset: 200,
			paginationURLs: []models.PaginationURL{
				{
					NextURL:   "/v2/example.com/stuff?limit=100",
					LastIndex: 100,
				},
			},
			expectedNextURL: "/v2/example.com/stuff?limit=100",
		},
		{
			name:   "return empty string if cache URL cannot be determined",
			offset: 40,
			paginationURLs: []models.PaginationURL{
				{
					NextURL:   "/v2/example.com/stuff?limit=100",
					LastIndex: 60,
				},
			},
			expectedNextURL: "",
		},
		{
			name:            "return empty string if no cache available",
			offset:          40,
			paginationURLs:  []models.PaginationURL{},
			expectedNextURL: "",
		},
	}

	assert := assert.New(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := CachedPaginationNextURL(tc.paginationURLs, tc.offset)
			assert.Equal(tc.expectedNextURL, url)
		})
	}
}

func TestAddPaginationURL(t *testing.T) {
	config := testhelpers.NewFakeCoreConfig()
	assert := assert.New(t)
	unsortedUrls := []models.PaginationURL{
		{
			NextURL:   "/v2/example.com/stuff?limit=200",
			LastIndex: 200,
		},
		{
			NextURL:   "/v2/example.com/stuff?limit=100",
			LastIndex: 100,
		},
	}

	var err error
	for _, p := range unsortedUrls {
		err = config.AddPaginationURL(p.LastIndex, p.NextURL)
		assert.Nil(err)
	}

	// expect url to be sorted in ascending order by LastIndex
	sortedUrls, err := config.PaginationURLs()
	assert.Nil(err)

	assert.Equal(2, len(sortedUrls))
	assert.Equal(sortedUrls[0].LastIndex, unsortedUrls[1].LastIndex)
	assert.Equal(sortedUrls[0].NextURL, unsortedUrls[1].NextURL)
	assert.Equal(sortedUrls[1].LastIndex, unsortedUrls[0].LastIndex)
	assert.Equal(sortedUrls[1].NextURL, unsortedUrls[0].NextURL)
}
