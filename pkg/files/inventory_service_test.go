package files_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"os"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"github.com/Meetic/blackbeard/pkg/files"
	"github.com/stretchr/testify/assert"
)

//Test Create method works as expected
func TestCreateInventoryOK(t *testing.T) {
	a := assert.New(t)
	defer cleanTestDir(t)

	namespace := "test"
	createDefaultTestDir(t)
	fClient := newDefaultClient()

	inv, err := fClient.InventoryService().Create(namespace)

	//Assert namespace is ok
	a.Equal(namespace, inv.Namespace)
	a.Nil(err)

	f, _ := ioutil.ReadDir(inventoryDir)
	//Test only one file is created
	a.Equal(1, len(f))
	//Test file name is ok
	a.Equal(fmt.Sprintf("%s_%s", namespace, "inventory.json"), f[0].Name())

}

//Test Create method with different values for namespace
func TestCreateCheckNamespace(t *testing.T) {
	a := assert.New(t)
	defer cleanTestDir(t)

	cases := []struct {
		namespace string
		err       bool
		errMsg    string
	}{
		{"", true, "A namespace cannot be empty"},
		{"test", false, ""},
		{"12345", false, ""},
	}

	createDefaultTestDir(t)

	for _, c := range cases {

		fClient := newDefaultClient()
		inv, err := fClient.InventoryService().Create(c.namespace)

		//Assert namespace is ok
		if c.err {
			a.Error(err)
			a.Equal(c.errMsg, err.Error())
		} else {
			a.Nil(err)
			a.Equal(c.namespace, inv.Namespace)
		}
	}
}

//Test Create when defaults.json file does not exist
func TestDefaultNotExist(t *testing.T) {
	a := assert.New(t)
	defer cleanTestDir(t)

	namespace := "test"
	createDefaultTestDir(t)
	fClient := files.NewClient(templateDir, configsDir, inventoryDir, "test")

	_, err := fClient.InventoryService().Create(namespace)

	if a.Error(err) {
		a.IsType(files.ErrorReadingDefaultsFile{}, err)
	}
}

//Test create inventory when an inventory named as the specified one already exist
func TestInventoryAlreadyExist(t *testing.T) {
	a := assert.New(t)
	defer cleanTestDir(t)

	namespace := "test"
	createDefaultTestDir(t)
	fClient := newDefaultClient()

	_, _ = fClient.InventoryService().Create(namespace)

	inv, err := fClient.InventoryService().Create(namespace)

	//Assert namespace is ok
	a.Equal(namespace, inv.Namespace)
	a.NotNil(err)
	a.IsType(files.ErrorInventoryAlreadyExist{}, err)
}

//Test Update func work as expected when trying to update the namespace
func TestUpdateInventoryOK(t *testing.T) {
	a := assert.New(t)
	defer cleanTestDir(t)

	namespace := "test"

	//Create a new inventory to replace the first one
	newInv := blackbeard.NewInventory("test2", []byte(`
{
  "namespace": "test2",
  "values": {
    "microservices": [
      {
        "name": "api-advertising",
        "version": "latest",
        "urls": [
          "api-advertising"
        ]
      },
      {
        "name": "api-algo",
        "version": "latest",
        "urls": [
          "api-algo"
        ]
      }
    ]
  }
}
`))

	createDefaultTestDir(t)
	fClient := newDefaultClient()

	//create a first inventory
	_, _ = fClient.InventoryService().Create(namespace)
	//update the inventory using the new created inventory
	err := fClient.InventoryService().Update(namespace, newInv)

	a.Nil(err)
	f, _ := ioutil.ReadDir(inventoryDir)
	//Test only one file is created
	a.Equal(1, len(f))
	//Test file name is ok
	a.Equal(fmt.Sprintf("%s_%s", "test2", "inventory.json"), f[0].Name())
}

