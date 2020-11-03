package types

import (
	"strings"
)

type Sentnft struct {
	Denom  string `json:"denom"`
	Id string `json:"id"`
}

type Sentnfts []Sentnft

func (Sentnfts Sentnfts) Empty() bool {
	return len(Sentnfts) == 0
}

func NewNft(denom string, id string) Sentnft {
	//if err := validate(denom, amount); err != nil {
	//	panic(err)
	//}

	return Sentnft{
		Denom:  denom,
		Id: id,
	}
}

//Parse 1 Nft as sentnft[]
func ParseNft(nftStr string) (Sentnfts, error) {
	nftStr = strings.TrimSpace(nftStr)
	if len(nftStr) == 0 {
		return nil, nil
	}

	nftStrs := strings.Split(nftStr, ",")
	denom := nftStrs[0]
	id := nftStrs[1]
	sentnft := NewNft(denom, id)

	sentnfts := make(Sentnfts, 1)
	sentnfts[0] = sentnft

	return sentnfts, nil
}

