package main

import (
	//"os"
	"testing"

	"github.com/kumarmapanip/toolkit"
	//"github.com/stretchr/testify/assert"
)

var TestCases = []struct {
	dirPath       string
	errorExpected bool
}{
	{dirPath: "./logs", errorExpected: false},
	{dirPath: "./logs", errorExpected: true},
}

func Test_CreateDirectory(t *testing.T) {
	toolKit := toolkit.ToolKit{}
	for _, testCase := range TestCases {
		dirPath := testCase.dirPath
		err := toolKit.CreateDirectoryIfNotExists(dirPath)

		if !testCase.errorExpected && err != nil {
			t.Errorf("failed test validation %s", testCase.dirPath)
		}

		if testCase.errorExpected && err == nil {
			t.Errorf("failed test validation %s", testCase.dirPath)
		}
		// one more test
		// fileInfo, err := os.Stat(dirPath)
		// assert.Equal(t, fileInfo != nil, true)
		// if err != nil {
		// 	t.Error(err)
		// }

		// cleanup
		//_ = os.RemoveAll(dirPath)
	}
}
