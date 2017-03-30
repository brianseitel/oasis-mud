package mud

const maxTrade = 5

type shop struct {
	keeper     *mobIndex
	KeeperID   int
	BuyType    [maxTrade]int
	ProfitBuy  int
	ProfitSell int
	OpenHour   int
	CloseHour  int
}

func (keeper *mob) getCost(item *item, buy bool) int {
	if item == nil || keeper.index.Shop == nil {
		return 0
	}
	var cost int

	shop := keeper.index.Shop
	if buy {
		cost = item.Cost * shop.ProfitBuy / 100
	} else {
		cost = 0

		for i := 0; i < maxTrade; i++ {
			if item.ItemType == uint(shop.BuyType[i]) {
				cost = item.Cost * shop.ProfitSell / 100
				break
			}
		}

		for _, i := range keeper.Inventory {
			if i.index.ID == item.index.ID {
				cost = 0
				break
			}
		}
	}

	return cost
}

func (s *shop) isBeforeOpen() bool {
	return false
}

func (s *shop) isAfterClose() bool {
	return false
}
