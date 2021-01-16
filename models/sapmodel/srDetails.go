package sapmodel

import (
	"encoding/json"
)

//IssueDetailsMarshal :
func (r *SRrequest) SrDetailsMarshal() ([]byte, error) {
	a, err := json.Marshal(r)
	return a, err
}
