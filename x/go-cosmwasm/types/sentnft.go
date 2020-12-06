package types

import (
	"encoding/json"
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

type sentnftsJSON Sentnfts

// MarshalJSON implements a custom JSON marshaller for the Coins type to allow
// nil Coins to be encoded as an empty array.
func (sentnfts Sentnfts) MarshalJSON() ([]byte, error) {
	if sentnfts == nil {
		return json.Marshal(sentnftsJSON(Sentnfts{}))
	}

	return json.Marshal(sentnftsJSON(sentnfts))
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

