package sapmodel

import (
	"encoding/json"
)

//ProductDetailsMarshal :
func (r *ProductDetailsReq) ProductDetailsMarshal() ([]byte, error) {
	a, err := json.Marshal(r)
	return a, err
}
