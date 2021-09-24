package path

import (
	"testing"
)

func TestPreTrim(t *testing.T) {
	TestData := make(map[string]string)
	TestData["/abc/def"] = "abc/def"
	TestData["/abcdef"] = "abcdef"
	TestData["\\abc\\def"] = "abc\\def"
	TestData["\\abcdef"] = "abcdef"

	for k := range TestData {
		if result := PreTrim(k); result != TestData[k] {
			t.Errorf("PreTrim: expected '%s', got '%s'", TestData[k], result)
		}
	}
}

func TestPostTrim(t *testing.T) {
	TestData := make(map[string]string)

	TestData["mojotx/git-goclone.git"] = "mojotx/git-goclone"
	TestData["mojotx/git-goclone"] = "mojotx/git-goclone"
	for k := range TestData {
		if result := PostTrim(k); result != TestData[k] {
			t.Errorf("PostTrim: expected '%s', got '%s'", TestData[k], result)
		}
	}
}

func TestSanitize(t *testing.T) {
	TestData := make(map[string]string)

	TestData["/mojotx/git-goclone.git"] = "mojotx/git-goclone"
	TestData["mojotx/git-goclone.git"] = "mojotx/git-goclone"
	TestData["/mojotx/git-goclone"] = "mojotx/git-goclone"
	TestData["mojotx/git-goclone"] = "mojotx/git-goclone"

	TestData[`\mojotx\git-goclone.git`] = `mojotx\git-goclone`
	TestData[`mojotx\git-goclone.git`] = `mojotx\git-goclone`
	TestData[`\mojotx\git-goclone`] = `mojotx\git-goclone`
	TestData[`mojotx\git-goclone`] = `mojotx\git-goclone`
	for k := range TestData {
		if result := Sanitize(k); result != TestData[k] {
			t.Errorf("Sanitize: expected '%s', got '%s'", TestData[k], result)
		}
	}
}
