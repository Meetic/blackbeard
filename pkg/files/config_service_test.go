package files_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"github.com/stretchr/testify/assert"
)

//Test Apply work as expected
func TestApplyOk(t *testing.T) {
	a := assert.New(t)

	defer cleanTestDir(t)

	//Create inventory
	def, _ := ioutil.ReadFile(filepath.Join(fixturesDir, defaultsFile))
	inv := blackbeard.NewInventory("test", def)
	inv.Namespace = "test"

	//Configure file client
	fClient := newDefaultClient()

	//Create folders
	createDefaultTestDir(t)
	createTemplateFiles(t, tplFile)

	//Test No Errors
	a.Nil(fClient.ConfigService().Apply(inv))

	//Test directory "test" exist under configsDir
	_, errD := os.Stat(filepath.Join(configsDir, "test"))
	a.Nil(errD)

	//Expected config file name generated from template
	ext := filepath.Ext(tplFile)
	_, genfile := filepath.Split(tplFile[0 : len(tplFile)-len(ext)])

	f, _ := ioutil.ReadDir(filepath.Join(configsDir, "test"))
	//Test only one file is created
	a.Equal(1, len(f))
	//Test file name is ok
	a.Equal(genfile, f[0].Name())

}

//Test apply an inventory when there are no template files
func TestApplyTemplateFileNotExist(t *testing.T) {
	a := assert.New(t)

	defer cleanTestDir(t)

	//Create inventory
	def, _ := ioutil.ReadFile(filepath.Join(fixturesDir, defaultsFile))
	inv := blackbeard.NewInventory("test", def)
	inv.Namespace = "test"

	//Configure file client
	fCLient := newDefaultClient()

	wrongTplFile := "template.yml.test"

	//Create folders
	createDefaultTestDir(t)
	createTemplateFiles(t, wrongTplFile)

	//est template error
	err := fCLient.ConfigService().Apply(inv)
	a.Equal(fmt.Sprintf("no template files found in directory %s", templateDir), err.Error())
}

//Test apply an inventory when config dir does not exist
func TestApplConfigDirNotExists(t *testing.T) {
	a := assert.New(t)

	defer cleanTestDir(t)

	//Create inventory
	def, _ := ioutil.ReadFile(filepath.Join(fixturesDir, defaultsFile))
	inv := blackbeard.NewInventory("test", def)
	inv.Namespace = "test"

	//Configure file client
	fCLient := newDefaultClient()

	//Create folders
	createTestDir(t, templateDir, os.ModePerm)
	createTemplateFiles(t, tplFile)

	//Test template error
	err := fCLient.ConfigService().Apply(inv)
	a.Equal(fmt.Sprintf("the configs dir '%s/%s' could not be created : mkdir %s/%s: no such file or directory", configsDir, "test", configsDir, "test"), err.Error())
}
