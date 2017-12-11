package files_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Meetic/blackbeard/pkg/files"
)

const (
	fixturesDir  = "test-fixtures"
	templateDir  = "template"
	configsDir   = "configs"
	inventoryDir = "inventories"
	defaultsFile = "defaults.json"
	tplFile      = "template.yml.tpl"
)

//-------------------
// Test Helper Functions
//-------------------

func createTemplateFiles(t *testing.T, tplFile string) {

	//Read fixtures
	tplString, _ := ioutil.ReadFile(filepath.Join(fixturesDir, tplFile))

	tFile := filepath.Join(templateDir, tplFile)
	if err := ioutil.WriteFile(tFile, tplString, 0644); err != nil {
		t.Fatalf("Template file %s could not be created", tFile)
	}

}

func createTestDir(t *testing.T, dir string, perm os.FileMode) {
	if err := os.Mkdir(dir, perm); err != nil {
		t.Fatalf("Directory %s could not be created for tests", dir)
	}
}

func createDefaultTestDir(t *testing.T) {
	dirs := []string{templateDir, inventoryDir, configsDir}
	for _, d := range dirs {
		createTestDir(t, d, os.ModePerm)
	}
}

func cleanTestDir(t *testing.T) {
	dirs := []string{templateDir, inventoryDir, configsDir}
	for _, d := range dirs {
		if err := os.RemoveAll(d); err != nil {
			t.Fatalf("Directory %s could not be removed", d)
		}
	}
}

//create a files client using default value for tests
func newDefaultClient() *files.Client {
	return files.NewClient(templateDir, configsDir, inventoryDir, filepath.Join(fixturesDir, defaultsFile))
}
