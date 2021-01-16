package sap

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/havells/nlp/models"
	sapmodel "github.com/havells/nlp/models/sapmodel"
	repo "github.com/havells/nlp/repo"
	log "github.com/sirupsen/logrus"
)

//AddrService :
type AddrService struct {
	GUID      string
	Phone     string
	DF        models.DfResp
	SessionID string
}

//AddrPinCode :
type AddrPinCode struct {
	PinCode string
}

//AddrAreaSrv :
type AddrAreaSrv struct {
	PinCode string
}

var (

	//VPinURL :
	VPinURL string
	//CAddURL :
	CAddURL string
	// ErrFallBack ...
	ErrFallBack = errors.New("Fallback Intent")
)

//InitValidatPin :
func InitValidatPin(url string) {
	log.Infof("Initializing Validate PinCode ....%v", url)
	VPinURL = url
	log.Infof(" Validate PinCode initialized ....")
	return
}

//InitCreateAddress :
func InitCreateAddress(url string) {
	log.Infof("Initializing Create Address ....%v", url)
	CAddURL = url
	log.Infof(" Create Address initialized ....")
	return
}

//GetFFResp :::::::::::::::::::::::::::::::::::::::::::::::::::;
func (b *AddrService) GetFFResp(ctx context.Context) (interface{}, error) {
	log.Infof("sap/address.go, fetching Brochure info %s,%s,%s", RegCustomerURL)

	data, err := sapmodel.NewAddressAdd(b.Phone, b.SessionID, b.DF)
	log.Infof("sap/address.go, data :", fmt.Sprintf("%v", data))
	if err != nil {
		a := b.DF.Fields["address"].GetStringValue()
		z := b.DF.Fields["zip-code"].GetStringValue()

		switch err {
		case sapmodel.ErrChkReg:
			log.Infof("sap/address.go, Checking registration : %v", b.Phone)
			ok, err := isRegistered(b.Phone, b.SessionID)
			if err != nil {
				return "Dear Customer, Sorry I'm still learning.\n \nIn case I am not able to assist you, please reach us at customercare@havells.com or visit our website at havells.com for further assistance.\n\nThanks !", nil
			}
			if ok {

				log.Infof("sap/address.go, address : %v", a)
				log.Infof("sap/address.go, pincode : %v", z)
				if a != "" && z != "" {
					ok, err := isValidPin(z, b.SessionID)
					if err != nil {
						return "Dear Customer, Sorry I'm still learning.\n \nIn case I am not able to assist you, please reach us at customercare@havells.com or visit our website at havells.com for further assistance.\n\nThanks !", nil
					}
					if ok {
						return "Would you like to set it as Primary Address. Reply with a Yes or No.", sapmodel.ErrPin
					}
					return "Sorry, I'm still learning.\nI would like you to please rephrase so that I can understand it better. Thanks!.", sapmodel.ErrAddressLines
				}
				return "Dear Customer, Please enter your ADDRESS with a VALID PINCODE.\n \nI will be happy to assist you.", sapmodel.ErrAddressLines
			}
			return "Dear Customer,\nKindly register yourself first with NAME and EMAIL ID, before adding a new address.", nil
		case sapmodel.ErrAddressLines:

			log.Infof("Pincode::::::;", z)

			if z != "" {
				ok, err := isValidPin(z, b.SessionID)
				if err != nil {
					return "Dear Customer, Sorry I'm still learning.\n \nIn case I am not able to assist you, please reach us at customercare@havells.com or visit our website at havells.com for further assistance.\n\nThanks !", nil
				}
				if ok {
					return "Would you like to set it as Primary Address. Reply with a Yes or No.", sapmodel.ErrPin
				}
				return "Sorry, we are unable to find valid pincode in address entered by you.\n \nPlease enter 6 digit VALID PINCODE", sapmodel.ErrAddressLines
			}
			return "Sorry, we are unable to find valid pincode in address entered by you.\n \nPlease enter 6 digit VALID PINCODE ", sapmodel.ErrAddressLines

		default:
			return "Dear Customer, I am still learning. In case I am not able to assist you, please reach us at customercare@havells.com !", err
		}
	}

	custGuid := repo.GetCustumerGUID(b.SessionID)
	address := b.DF.Fields["address"].GetStringValue()

	pinGuid := repo.GetPinCodeGUID(b.SessionID)
	isPrim := "2"
	if strings.EqualFold(b.DF.Fields["Confimation"].GetStringValue(), "yes") {
		isPrim = "1"
	}
	return CreateAddress(address, pinGuid, custGuid, isPrim, b.SessionID)
}

func isValidPin(pin, key string) (bool, error) {
	r := sapmodel.PinCodeReq{PinCode: pin}
	res, err := validatePin(&r)
	if err != nil {
		return false, err
	}
	if res.StatusCode == "200" {
		PinGUID := res.PINCodeGUID
		repo.UpdatePinCodeGUID(key, PinGUID)
		return true, nil
	}
	return false, nil
}

func CreateAddress(address, pinGuid, custGuid, isPrim, key string) (string, error) {

	req := sapmodel.AddrReq{AddressLine1: address, PINCodeGUID: pinGuid, CustomerGUID: custGuid, IsPrim: isPrim}
	res, err := addressDetails(&req)

	if err != nil {
		return "Dear Customer, Sorry I'm still learning.\n \nIn case I am not able to assist you, please reach us at customercare@havells.com or visit our website at havells.com for further assistance.\n\nThanks !", err
	}
	if res.StatusCode != "200" {

		return "Dear Customer, Sorry I'm still learning.\nPlease try to register after sometime", nil
	}
	return sendAddrSuccessMsg(key, res.AddressGUID, res.CustomerGUID)
}

func sendAddrSuccessMsg(key, addrguid, custGuid string) (string, error) {
	msgType := repo.GetMasterFlowMsg(key)
	if msgType == nil {
		return "Dear Customer, Your have successfully registered a new address.\n \n You may register a new product or proceed in raising a service request without registering any product  !", ErrFallBack
	}
	switch *msgType {
	case "SRMSG":
		err := repo.UpdateAddressGuid(key, addrguid)
		log.Infof("sap/address.go, updated address guid : %v,%v ", addrguid, err)
		return isProducts(custGuid, key)
	default:
		return "Dear Customer, Your have successfully registered a new address. You may go ahead and raise a service reuest if you are facing any problems with Havells Products.\n \n Thanks!", nil
	}
}
