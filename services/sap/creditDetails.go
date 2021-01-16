package sap

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	//"fmt"
	//"github.com/havells/nlp/models"
	sapmodel "github.com/havells/nlp/models/sapmodel"
	//repo "github.com/havells/nlp/repo"
	log "github.com/sirupsen/logrus"
)

func creditDetails(r *sapmodel.CreditReq) (*sapmodel.CreditResp, error) {
	data, err := r.CreditDetailsMarshal()
	if err != nil {
		log.Errorf("Error marhsaling credit limit request : %v", err)
		return nil, err
	}
	log.Info("payload : ", string(data))
	c := http.Client{Timeout: time.Second * 60}
	req, err := NewHTTPReq("POST", creditURL, creditUser, creditPswd, data)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		log.Errorf("Error fetching Credit Request info : %v", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Infof("Unknown response status code %v", res.StatusCode)
		return nil, errors.New("Unknown status code")
	}
	br := &sapmodel.CreditResp{}
	if err := json.NewDecoder(res.Body).Decode(br); err != nil {
		log.Errorf("Error parsing Credit Request response : %v", err)
		return nil, err
	}
	return br, err
}
