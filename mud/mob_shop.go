package mud

import (
	"fmt"
	"strings"
)

func doBuy(player *mob, argument string) {
	if len(argument) == 0 {
		player.notify("Buy what?")
		return
	}

	if hasBit(player.Room.RoomFlags, roomPetShop) {
		// TODO
	} else {
		keeper := player.findKeeper()
		if keeper == nil {
			return
		}

		obj := keeper.carrying(argument)

		if obj == nil {
			act("$n tells you 'I don't sell that. Try 'list'.'", keeper, nil, player, actToVict)
			return
		}

		cost := keeper.getCost(obj, true)
		if cost <= 0 || !player.canSeeItem(obj) {
			act("$n tells you 'I don't sell that. Try 'list'.'", keeper, nil, player, actToVict)
			return
		}

		if player.Gold < cost {
			act("$n tells you 'You cannot afford to buy $p.'", keeper, obj, player, actToVict)
			player.replyTarget = keeper
			return
		}

		if obj.Level > player.Level {
			act("$n tells you 'You cannot use $p yet.", keeper, obj, player, actToVict)
			player.replyTarget = keeper
		}

		if player.Carrying+1 > player.CarryMax {
			player.notify("You can't carry that many items.")
			return
		}

		if player.CarryWeight+obj.Weight > player.CarryWeightMax {
			player.notify("You can't carry that much weight.")
			return
		}

		act("$n buys $p.", player, obj, nil, actToRoom)
		act("You buy $p.", player, obj, nil, actToChar)
		player.Gold -= cost
		keeper.Gold += cost

		var item *item
		if hasBit(obj.ExtraFlags, itemInventory) {
			item = createItem(&obj.index)
			item.Level = player.Level
		} else {
			item = obj
			keeper.removeItem(obj)
			return
		}
		player.addItem(item)
	}
}

func (player *mob) findKeeper() *mob {

	var keeper *mob
	for _, m := range player.Room.Mobs {
		if m.isNPC() && m.index.Shop != nil {
			keeper = m
			break
		}
	}

	if keeper == nil {
		player.notify("There is no shopkeep here.")
		return nil
	}

	store := keeper.index.Shop

	if store == nil {
		player.notify("You can't do that.")
		return nil
	}

	if store.isBeforeOpen() {
		doSay(keeper, "Sorry, come back later.")
		return nil
	}

	if store.isAfterClose() {
		doSay(keeper, "Sorry, come back tomorrow.")
		return nil
	}

	if !keeper.canSee(player) {
		doSay(keeper, "I don't trade with folks I can't see.")
		return nil
	}

	return keeper
}

func doList(player *mob, argument string) {

	if hasBit(player.Room.RoomFlags, roomPetShop) {
		// TODO
	} else {
		keeper := player.findKeeper()
		if keeper == nil {
			return
		}

		argument, arg1 := oneArgument(argument)
		found := false
		for _, i := range keeper.Inventory {
			cost := keeper.getCost(i, true)
			if i.WearLocation == wearNone && player.canSeeItem(i) && cost > 0 && len(argument) == 0 && matchesSubject(i.Name, arg1) {
				if !found {
					found = true
					player.notify("[Lv Price] Item")
				}

				player.notify("[%2d %5d] %s", i.Level, cost, strings.Title(i.Name))
			}
		}

		if !found {
			if len(argument) == 0 {
				player.notify("You can't buy anything here.")
			} else {
				player.notify("You can't buy that here.")
			}
		}
		return
	}
}

func doSell(player *mob, argument string) {
	if len(argument) < 1 {
		player.notify("Sell what?")
		return
	}

	keeper := player.findKeeper()
	if keeper == nil {
		return
	}

	argument, arg1 := oneArgument(argument)
	obj := player.carrying(arg1)

	if obj == nil {
		act("$n tells you 'You do not have that item.'", keeper, nil, player, actToVict)
		player.replyTarget = keeper
		return
	}

	if !player.canDropItem(obj) {
		player.notify("You can't let go of it.")
		return
	}

	cost := keeper.getCost(obj, false)
	if cost <= 0 {
		act("$n looks uninterested in $p.", keeper, obj, player, actToVict)
		return
	}

	act("$n sells $p.", player, obj, nil, actToRoom)
	suffix := "s"
	if cost == 1 {
		suffix = ""
	}
	act(fmt.Sprintf("You sell $p for %d gold piece%s.", cost, suffix), player, obj, nil, actToChar)

	player.Gold += cost
	keeper.Gold -= cost

	if keeper.Gold <= 0 {
		keeper.Gold = 0
	}

	if obj.ItemType == itemTrash {
		obj = nil // destroy item
	} else {
		player.removeItem(obj)
		keeper.addItem(obj)
		return
	}
}

func doValue(player *mob, argument string) {
	if len(argument) < 1 {
		player.notify("Value what?")
		return
	}

	keeper := player.findKeeper()

	if keeper == nil {
		return
	}

	argument, arg1 := oneArgument(argument)
	var obj *item
	for _, i := range player.Inventory {
		if matchesSubject(i.Name, arg1) {
			obj = i
			break
		}
	}

	if obj == nil {
		act("$n tells you 'You do not have that item.'", keeper, nil, player, actToVict)
		player.replyTarget = keeper
		return
	}

	if !player.canDropItem(obj) {
		player.notify("You can't let go of it.")
		return
	}

	cost := keeper.getCost(obj, false)
	if cost <= 0 {
		act("$n looks uninterested in $p.", keeper, obj, player, actToVict)
		return
	}

	buf := fmt.Sprintf("$n tells you \"I'll give you %d gold coins for $p.\"", cost)
	act(buf, keeper, obj, player, actToVict)
	player.replyTarget = keeper

	return
}
