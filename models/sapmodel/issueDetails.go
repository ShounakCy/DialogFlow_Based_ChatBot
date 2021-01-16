package sapmodel

import (
	"encoding/json"
)

//IssueDetailsMarshal :
func (r *ProductIssueReq) IssueDetailsMarshal() ([]byte, error) {
	a, err := json.Marshal(r)
	return a, err
}
