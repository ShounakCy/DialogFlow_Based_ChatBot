package sapmodel

import (
	"encoding/json"
)

//AddressBookMarshal :
func (r *AddressReq) AddressBookMarshal() ([]byte, error) {
	a, err := json.Marshal(r)
	return a, err
}
