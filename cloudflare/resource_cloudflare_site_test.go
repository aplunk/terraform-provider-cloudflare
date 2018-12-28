package cloudflare

import (
	"io/ioutil"
	"testing"
)

func TestUploadSite(t *testing.T) {
	testDir, err := ioutil.TempDir("", "terraform-upload-test")
	if err != nil {
		t.Fatal(err)
	}

	first, err := ioutil.TempFile(testDir, "one")
	if err != nil {
		t.Fatal(err)
	}

	fContent := make([]byte, 30)
	first.Write(fContent)

	second, err := ioutil.TempFile(testDir, "two")
	if err != nil {
		t.Fatal(err)
	}

	test := func(key string, value []byte) error {
		return nil
	}
}
