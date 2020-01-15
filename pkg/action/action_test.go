package action

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestMakeLogWriter(t *testing.T) {
	testFileNamePrefix := filepath.Join(os.TempDir(), "TestMakeLogWriter")
	testFileContent := []byte("test makeLogWriter content")

	logName, logFile, err := makeLogWriter(testFileNamePrefix)
	if err != nil {
		t.Errorf("makeLogWriter function returned error: [%s]", err)
	}
	defer logFile.Close()

	n, err := logFile.Write(testFileContent)
	if err != nil {
		t.Errorf("error occurred while trying to write to the log file: [%s]", err)
	}

	if n != len(testFileContent) {
		t.Errorf("unexpected written bytes count: %d, expected: %d", n, len(testFileContent))
	}

	readContent, err := ioutil.ReadFile(logName)
	if err != nil {
		t.Errorf("error occurred while trying to read from the log file: [%s]", err)
	}

	if equal := reflect.DeepEqual(readContent, testFileContent); equal != true {
		t.Errorf("recorded content does not match read content")
	}

	// still need to close explicitly instead of relying on defer so os.Remove succeeds
	logFile.Close()

	if err = os.Remove(logName); err != nil {
		t.Errorf("error deleting file; file: [%s], error: [%s]", logName, err)
	}
}

func TestLazyFile(t *testing.T) {
	testFileNamePrefix := filepath.Join(os.TempDir(), "TestLazyFile")
	testFileContent := []byte("test lazy file content")

	lzFile := &lazyFile{prefix: testFileNamePrefix}

	writtenBytes, err := lzFile.Write(testFileContent)
	if err != nil {
		t.Errorf("error occurred while trying to write to the lazy file: [%s]", err)
	}
	defer lzFile.Close()

	// after first Write call, out should be not nil
	if lzFile.out == nil {
		t.Errorf("something went wrong")
	}

	if writtenBytes != len(testFileContent) {
		t.Errorf("unexpected written bytes count: %d, expected: %d", writtenBytes, len(testFileContent))
	}

	readContent, err := ioutil.ReadFile(lzFile.name)
	if err != nil {
		t.Errorf("error occurred while trying to read from the lazy file: [%s]", err)
	}

	if equal := reflect.DeepEqual(readContent, testFileContent); equal != true {
		t.Errorf("recorded content does not match read content")
	}

	// still need to close explicitly instead of relying on defer so os.Remove succeeds
	lzFile.Close()

	// after call Close method, Write method should return error != nil
	_, err = lzFile.Write(testFileContent)
	if err == nil {
		t.Errorf("Close method does not actually close the file")
	}

	if err = os.Remove(lzFile.name); err != nil {
		t.Errorf("error deleting file; file: [%s], error: [%s]", lzFile.name, err)
	}
}

func TestRunLoggedCmdDirOutput(t *testing.T) {
	testLogFilePrefix := filepath.Join(os.TempDir(), "TestRunLoggedCmdDirOutput")
	var buffer bytes.Buffer
	workdir := ""
	name := "go"
	arg := "version"

	logName, err := RunLoggedCmdDirOutput(testLogFilePrefix, workdir, &buffer, name, arg)
	if err != nil {
		t.Errorf("error occurred while trying to run command: [%s %s], error: [%s]", name, arg, err)
	}

	content := buffer.String()
	expectedBufferContent := "go version go1.13.3 windows/amd64"
	expectedLogFileContent := "Rhoc: running command: go version"

	if strings.HasPrefix(content, expectedBufferContent) != true {
		t.Errorf("recorded content: [%s] does not match read content: [%s]", content, expectedBufferContent)
	}

	logContent, err := ioutil.ReadFile(logName)
	if err != nil {
		t.Errorf("error occurred while trying to read from the log file: [%s]", err)
	}

	if strings.HasPrefix(string(logContent), expectedLogFileContent) != true {
		t.Errorf("recorded content: [%s] does not match read content: [%s]", content, expectedBufferContent)
	}

	logName2, err := RunLoggedCmd(testLogFilePrefix, name, arg)
	if err != nil {
		t.Errorf("error occurred while trying to run command: [%s %s], error: [%s]", name, arg, err)
	}

	if logName != logName2 {
		t.Errorf("something went wrong")
	}

}
