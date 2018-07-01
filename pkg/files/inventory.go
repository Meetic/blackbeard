package files

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Meetic/blackbeard/pkg/playbook"
)

const (
	inventoryFileSuffix = "inventory.json"
)

type inventories struct {
	inventoryPath string
}

// NewInventoryRepository returns a InventoryRepository
// The parameter is the directory where are stored the inventories
func NewInventoryRepository(inventoryPath string) playbook.InventoryRepository {
	return &inventories{
		inventoryPath: inventoryPath,
	}
}

// Get returns an inventory for a given namespace.
// If the inventory cannot be found based on its path, Get returns an empty inventory and an error
func (ir *inventories) Get(namespace string) (playbook.Inventory, error) {

	if !ir.Exists(namespace) {
		return playbook.Inventory{}, playbook.NewErrorInventoryNotFound(namespace)
	}

	return ir.read(ir.path(namespace))
}

// Create writes an inventory file containing the inventory passed as parameter.
func (ir *inventories) Create(inventory playbook.Inventory) error {

	// Check if an inventory file already exist for this namespace
	if ir.Exists(inventory.Namespace) {
		return playbook.NewErrorInventoryAlreadyExist(inventory.Namespace)
	}

	j, _ := json.MarshalIndent(inventory, "", "    ")
	return ioutil.WriteFile(ir.path(inventory.Namespace), j, 0644)

}

// Delete remove an inventory file.
// if the specified inventory does not exist, Delete return nil and does nothing.
func (ir *inventories) Delete(namespace string) error {
	if !ir.Exists(namespace) {
		return nil
	}
	return os.Remove(ir.path(namespace))
}

// Update will update inventory for a given namespace.
// If the namespace in the inventory is not the same as the namespace given as first parameters of Update
// this function will rename the inventory file to match ne new namespace.
func (ir *inventories) Update(namespace string, inv playbook.Inventory) error {

	//check if the namespace name has change
	if namespace != inv.Namespace {
		//Check if a inventory file already exist for this usr.
		if ir.Exists(inv.Namespace) {
			return playbook.NewErrorInventoryAlreadyExist(inv.Namespace)
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
func (ir *inventories) List() ([]playbook.Inventory, error) {
	var inventories []playbook.Inventory

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

func (ir *inventories) read(path string) (playbook.Inventory, error) {
	var inv playbook.Inventory

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return inv, err
	}

	json.Unmarshal(raw, &inv)
	return inv, nil
}

// Exists return true if an inventory for the given namespace already exist.
// Else, it return false.
func (ir *inventories) Exists(namespace string) bool {
	if _, err := os.Stat(ir.path(namespace)); os.IsNotExist(err) {
		return false
	} else if err == nil {
		return true
	}
	return false
}

// path return the inventory file path of a given namespace
func (ir *inventories) path(namespace string) string {
	return filepath.Join(ir.inventoryPath, fmt.Sprintf("%s_%s", namespace, inventoryFileSuffix))
}
