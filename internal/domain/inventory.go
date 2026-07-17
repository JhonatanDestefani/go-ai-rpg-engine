package domain

type Inventory struct {
	Items    []Item
	Capacity int
}

func NewInventory(capacity int) Inventory {
	return Inventory{
		Items:    []Item{},
		Capacity: capacity,
	}
}

func (i *Inventory) AddItem(item Item) bool {
	for index := range i.Items {
		if i.Items[index].ID == item.ID {
			i.Items[index].Quantity += item.Quantity
			return true
		}
	}

	if len(i.Items) >= i.Capacity {
		return false
	}

	i.Items = append(i.Items, item)
	return true
}

func (i *Inventory) RemoveItem(itemID string, quantity int) bool {
	for index := range i.Items {
		if i.Items[index].ID != itemID {
			continue
		}

		if i.Items[index].Quantity < quantity {
			return false
		}

		i.Items[index].Quantity -= quantity

		if i.Items[index].Quantity <= 0 {
			i.Items = append(i.Items[:index], i.Items[index+1:]...)
		}

		return true
	}

	return false
}

func (i *Inventory) FindItem(itemID string) (*Item, bool) {
	for index := range i.Items {
		if i.Items[index].ID == itemID {
			return &i.Items[index], true
		}
	}

	return nil, false
}

func (inventory *Inventory) GetUsableItems() []Item {
	usableItems := []Item{}

	for _, item := range inventory.Items {
		if item.Type != ItemConsumable {
			continue
		}

		if item.Quantity <= 0 {
			continue
		}

		usableItems = append(usableItems, item)
	}

	return usableItems
}

func (i *Inventory) GetEquippableItems() []Item {
	equippableItems := []Item{}

	for _, item := range i.Items {
		if item.Quantity <= 0 {
			continue
		}

		if item.Type != ItemWeapon && item.Type != ItemArmor {
			continue
		}

		equippableItems = append(equippableItems, item)
	}

	return equippableItems
}
