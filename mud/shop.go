package mud

const maxTrade = 5

type shop struct {
	keeper     *mobIndex
	KeeperID   int           `json:"keeper_id"`
	BuyType    [maxTrade]int `json:"buy_types"`
	ProfitBuy  int           `json:"profit_buy"`
	ProfitSell int           `json:"profit_sell"`
	OpenHour   int           `json:"open_hour"`
	CloseHour  int           `json:"close_hour"`
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
			if item.ItemType == shop.BuyType[i] {
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

// TODO
func (s *shop) isBeforeOpen() bool {
	return false
}

// TODO
func (s *shop) isAfterClose() bool {
	return false
}
