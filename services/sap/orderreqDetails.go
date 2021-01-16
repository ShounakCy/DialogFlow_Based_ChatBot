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

func orderreqDetails(r *sapmodel.OrderReq) (*sapmodel.OrderResp, error) {
	data, err := r.OrderreqDetailsMarshal()
	if err != nil {
		log.Errorf("Error marhsaling order request : %v", err)
		return nil, err
	}
	log.Info("payload : ", string(data))
	c := http.Client{Timeout: time.Second * 60}
	req, err := NewHTTPReq("POST", OrderRequestURL, OrderRequestUser, OrderRequestPass, data)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		log.Errorf("Error fetching OrderRequest info : %v", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Infof("Unknown response status code %v", res.StatusCode)
		return nil, errors.New("Unknown status code")
	}
	br := &sapmodel.OrderResp{}
	if err := json.NewDecoder(res.Body).Decode(br); err != nil {
		log.Errorf("Error parsing orderRequest response : %v", err)
		return nil, err
	}
	return br, err
}
