package mock

import "github.com/Meetic/blackbeard/pkg/playbook"

type inventoryRepository struct{}

// NewInventoryRepository returns a Mock InventoryRepository
func NewInventoryRepository() playbook.InventoryRepository {
	return &inventoryRepository{}
}

func (ir *inventoryRepository) Get(namespace string) (playbook.Inventory, error) {
	playbooks := NewPlaybookRepository()
	inv, _ := playbooks.GetDefault()
	inv.Namespace = namespace

	return inv, nil
}

func (ir *inventoryRepository) Create(inventory playbook.Inventory) error {
	return nil
}

func (ir *inventoryRepository) Delete(namespace string) error {
	return nil
}

func (ir *inventoryRepository) Update(namespace string, inv playbook.Inventory) error {
	return nil
}

func (ir *inventoryRepository) Exists(namespace string) bool {
	return true
}

func (ir *inventoryRepository) List() ([]playbook.Inventory, error) {
	var inventories []playbook.Inventory

	inv1, _ := ir.Get("test1")
	inv2, _ := ir.Get("test2")

	inventories = append(inventories, inv1, inv2)

	return inventories, nil
}
