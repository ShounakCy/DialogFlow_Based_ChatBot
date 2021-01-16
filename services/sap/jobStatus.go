package sap

import (
	"context"
	"strings"

	"github.com/havells/nlp/models"
	sapmodel "github.com/havells/nlp/models/sapmodel"
	log "github.com/sirupsen/logrus"
)

type StatusSrv struct {
	Phone     string
	SessionID string
	DF        models.DfResp
}

var (
	// ServiceStatusURL :
	ServiceStatusURL string
)

// InitServiceStatus :
func InitServiceStatus(url string) {
	log.Infof("Initializing Service Status ....%v", url)
	ServiceStatusURL = url
	log.Infof(" Create Service Status initialized ....")
	return
}

func (b *StatusSrv) GetFFResp(ctx context.Context) (interface{}, error) {

	_, err := sapmodel.NewServiceStatus(b.Phone, b.SessionID, b.DF)
	if err != nil {
		switch err {
		case sapmodel.ErrChkStatus:
			ID := b.DF.Fields["num"].GetStringValue()
			log.Infof("JobId::::::::::::", ID)
			if ID != "" {
				sSresp, _ := ServiceStatus(ID)
				return sSresp, nil
			}
			return "Please provide JOB ID to know status", sapmodel.ErrJobID
		default:
			return "Dear Customer, I am still learning. In case I am not able to assist you, please reach us at customercare@havells.com !", err
		}
	}
	JobID := b.DF.Fields["num"].GetStringValue()
	sSresp, err := ServiceStatus(JobID)
	return sSresp, nil
}

func ServiceStatus(Jobid string) (string, error) {

	req := sapmodel.JobStatusReq{ID: Jobid}
	res, err := jobDetails(&req)
	if err != nil {
		return "", nil
	}

	for _, j := range res {
		resp := `Customer NAME  : ` + j.CustomerName + `

	MOBILE NO  : ` + j.MobileNumber + `

	ADDRESS  : ` + strings.ReplaceAll(j.CustomerAddress, " ,", "") + `

	SERIAL NO  :  ` + j.Asset + `

	PRODUCT  :  ` + j.Product + `

	DESCRIPTION  :  ` + j.ChiefComplaint + `

	JOB ID  : ` + j.ID + `

	JOB STATUS  : ` + j.Status + `

	CREATED ON  : ` + j.Loggedon + `

	CALL TYPE   : ` + j.CallSubType + `

	ASSIGNED TO  : ` + j.AssignedTo + "\n\nPlease find the above Job Status. Thanks !"
		if resp != "" {
			return resp, sapmodel.ErrEndSessionService

		}
		if j.ClosedOn != "" {
			resp = resp + "\n\n*CLOSED ON* : " + j.ClosedOn
			return resp, nil
		}
	}
	return "Dear Customer ! Please check the Job ID entered.\nPlease EXIT if you want to start over again !", nil
}
