package sapmodel

import (
	"encoding/json"
)

//JobDetailsMarshal :
func (r *OrderReq) OrderreqDetailsMarshal() ([]byte, error) {
	a, err := json.Marshal(r)
	return a, err
}
