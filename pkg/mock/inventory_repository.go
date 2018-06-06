package mock

import (
	"encoding/json"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
)

const (
	def = `{
  "namespace": "default",
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
}`
)

type inventoryRepository struct{}

// NewInventoryRepository returns a Mock InventoryRepository
func NewInventoryRepository() blackbeard.InventoryRepository {
	return &inventoryRepository{}
}

func (ir *inventoryRepository) GetDefault() (blackbeard.Inventory, error) {

	var inventory blackbeard.Inventory

	if err := json.Unmarshal([]byte(def), &inventory); err != nil {
		return blackbeard.Inventory{}, blackbeard.NewErrorReadingDefaultsFile(err)
	}

	return inventory, nil
}

func (ir *inventoryRepository) Get(namespace string) (blackbeard.Inventory, error) {

	inv, _ := ir.GetDefault()
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
