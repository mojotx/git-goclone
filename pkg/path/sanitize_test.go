package path

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreTrim(t *testing.T) {
	TestData := make(map[string]string)
	TestData["/abc/def"] = "abc/def"
	TestData["/abcdef"] = "abcdef"
	TestData["\\abc\\def"] = "abc\\def"
	TestData["\\abcdef"] = "abcdef"

	for k := range TestData {
		result := PreTrim(k)
		assert.Equalf(t, result, TestData[k], "PreTrim: expected '%s', got '%s'", TestData[k], result)
	}
}

func TestPostTrim(t *testing.T) {
	TestData := make(map[string]string)

	TestData["mojotx/git-goclone.git"] = "mojotx/git-goclone"
	TestData["mojotx/git-goclone"] = "mojotx/git-goclone"
	for k := range TestData {
		result := PostTrim(k)
		assert.Equalf(t, result, TestData[k], "PostTrim: expected '%s', got '%s'", TestData[k], result)
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
		result := Sanitize(k)
		assert.Equalf(t, result, TestData[k], "Sanitize: expected '%s', got '%s'", TestData[k], result)
	}
}
