package files

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
)

const (
	inventoryFileSuffix = "inventory.json"
)

type inventoryRepository struct {
	inventoryPath string
	defaultsPath  string
}

// NewInventoryRepository returns a InventoryRepository
// Parameters are the directory where are stored the inventories and the location of the default.json file
// aka : the default inventory
func NewInventoryRepository(inventoryPath string, defaultsPath string) blackbeard.InventoryRepository {
	return &inventoryRepository{
		inventoryPath: inventoryPath,
		defaultsPath:  defaultsPath,
	}
}

// GetDefault reads the default inventory file and return an Inventory where namespace is set to "default"
func (ir *inventoryRepository) GetDefault() (blackbeard.Inventory, error) {

	defaults, err := ioutil.ReadFile(ir.defaultsPath)

	if err != nil {
		return blackbeard.Inventory{}, blackbeard.NewErrorReadingDefaultsFile(err)
	}

	var inventory blackbeard.Inventory

	if err := json.Unmarshal(defaults, &inventory); err != nil {
		return blackbeard.Inventory{}, blackbeard.NewErrorReadingDefaultsFile(err)
	}

	return inventory, nil
}

// Get returns an inventory for a given namespace.
// If the inventory cannot be found based on its path, Get returns an empty inventory and an error
func (ir *inventoryRepository) Get(namespace string) (blackbeard.Inventory, error) {

	if !ir.exist(namespace) {
		return blackbeard.Inventory{}, blackbeard.NewErrorInventoryNotFound(namespace)
	}

	return ir.read(ir.path(namespace))
}

// Create writes an inventory file containing the inventory passed as parameter.
func (ir *inventoryRepository) Create(inventory blackbeard.Inventory) error {

	//Check if an inventory file already exist for this namespace
	if ir.exist(inventory.Namespace) {
		return blackbeard.NewErrorInventoryAlreadyExist(inventory.Namespace)
	}

	j, _ := json.MarshalIndent(inventory, "", "    ")
	return ioutil.WriteFile(ir.path(inventory.Namespace), j, 0644)

}

// Delete remove an inventory file.
// if the specified inventory does not exist, Delete return nil and does nothing.
func (ir *inventoryRepository) Delete(namespace string) error {
	if !ir.exist(namespace) {
		return nil
	}
	return os.Remove(ir.path(namespace))
}

// Update will update inventory for a given namespace.
// If the namespace in the inventory is not the same as the namespace given as first parameters of Update
// this function will rename the inventory file to match ne new namespace.
func (ir *inventoryRepository) Update(namespace string, inv blackbeard.Inventory) error {

	//check if the namespace name has change
	if namespace != inv.Namespace {
		//Check if a inventory file already exist for this usr.
		if ir.exist(inv.Namespace) {
			return blackbeard.NewErrorInventoryAlreadyExist(inv.Namespace)
		}
		err := os.Rename(ir.path(namespace), ir.path(inv.Namespace))
		if err != nil {
			return err
		}
	}

	iJSON, _ := json.MarshalIndent(inv, "", "    ")

	err := ioutil.WriteFile(ir.path(inv.Namespace), iJSON, 0644)

	if err != nil {
		return err
	}

	return nil
}

// List return the list of existing inventories
// If no inventory file exist, the function returns an empty slice.
func (ir *inventoryRepository) List() ([]blackbeard.Inventory, error) {
	var inventories []blackbeard.Inventory

	invFiles, _ := filepath.Glob(filepath.Join(ir.inventoryPath, fmt.Sprintf("*_%s", inventoryFileSuffix)))

	for _, invFile := range invFiles {
		inv, err := ir.read(invFile)
		if err != nil {
			return inventories, err
		}
		inventories = append(inventories, inv)
	}

	return inventories, nil

}

func (ir *inventoryRepository) read(path string) (blackbeard.Inventory, error) {
	var inv blackbeard.Inventory

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return inv, err
	}

	json.Unmarshal(raw, &inv)
	return inv, nil
}

// exist return true if an inventory for the given namespace already exist.
// Else, it return false.
func (ir *inventoryRepository) exist(namespace string) bool {
	if _, err := os.Stat(ir.path(namespace)); os.IsNotExist(err) {
		return false
	} else if err == nil {
		return true
	}
	return false
}

// Path return the inventory file path of a given namespace
func (ir *inventoryRepository) path(namespace string) string {
	return filepath.Join(ir.inventoryPath, fmt.Sprintf("%s_%s", namespace, inventoryFileSuffix))
}
