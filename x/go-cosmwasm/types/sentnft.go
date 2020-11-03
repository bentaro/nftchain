package types

type Sentnft struct {
	Denom  string `json:"denom"`
	Id string `json:"id"`
}

type Sentnfts []Sentnft

func (Sentnfts Sentnfts) Empty() bool {
	return len(Sentnfts) == 0
}
