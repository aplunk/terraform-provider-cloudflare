package cloudflare

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

type testFile struct {
	path    []string
	name    string
	size    int
	pathStr string
	key     string
	data    []byte
}

func (t *testFile) make() error {
	t.pathStr = path.Join(os.TempDir(), path.Join(t.path...))
	if err := os.MkdirAll(t.pathStr, 0700); err != nil {
		return err
	}

	fh, err := ioutil.TempFile(t.pathStr, t.name)
	if err != nil {
		return err
	}
	defer fh.Close()

	t.key = strings.Replace(fh.Name(), string(os.PathSeparator), "_", -1)

	t.data = make([]byte, t.size)
	_, err = fh.Write(t.data)
	return err
}

func (t *testFile) cleanup() {
	os.RemoveAll(t.pathStr)
}

func TestUploadSite(t *testing.T) {
	tests := []*testFile{
		&testFile{
			[]string{"terraform-site-test"},
			"one",
			50,
			"",
			"",
			nil,
		},
		&testFile{
			[]string{"terraform-site-test"},
			"two",
			49,
			"",
			"",
			nil,
		},
		&testFile{
			[]string{"terraform-site-test", "nested"},
			"three",
			30,
			"",
			"",
			nil,
		},
	}

	keys := make([]string, len(tests))
	for i, setup := range tests {
		if err := setup.make(); err != nil {
			t.Fatal(err)
		}
		defer setup.cleanup()
		keys[i] = setup.key
	}

	test := func(key string, value []byte) error {
		for _, test := range tests {
			if test.key == key && bytes.Equal(test.data, value) {
				return nil
			}
			t.Logf("received key:%s value:%d", key, len(value))
			t.Logf("skipping key:%s value:%d", test.key, len(test.data))
		}
		t.Fatalf("key:%s value:%d not in %+v", key, len(value), keys)
		return nil
	}

	_, err := uploadSite("test_namespace", path.Join(os.TempDir(), "terraform-site-test"), 49, test)
	if err != nil && err != io.EOF {
		t.Fatalf("%+v", err)
	}
}
