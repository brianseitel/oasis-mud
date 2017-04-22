package mud

import "testing"

func TestGetCost(t *testing.T) {
	keeper := &mob{Inventory: nil, index: &mobIndex{}}
	item := &item{ItemType: itemLight, Cost: 100}
	var buyTypes = [5]int{itemLight, 0, 0, 0, 0}
	shop := &shop{ProfitBuy: 115, ProfitSell: 115, BuyType: buyTypes}
	if keeper.getCost(item, true) != 0 {
		t.Error("Since keeper has no index, should be 0")
	}

	keeper.index = &mobIndex{Shop: shop}

	if keeper.getCost(item, true) != 115 {
		t.Error("Expected cost is 115")
	}

	if keeper.getCost(item, false) != 115 {
		t.Error("Expected cost is 115")
	}

	keeper.Inventory = append(keeper.Inventory, item)

	if keeper.getCost(item, false) != 0 {
		t.Error("Expected cost is 0")
	}
}

func TestIsBeforeOpen(t *testing.T) {
	shop := &shop{}

	if shop.isBeforeOpen() {
		t.Error("It's never before open")
	}
}

func TestIsAfterClose(t *testing.T) {
	shop := &shop{}

	if shop.isAfterClose() {
		t.Error("It's never after close")
	}
}
