package sapmodel

import (
	"encoding/json"
)

//ValidateMarshal :
func (r *PinCodeReq) ValidateMarshal() ([]byte, error) {
	a, err := json.Marshal(r)
	return a, err
}
