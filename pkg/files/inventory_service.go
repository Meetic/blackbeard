package files

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	inventoryFileSuffix = "inventory.json"
)

//Create instantiate a new Inventory and write a json file containing the inventory
func (is *InventoryService) Create(namespace string) (blackbeard.Inventory, error) {

	var inv blackbeard.Inventory

	if namespace == "" {
		return inv, fmt.Errorf("A namespace cannot be empty")
	}

	defaults, errR := ioutil.ReadFile(is.defaultsPath)

	if errR != nil {
		return inv, NewErrorReadingDefaultsFile(errR)
	}

	inv = blackbeard.NewInventory(namespace, defaults)

	//Check if an inventory file already exist for this namespace
	if is.exists(inv.Namespace) {
		return inv, NewErrorInventoryAlreadyExist(namespace)
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
			return NewErrorInventoryAlreadyExist(inv.Namespace)
		}
		err := os.Rename(is.path(namespace), is.path(inv.Namespace))
		if err != nil {
			return err
		}
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
		return inv, NewErrorInventoryNotFound(namespace)
	}

	return is.read(is.path(namespace))
}

//GetDefaults return the defaults value for an inventory
func (is *InventoryService) GetDefaults() (blackbeard.Inventory, error) {
	return is.read(is.defaultsPath)
}

//List return the list of existing inventories
//If no inventory file exist, the function returns an empty slice.
func (is *InventoryService) List() ([]blackbeard.Inventory, error) {
	var inventories []blackbeard.Inventory

	invFiles, _ := filepath.Glob(filepath.Join(is.inventoryPath, fmt.Sprintf("*_%s", inventoryFileSuffix)))

	for _, invFile := range invFiles {
		inv, err := is.read(invFile)
		if err != nil {
			return inventories, err
		}
		inventories = append(inventories, inv)
	}

	return inventories, nil

}

//Delete remove an inventory file.
//if the specified inventory does not exist, Delete return nil and does nothing.
func (is *InventoryService) Delete(namespace string) error {
	if !is.exists(namespace) {
		return nil
	}
	return os.Remove(is.path(namespace))
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

//exists return true if an inventory for the given namespace already exist.
//Else, it return false.
func (is *InventoryService) exists(namespace string) bool {
	if _, err := os.Stat(is.path(namespace)); os.IsNotExist(err) {
		return false
	} else if err == nil {
		return true
	}
	return false
}

//Path return the inventory file path of a given namespace
func (is *InventoryService) path(namespace string) string {
	return filepath.Join(is.inventoryPath, fmt.Sprintf("%s_%s", namespace, inventoryFileSuffix))
}
