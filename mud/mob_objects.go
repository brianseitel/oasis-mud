package mud

import (
	"fmt"

	"github.com/brianseitel/oasis-mud/helpers"
)

func (player *mob) get(item *item, container *item) {
	if !item.canWear(itemTake) {
		player.notify(fmt.Sprintf("You can't take that.%s", helpers.Newline))
		return
	}

	if player.Carrying+1 > player.CarryMax {
		player.notify(fmt.Sprintf("You can't carry that many items.%s", helpers.Newline))
		return
	}

	if player.CarryWeight+item.Weight > player.CarryWeightMax {
		player.notify(fmt.Sprintf("You can't carry that much weight.%s", helpers.Newline))
		return
	}

	if container != nil {
		player.notify(fmt.Sprintf("You get %s from %s.%s", item.Name, container.Name, helpers.Newline))
		player.Room.notify(fmt.Sprintf("%s gets %s from %s.%s", player.Name, item.Name, container.Name, helpers.Newline), player)
		container.removeObject(item)
	} else {
		player.notify(fmt.Sprintf("You get %s.%s", item.Name, helpers.Newline))
		player.Room.notify(fmt.Sprintf("%s gets %s.%s", player.Name, item.Name, helpers.Newline), player)
		player.Room.removeObject(item)
	}

	if item.ItemType == itemMoney {
		player.Gold += uint(item.Value)
		item.extract()
	} else {
		player.Inventory = append(player.Inventory, item)
	}
}