//Test get an inventory
func TestGetInventory(t *testing.T) {
	a := assert.New(t)
	defer cleanTestDir(t)

	createDefaultTestDir(t)
	//Create an inventory file to get
	def, _ := ioutil.ReadFile(filepath.Join(fixturesDir, defaultsFile))
	_ = ioutil.WriteFile(filepath.Join(inventoryDir, "default_inventory.json"), def, 0644)

	fClient := newDefaultClient()
	inv, err := fClient.InventoryService().Get("default")

	a.Equal("default", inv.Namespace)
	a.Nil(err)
}

//Test Get func when inventory does not exist
func TestGetInventoryNotExist(t *testing.T) {
	a := assert.New(t)
	defer cleanTestDir(t)

	createDefaultTestDir(t)

	fClient := newDefaultClient()
	_, err := fClient.InventoryService().Get("default")

	a.NotNil(err)
	a.IsType(files.ErrorInventoryNotFound{}, err)
}

//Test GetDefault func
func TestGetDefaults(t *testing.T) {
	a := assert.New(t)
	defer cleanTestDir(t)

	createDefaultTestDir(t)

	fClient := newDefaultClient()
	inv, err := fClient.InventoryService().GetDefaults()

	a.Equal("default", inv.Namespace)
	a.Nil(err)
}

//Test GetDefault func when defaults file doesn't exist
func TestGetDefaultsNotFound(t *testing.T) {
	a := assert.New(t)
	defer cleanTestDir(t)

	createDefaultTestDir(t)

	fClient := files.NewClient(templateDir, configsDir, inventoryDir, "test")

	_, err := fClient.InventoryService().GetDefaults()

	a.NotNil(err)
}

//Test get a list of inventories
func TestGetInventoryList(t *testing.T) {
	a := assert.New(t)
	defer cleanTestDir(t)

	createDefaultTestDir(t)
	//Create an inventory file to get
	def, _ := ioutil.ReadFile(filepath.Join(fixturesDir, defaultsFile))
	_ = ioutil.WriteFile(filepath.Join(inventoryDir, "default_inventory.json"), def, 0644)
	_ = ioutil.WriteFile(filepath.Join(inventoryDir, "default2_inventory.json"), def, 0644)

	fClient := newDefaultClient()
	inventories, err := fClient.InventoryService().List()

	a.Nil(err)

	a.Len(inventories, 2)

	for _, inv := range inventories {
		a.Equal("default", inv.Namespace)
	}
}

//Test get a list of inventories when there are no inventories.
func TestGetInventoryListNoFiles(t *testing.T) {
	a := assert.New(t)
	defer cleanTestDir(t)

	createDefaultTestDir(t)

	fClient := newDefaultClient()
	inventories, err := fClient.InventoryService().List()

	a.Nil(err)

	a.Len(inventories, 0)
	a.Empty(inventories)
}

//Test Delete method works as expected
func TestDeleteInventoryOK(t *testing.T) {
	a := assert.New(t)
	defer cleanTestDir(t)

	namespace := "test"
	createDefaultTestDir(t)
	fClient := newDefaultClient()

	fClient.InventoryService().Create(namespace)

	a.Nil(fClient.InventoryService().Delete(namespace))
	//Test directory "test" no more exit
	_, errD := os.Stat(filepath.Join(inventoryDir, fmt.Sprintf("%s_%s", namespace, "inventory.json")))
	a.True(os.IsNotExist(errD))
}

//Test Delete method return nil when inventory does not exist
func TestDeleteInventoryNotExist(t *testing.T) {
	a := assert.New(t)
	defer cleanTestDir(t)

	namespace := "test"
	createDefaultTestDir(t)
	fClient := newDefaultClient()
	a.Nil(fClient.InventoryService().Delete(namespace))
	//Test directory "test" no more exit
	_, errD := os.Stat(filepath.Join(inventoryDir, fmt.Sprintf("%s_%s", namespace, "inventory.json")))
	a.True(os.IsNotExist(errD))
}
