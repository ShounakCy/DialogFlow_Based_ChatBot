package sapmodel

import (
	"encoding/json"
)

//JobDetailsMarshal :
func (r *CreditReq) CreditDetailsMarshal() ([]byte, error) {
	a, err := json.Marshal(r)
	return a, err
}
