package files

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
)

//InventoryService define the way inventory should be managed
type InventoryService struct {
	inventoryPath string
	defaultsPath  string
}

//Ensure that InventoryService implements the interface
var _ blackbeard.InventoryService = (*InventoryService)(nil)

const (
	inventoryFileSuffix = "_inventory.json"
)

//Create instantiate a new Inventory and write a json file containing the inventory
func (is *InventoryService) Create(namespace string) (blackbeard.Inventory, error) {

	var inv blackbeard.Inventory

	if namespace == "" {
		return inv, fmt.Errorf("A namespace cannot be empty")
	}

	defaults, errR := ioutil.ReadFile(is.defaultsPath)

	if errR != nil {
		log.Fatalf("Error when reading defaults file : %s", errR.Error())
	}

	inv = blackbeard.NewInventory(namespace, defaults)

	//Check if a inventory file already exist for this usr.
	if is.exists(inv.Namespace) {
		err := fmt.Errorf("An inventory for the namespace %s already exist", namespace)
		return inv, err
	}

	iJSON, _ := json.MarshalIndent(inv, "", "    ")

	err := ioutil.WriteFile(is.path(inv.Namespace), iJSON, 0644)

	if err != nil {
		return inv, err
	}

	return inv, nil
}

//Update will update inventory for a given namespace.
//If the namespace in the inventory is not the same as the namespace given as first parameters of Update
//this function will rename the inventory file to match ne new namespace.
func (is *InventoryService) Update(namespace string, inv blackbeard.Inventory) error {

	//check if the namespace name has change
	if namespace != inv.Namespace {
		//Check if a inventory file already exist for this usr.
		if is.exists(inv.Namespace) {
			err := fmt.Errorf("An inventory for the namespace %s already exist", inv.Namespace)
			return err
		}
		os.Rename(is.path(namespace), is.path(inv.Namespace))
	}

	iJSON, _ := json.MarshalIndent(inv, "", "    ")

	err := ioutil.WriteFile(is.path(inv.Namespace), iJSON, 0644)

	if err != nil {
		return err
	}

	return nil
}

//Get return an inventory for a given namespace.
//If the inventory cannot be found based on his path, Get return an empty inventory and an error
func (is *InventoryService) Get(namespace string) (blackbeard.Inventory, error) {

	var inv blackbeard.Inventory

	if !is.exists(namespace) {
		return inv, fmt.Errorf("The namespace %s does not exist.", namespace)
	}

	return is.read(is.path(namespace))
}

//GetDefaults return the defaults value for an inventory
func (is *InventoryService) GetDefaults() (blackbeard.Inventory, error) {
	return is.read(is.defaultsPath)
}

//List return the list of existing inventories
func (is *InventoryService) List() ([]blackbeard.Inventory, error) {
	var inventories []blackbeard.Inventory

	invFiles, _ := filepath.Glob(is.inventoryPath + "*" + inventoryFileSuffix)

	for _, invFile := range invFiles {
		inv, err := is.read(invFile)
		if err != nil {
			return inventories, err
		}
		inventories = append(inventories, inv)
	}

	return inventories, nil

}

func (is *InventoryService) read(path string) (blackbeard.Inventory, error) {
	var inv blackbeard.Inventory

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return inv, err
	}

	json.Unmarshal(raw, &inv)
	return inv, nil
}

//exists return true if an inventory for the given user already exist.
//Else, it return false.
func (is *InventoryService) exists(namespace string) bool {
	if _, err := os.Stat(is.path(namespace)); err == nil {
		return true
	}
	return false
}

//Path return the inventory file path of a given namespace
func (is *InventoryService) path(namespace string) string {
	return is.inventoryPath + "/" + namespace + inventoryFileSuffix
}
