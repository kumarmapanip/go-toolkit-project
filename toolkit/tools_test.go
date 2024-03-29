package toolkit

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateRandomPassword(t *testing.T) {
	tools := ToolKit{}

	password := tools.RandomPasswordGenerator()

	n := len(password)
	assert.EqualValues(t, n > minL && n < maxL, true)
}

// another way of testing from stdout
func Test_CheckTimeElapsed(t *testing.T) {
	stdOut := os.Stdout
	read, write, _ := os.Pipe()
	os.Stdout = write

	tools := ToolKit{}

	_ = tools.RandomPasswordGenerator()

	write.Close()

	data, _ := io.ReadAll(read)
	result := string(data)

	os.Stdout = stdOut

	if !strings.Contains(result, "password") {
		t.Errorf("Exception failed running test cases #02")
	}
}

// Table testing

var uploadTests = []struct {
	name          string
	allowedTypes  []string
	renameFile    bool
	errorExpected bool
}{
	{
		name:          "allowed no rename",
		allowedTypes:  []string{"image/jpg", "image/png"},
		renameFile:    false,
		errorExpected: false,
	},
	{
		name:          "allowed rename",
		allowedTypes:  []string{"image/jpg", "image/png"},
		renameFile:    true,
		errorExpected: false,
	},
	{
		name:          "not allowed",
		allowedTypes:  []string{"image/jpgp"},
		renameFile:    true,
		errorExpected: true,
	},
}

func Test_UploadMultipleFiles(t *testing.T) {
	for _, testcase := range uploadTests {
		// set up pipe to avoid buffering
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer func() {
				writer.Close()
				wg.Done()
			}()

			// create the form data file for multipart file upload
			part, err := writer.CreateFormFile("file", "./testdata/test-1.png")
			if err != nil {
				t.Error(err)
			}

			f, err := os.Open("./testdata/test-1.png")
			if err != nil {
				t.Error(err)
			}

			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				t.Error(err)
			}

			err = png.Encode(part, img)
			if err != nil {
				t.Error(err)
			}
		}()

		// read the  pipe whihc recieves data
		request := httptest.NewRequest("POST", "/", pr)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var testTools ToolKit
		testTools.AllowedFileTypes = testcase.allowedTypes

		uploadedFiles, err := testTools.UploadFiles(request, "./testdata/uploads/", testcase.renameFile)
		if err != nil && !testcase.errorExpected {
			t.Error(err)
		}

		// check file is presnt actually or not
		if !testcase.errorExpected {
			_, err = os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName))
			if os.IsNotExist(err) {
				t.Errorf("%s expected file to exist err: %v", testcase.name, err)
			}

			log.Println("Testcase 1: file is uploaded")
			// clean up
			_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName))
		}

		if testcase.errorExpected && err == nil {
			t.Errorf("Expected error %s", testcase.name)
		}

		wg.Wait()
	}
}

func Test_UploadOneFile(t *testing.T) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer func() {
			writer.Close()
		}()

		// create the form data file for multipart file upload
		part, err := writer.CreateFormFile("file", "./testdata/test-1.png")
		if err != nil {
			t.Error(err)
		}

		f, err := os.Open("./testdata/test-1.png")
		if err != nil {
			t.Error(err)
		}

		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			t.Error(err)
		}

		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	// read the  pipe whihc recieves data
	request := httptest.NewRequest("POST", "/", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	var testTools ToolKit

	uploadedFile, err := testTools.UploadOneFile(request, "./testdata/uploads/", true)
	if err != nil {
		t.Error(err)
	}

	// check file is presnt actually or not
	_, err = os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFile.NewFileName))
	if os.IsNotExist(err) {
		t.Errorf("expected file to exist err: %v", err)
	}

	log.Println("Testcase 1: file is uploaded")
	// clean up
	_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFile.NewFileName))
}

func Test_DownloadFile(t *testing.T) {
	path := "./testdata"
	fileNmae := "test-1.png"

	tlk := ToolKit{}

	respWriter := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	tlk.DownloadStaticFile(respWriter, req, path, fileNmae, "sinoi.png")

	res := respWriter.Result()
	defer res.Body.Close()

	if res.Header["Content-Length"][0] != "1392758" {
		t.Error("Wrong content length: ", res.Header["Content-Length"][0])
	}

	if res.Header["Content-Disposition"][0] != "attachment; filename=\"sinoi.png\"" {
		t.Error("Wrong disposition: ", res.Header["Content-Disposition"][0])
	}

	// _, err := io.ReadAll(req.Body)
	// if err != nil {
	// 	t.Error(err)
	// }
}
