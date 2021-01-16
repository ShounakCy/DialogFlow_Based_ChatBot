package sapmodel

import (
	"encoding/json"
)

//WalletDetailsMarshal :
func (r *WalletReq) WalletDetailsMarshal() ([]byte, error) {
	a, err := json.Marshal(r)
	return a, err
}
