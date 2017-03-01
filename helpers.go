package main

func RemoveItem(s []Item, index int) []Item {
	return append(s[:index], s[index+1:]...)
}
