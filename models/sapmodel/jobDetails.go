package sapmodel

import (
	"encoding/json"
)

//JobDetailsMarshal :
func (r *JobStatusReq) JobDetailsMarshal() ([]byte, error) {
	a, err := json.Marshal(r)
	return a, err
}
