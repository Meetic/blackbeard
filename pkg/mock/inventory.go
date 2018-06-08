package mock

import (
	"github.com/Meetic/blackbeard/pkg/blackbeard"
)

type inventoryRepository struct{}

// NewInventoryRepository returns a Mock InventoryRepository
func NewInventoryRepository() blackbeard.InventoryRepository {
	return &inventoryRepository{}
}

func (ir *inventoryRepository) Get(namespace string) (blackbeard.Inventory, error) {
	playbooks := NewPlaybookRepository()
	inv, _ := playbooks.GetDefault()
	inv.Namespace = namespace

	return inv, nil
}

func (ir *inventoryRepository) Create(inventory blackbeard.Inventory) error {
	return nil
}

func (ir *inventoryRepository) Delete(namespace string) error {
	return nil
}

func (ir *inventoryRepository) Update(namespace string, inv blackbeard.Inventory) error {
	return nil
}

func (ir *inventoryRepository) List() ([]blackbeard.Inventory, error) {
	var inventories []blackbeard.Inventory

	inv1, _ := ir.Get("test1")
	inv2, _ := ir.Get("test2")

	inventories = append(inventories, inv1, inv2)

	return inventories, nil
}
