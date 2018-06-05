package blackbeard

import (
	"fmt"
)

//Inventory represents a group of variable to use in the templates.
// Namespace is the namespace dedicated files where to apply the variables contains into Values
// Values is map of string can contains whatever the user set in the default.json file inside a playbook
type Inventory struct {
	Namespace string                 `json:"namespace"`
	Values    map[string]interface{} `json:"values"`
}

//InventoryService define the way inventory should be managed.
type InventoryService interface {
	Create(namespace string) (Inventory, error)
	Update(namespace string, inventory Inventory) error
	Get(namespace string) (Inventory, error)
	GetDefault() (Inventory, error)
	List() ([]Inventory, error)
	Delete(namespace string) error
	Reset(namespace string) (Inventory, error)
}

type InventoryRepository interface {
	GetDefault() (Inventory, error)
	Get(namespace string) (Inventory, error)
	Create(Inventory) error
	Delete(namespace string) error
	Update(namespace string, inventory Inventory) error
	List() ([]Inventory, error)
}

type inventoryService struct {
	inventories InventoryRepository
}

func NewInventoryService(inventories InventoryRepository) InventoryService {
	return &inventoryService{
		inventories,
	}
}

//Create instantiate a new Inventory and write a json file containing the inventory
func (is *inventoryService) Create(namespace string) (Inventory, error) {

	if namespace == "" {
		return Inventory{}, fmt.Errorf("A namespace cannot be empty")
	}

	def, err := is.inventories.GetDefault()
	if err != nil {
		return Inventory{}, err
	}

	var inv Inventory
	inv.Namespace = namespace
	inv.Values = def.Values

	if err := is.inventories.Create(inv); err != nil {
		return Inventory{}, err
	}

	return inv, nil
}

func (is *inventoryService) Get(namespace string) (Inventory, error) {
	if namespace == "" {
		return Inventory{}, fmt.Errorf("A namespace cannot be empty")
	}

	return is.inventories.Get(namespace)
}

func (is *inventoryService) Delete(namespace string) error {
	return is.inventories.Delete(namespace)
}

func (is *inventoryService) GetDefault() (Inventory, error) {
	return is.inventories.GetDefault()
}

func (is *inventoryService) List() ([]Inventory, error) {
	return is.inventories.List()
}

func (is *inventoryService) Update(namespace string, inv Inventory) error {
	return is.inventories.Update(namespace, inv)
}

//Reset override the inventory file for the given namespace base on the content of the default inventory.
func (is *inventoryService) Reset(namespace string) (Inventory, error) {
	def, err := is.inventories.GetDefault()
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

type ErrorReadingDefaultsFile struct {
	msg string
}

func (err ErrorReadingDefaultsFile) Error() string {
	return err.msg
}

func NewErrorReadingDefaultsFile(err error) ErrorReadingDefaultsFile {
	return ErrorReadingDefaultsFile{fmt.Sprintf("Error when reading defaults file : %s", err.Error())}
}

type ErrorInventoryAlreadyExist struct {
	msg string
}

func (err ErrorInventoryAlreadyExist) Error() string {
	return err.msg
}

func NewErrorInventoryAlreadyExist(namespace string) ErrorInventoryAlreadyExist {
	return ErrorInventoryAlreadyExist{fmt.Sprintf("An inventory for the namespace %s already exist", namespace)}
}

type ErrorInventoryNotFound struct {
	msg string
}

func (err ErrorInventoryNotFound) Error() string {
	return err.msg
}

func NewErrorInventoryNotFound(namespace string) ErrorInventoryNotFound {
	return ErrorInventoryNotFound{fmt.Sprintf("The inventory for %s does not exist.", namespace)}
}
