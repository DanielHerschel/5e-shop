package domain

type Campaign struct {
	Name        string `json:"name"`
	Shops       []Shop `json:"shops"`
	CurrentShop Shop `json:"current_shop"`

	Players []string `json:"players"`
}

type Shop struct {
	Name  string `json:"name"`
	Items []Item `json:"items"`
}

type Item struct {
	Name   string `json:"name"`
	Rarity Rarity `json:"rarity"`
	Copper int    `json:"copper"`
	Silver int    `json:"silver"`
	Gold   int    `json:"gold"`
	Plat   int    `json:"plat"`
}

type Rarity int

const (
	Common Rarity = iota
	Uncommon
	Rare
	VeryRare
	Legendary
)