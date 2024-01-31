package cryptutil

import "github.com/speps/go-hashids"

type Hash struct {
	secret string
	length int
}

func NewHashids(secret string, length int) *Hash {
	return &Hash{
		secret: secret,
		length: length,
	}
}

func (h *Hash) HashidsEncode(params []int) (string, error) {
	hd := hashids.NewData()
	hd.Salt = h.secret
	hd.MinLength = h.length

	hashID, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}

	hashStr, err := hashID.Encode(params)
	if err != nil {
		return "", err
	}

	return hashStr, nil
}

func (h *Hash) HashidsDecode(hash string) ([]int, error) {
	hd := hashids.NewData()
	hd.Salt = h.secret
	hd.MinLength = h.length

	hashID, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, err
	}

	ids, err := hashID.DecodeWithError(hash)
	if err != nil {
		return nil, err
	}

	return ids, nil
}
