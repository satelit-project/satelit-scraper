package anidb

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"gotest.tools/assert"
)

const _testdataDir = "testdata"

func TestParser_EmbeddedHTML(t *testing.T) {
	files, err := ioutil.ReadDir(_testdataDir)
	if err != nil {
		assert.Error(t, err, "Failed to read HTML dir")
	}

	tmpDir, err := ioutil.TempDir("", "anidbunzip")
	if err != nil {
		assert.Error(t, err, "Failed to create temp directory")
	}

	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), "zip") {
			continue
		}

		t.Logf("Testing %s", file.Name())
		zipPath := filepath.Join(_testdataDir, file.Name())
		err = unzip(zipPath, tmpDir)
		if err != nil {
			assert.Error(t, err, "Failed to unzip archive")
		}

		id, url := anidbZipInfo(file)
		testAnidbEntry(t, id, url, tmpDir)
	}
}

func testAnidbEntry(t *testing.T, id string, url *url.URL, dir string) {
	htmlPath := filepath.Join(dir, fmt.Sprintf("%s.html", id))
	jsonPath := filepath.Join(dir, fmt.Sprintf("%s.json", id))

	html, err := os.Open(htmlPath)
	if err != nil {
		assert.Error(t, err, "Failed to read HTML")
	}

	json, err := os.Open(jsonPath)
	if err != nil {
		assert.Error(t, err, "Failed to read JSON")
	}

	testAnidbData(t, url, html, json)
}

func testAnidbData(t *testing.T, url *url.URL, html io.Reader, json io.Reader) {
	parser, err := NewParser(url, html, nil)
	if err != nil {
		assert.Error(t, err, "Failed to parse html")
	}

	parsed, err := parser.Anime()
	if err != nil {
		assert.Error(t, err, "Failed to parse anime")
	}

	marshaler := jsonpb.Marshaler{}
	gotJSON, err := marshaler.MarshalToString(parsed)
	if err != nil {
		assert.Error(t, err, "Failed to generate Anime JSON")
	}

	expectedJSON, err := ioutil.ReadAll(json)
	if err != nil {
		assert.Error(t, err, "Failed to read Anime JSON")
	}

	assert.Assert(t, bytes.Equal([]byte(gotJSON), expectedJSON))
}

func anidbZipInfo(file os.FileInfo) (string, *url.URL) {
	rg := regexp.MustCompile(`\[(\d+)]`)
	id := rg.FindStringSubmatch(file.Name())[1]

	rawURL := fmt.Sprintf("https://anidb.net/anime/%s", id)
	parsed, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	return id, parsed
}

func unzip(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	defer func() {
		_ = reader.Close()
	}()

	for _, file := range reader.File {
		if err = unzipFile(file, dest); err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(file *zip.File, dest string) error {
	path := filepath.Join(dest, file.Name)
	fileReader, err := file.Open()
	if err != nil {
		return err
	}

	fileWriter, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}

	if _, err = io.Copy(fileWriter, fileReader); err != nil {
		return err
	}

	return nil
}
