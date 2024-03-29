package toolkit

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	minL = 8
	maxL = 30
)

type ToolKit struct {
	MaxFileSize      int64
	AllowedFileTypes []string
}

type UploadedFile struct {
	NewFileName      string
	OriginalFileName string
	FileSize         int64
}

func (t *ToolKit) RandomPasswordGenerator() string {
	currTime := time.Now()

	sourceLetter := []string{
		"abcdefghijklomnopqrstuvwxyz",
		"ABCDEFGHIJKLOMNOPQRSTUVWXYZ",
		"@#$%^&*",
		"01233456789",
	}

	source := rand.NewSource(time.Now().UnixNano())

	randn := rand.New(source)

	passwordLen := minL + randn.Intn(maxL)

	password := []byte{}

	for {
		if len(password) == passwordLen {
			break
		}

		indexI := randn.Intn(4)
		indexJ := randn.Intn(len(sourceLetter[indexI]))
		password = append(password, sourceLetter[indexI][indexJ])
	}

	totalTimeTaken := time.Since(currTime)
	fmt.Printf("Total time taken to generate password: %v\n", totalTimeTaken)

	return string(password)
}

func (t *ToolKit) RandomStringGenerator(n int) string {
	sourceLetter := "abcdefghijklomnopqrstuvwxyzABCDEFGHIJKLOMNOPQRSTUVWXYZ01233456789"

	source := rand.NewSource(time.Now().UnixNano())

	randn := rand.New(source)

	str := []byte{}

	for {
		if len(str) == n {
			break
		}
		indexI := randn.Intn(len(sourceLetter))
		str = append(str, sourceLetter[indexI])
	}

	return string(str)
}

func (t *ToolKit) UploadFiles(r *http.Request, directory string, rename ...bool) ([]*UploadedFile, error) {
	renameFile := true

	if len(rename) > 0 {
		renameFile = rename[0]
	}

	if t.MaxFileSize == 0 {
		t.MaxFileSize = int64(math.Pow(2, 30))
	}
	uploadedFiles := make([]*UploadedFile, 0)

	err := r.ParseMultipartForm(t.MaxFileSize)
	if err != nil {
		return nil, fmt.Errorf("file is too large to upload")
	}

	// lets dive into the file
	for _, fHeaders := range r.MultipartForm.File {
		for _, hdr := range fHeaders {
			// lets create an inline function
			uploadedFiles, err = func(UploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
				uploadedFile := &UploadedFile{}

				// open contents
				infile, err := hdr.Open()
				if err != nil {
					return nil, err
				}

				defer infile.Close()

				// check type of file from first 512 bytes
				buff := make([]byte, 512)
				_, err = infile.Read(buff)
				if err != nil {
					return nil, err
				}

				// check type
				allowed := false
				fileType := http.DetectContentType(buff)

				if len(t.AllowedFileTypes) > 0 {
					for _, ftype := range t.AllowedFileTypes {
						if strings.EqualFold(fileType, ftype) {
							allowed = true
							break
						}
					}
				}

				if !allowed {
					return nil, fmt.Errorf("the uploaded file is of incorrect type")
				}

				// get back to starting of the file since ew are in line 512
				_, err = infile.Seek(0, 0)
				if err != nil {
					return nil, err
				}

				uploadedFile.OriginalFileName = hdr.Filename
				if renameFile {
					uploadedFile.NewFileName = fmt.Sprintf("%s%s", t.RandomStringGenerator(10), filepath.Ext(hdr.Filename))
				} else {
					uploadedFile.NewFileName = hdr.Filename
				}

				// save to disk
				var outFile *os.File
				defer outFile.Close()

				// create the output file in the dir
				outFile, err = os.Create(filepath.Join(directory, uploadedFile.NewFileName))
				if err != nil {
					return nil, err
				}

				// copy uploaded(input) file to output one inside directory
				filesize, err := io.Copy(outFile, infile)
				if err != nil {
					return nil, err
				}

				uploadedFile.FileSize = filesize
				uploadedFiles = append(uploadedFiles, uploadedFile)

				return uploadedFiles, nil
			}(uploadedFiles)

			if err != nil {
				return nil, err
			}
		}
	}

	return uploadedFiles, nil
}

func (t *ToolKit) UploadOneFile(r *http.Request, directory string, rename ...bool) (*UploadedFile, error) {
	files, err := t.UploadFiles(r, directory, rename...)
	if err != nil {
		return nil, fmt.Errorf("filed to upload file")
	}

	return files[0], nil
}
