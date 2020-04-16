package playbook

import (
	"fmt"
)

// Inventory represents a set of variable to apply to the templates (see config).
// Namespace is the namespace dedicated files where to apply the variables contains into Values
// Values is map of string that contains whatever the user set in the default inventory from a playbook
type Inventory struct {
	Namespace string                 `json:"namespace"`
	Values    map[string]interface{} `json:"values"`
}

// InventoryService define the way inventories are managed.
type InventoryService interface {
	Create(namespace string) (Inventory, error)
	Update(namespace string, inventory Inventory) error
	Get(namespace string) (Inventory, error)
	Exists(namespace string) bool
	List() ([]Inventory, error)
	Delete(namespace string) error
	Reset(namespace string) (Inventory, error)
}

// InventoryRepository define the way inventories are actually managed
type InventoryRepository interface {
	Get(namespace string) (Inventory, error)
	Exists(namespace string) bool
	Create(Inventory) error
	Delete(namespace string) error
	Update(namespace string, inventory Inventory) error
	List() ([]Inventory, error)
}

type inventoryService struct {
	inventories InventoryRepository
	playbooks   PlaybookService
}

// NewInventoryService create an InventoryService
func NewInventoryService(inventories InventoryRepository, playbooks PlaybookService) InventoryService {
	return &inventoryService{
		inventories,
		playbooks,
	}
}

// Create instantiate a new Inventory from the default inventory of a playbook and save it
func (is *inventoryService) Create(namespace string) (Inventory, error) {

	if namespace == "" {
		return Inventory{}, fmt.Errorf("A namespace cannot be empty")
	}

	def, err := is.playbooks.GetDefault()
	if err != nil {
		return Inventory{}, err
	}

	inv := Inventory{
		Namespace: namespace,
		Values:    def.Values,
	}

	if err := is.inventories.Create(inv); err != nil {
		return Inventory{}, err
	}

	return inv, nil
}

// Get returns the Inventory for a given namespace
func (is *inventoryService) Get(namespace string) (Inventory, error) {
	if namespace == "" {
		return Inventory{}, fmt.Errorf("A namespace cannot be empty")
	}

	return is.inventories.Get(namespace)
}

// Exists return true if an inventory for the given namespace already exists.
// Else, it return false.
func (is *inventoryService) Exists(namespace string) bool {
	return is.inventories.Exists(namespace)
}

// Delete deletes the inventory for the given namespace
func (is *inventoryService) Delete(namespace string) error {
	return is.inventories.Delete(namespace)
}

// List returns the list of available inventories
func (is *inventoryService) List() ([]Inventory, error) {
	return is.inventories.List()
}

// Update replace the inventory associated to the given namespace by the given inventory
func (is *inventoryService) Update(namespace string, inv Inventory) error {
	return is.inventories.Update(namespace, inv)
}

// Reset override the inventory file for the given namespace base on the content of the default inventory.
func (is *inventoryService) Reset(namespace string) (Inventory, error) {
	def, err := is.playbooks.GetDefault()
	if err != nil {
		return Inventory{}, err
	}

	var inv Inventory

	inv.Namespace = namespace
	inv.Values = def.Values

	if err := is.inventories.Update(namespace, inv); err != nil {
		return Inventory{}, err
	}

	return inv, nil
}

// ErrorReadingDefaultsFile represents an error due to unreadable default inventory
type ErrorReadingDefaultsFile struct {
	msg string
}

// Error returns the error message
func (err ErrorReadingDefaultsFile) Error() string {
	return err.msg
}

// NewErrorReadingDefaultsFile creates an ErrorReadingDefaultsFile error
func NewErrorReadingDefaultsFile(err error) ErrorReadingDefaultsFile {
	return ErrorReadingDefaultsFile{fmt.Sprintf("Error when reading defaults file : %s", err.Error())}
}

// ErrorInventoryAlreadyExist represents an error due to an already existing inventory for a given namespace
type ErrorInventoryAlreadyExist struct {
	msg string
}

// Error returns the error message
func (err ErrorInventoryAlreadyExist) Error() string {
	return err.msg
}

// NewErrorInventoryAlreadyExist creates a new ErrorInventoryAlreadyExist error
func NewErrorInventoryAlreadyExist(namespace string) ErrorInventoryAlreadyExist {
	return ErrorInventoryAlreadyExist{fmt.Sprintf("An inventory for the namespace %s already exist", namespace)}
}

// ErrorInventoryNotFound represents an error due to a missing inventory for the given namespace
type ErrorInventoryNotFound struct {
	msg string
}

// Error returns the error message
func (err ErrorInventoryNotFound) Error() string {
	return err.msg
}

// NewErrorInventoryNotFound creates a new ErrorInventoryNotFound error
func NewErrorInventoryNotFound(namespace string) ErrorInventoryNotFound {
	return ErrorInventoryNotFound{fmt.Sprintf("The inventory for %s does not exist.", namespace)}
}
