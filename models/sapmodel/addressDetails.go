package sapmodel

import (
	"encoding/json"
)

//AddressDetailsMarshal :
func (r *AddrReq) AddressDetailsMarshal() ([]byte, error) {
	a, err := json.Marshal(r)
	return a, err
}
