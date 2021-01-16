package sap

import (
	"context"
	"fmt"
	"strings"

	"github.com/havells/nlp/models"
	sapmodel "github.com/havells/nlp/models/sapmodel"
	repo "github.com/havells/nlp/repo"
	log "github.com/sirupsen/logrus"
)

//AddrUpdate :
type AddrUpdate struct {
	GUID      string
	Phone     string
	DF        models.DfResp
	SessionID string
}

var (
	//UAddURL :
	UAddURL string
)

//InitUpdateAddress :
func InitUpdateAddress(url string) {
	log.Infof("Initializing Update Address ....%v", url)
	UAddURL = url
	log.Infof(" Create Update initialized ....")
	return
}

//GetFFResp :::::::::::::::::::::::::::::::::::::::::::::::::::;
func (b *AddrUpdate) GetFFResp(ctx context.Context) (interface{}, error) {
	log.Infof("sap/updateAddress.go, fetching Brochure info %s,%s,%s", RegCustomerURL)

	_, err := sapmodel.NewAddressUpdate(b.Phone, b.SessionID, b.DF)

	if err != nil {
		switch err {
		case sapmodel.ErrChkReg:
			log.Infof("sap/updateAddress.go, Checking registration : %v", b.Phone)
			ok, err := isRegistered(b.Phone, b.SessionID)
			if err != nil {
				return "Sorry I'm still learning. Iyris will be able to assist you after sometime.\n \nFor more details please visit our website at havells.com .", nil
			}
			if ok {
				log.Infof("Checking address book", b.SessionID)
				respAddresses, err := isNewAddressBook(b.SessionID)
				if err == sapmodel.ErrFallBackAddress {
					return respAddresses, sapmodel.ErrFallBackProduct
				}
				return "Dear Customer,\nWe don't have your Address. So please add your ADDRESS with a VALID PINCODE first before updation.", nil
			}
			return "Dear Customer,\nKindly register yourself first with NAME and EMAIL ID, before adding a new address.", nil
		case sapmodel.ErrAddressNumber:
			return "Please enter complete address to replace existing address details", sapmodel.ErrAddressLines
		case sapmodel.ErrUpdatedAddress:
			log.Info("Checking Primary")
			return "Would you like to set this address as permanent ? Reply with yes or no", sapmodel.ErrAddressLines
		default:
			return "Dear Customer, I am still learning. In case I am not able to assist you, please reach us at customercare@havells.com !", err
		}
	}
	addressNum := b.DF.Fields["addressNumber"].GetNumberValue()
	log.Infof("sap/updateAddress.go, Updated Address : %v", addressNum)
	pinGuid := repo.GetAddressPinCodeGUID(b.SessionID, int(addressNum-1))
	addrGuid := repo.GetAddressGUID(b.SessionID, int(addressNum-1))
	addressUpdated := b.DF.Fields["updatedAddress"].GetStringValue()
	log.Infof("sap/updateAddress.go,, isPRIM : %v", b.DF.Fields["Confirmation"].GetStringValue())
	custGuid := repo.GetCustumerGUID(b.SessionID)

	isPrim := "2"
	if strings.EqualFold(b.DF.Fields["Confimation"].GetStringValue(), "yes") {
		isPrim = "1"
	}
	resp, _ := UpdateAddress(addressUpdated, addrGuid, pinGuid, custGuid, isPrim)
	return resp, nil

}

func isNewAddressBook(key string) (string, error) {

	custGUID := repo.GetCustumerGUID(key)
	log.Infof("sap/updateAddress.go, Forming address request: %v", custGUID)
	reg := sapmodel.AddressReq{CustomerGUID: custGUID}
	res, err := addressBook(&reg)

	log.Infof("sap/updateAddress.go, Address Response", fmt.Sprintf("%v", res))
	if err != nil {
		return "", err
	}

	var ans []string
	var addresses []string
	var fulladdress []string
	var pinGuids []string
	for i, v := range res {
		if v.StatusCode == "200" {
			FullAddresses := fmt.Sprintf("%v", v.FullAddress)
			resp := fmt.Sprintf("%d", i+1) + ". " + FullAddresses + "\n\n"
			ans = append(ans, resp)
			addresses = append(addresses, fmt.Sprintf("%v", v.AddressGUID))
			fulladdress = append(fulladdress, fmt.Sprintf("%v", v.FullAddress))
			pinGuids = append(pinGuids, fmt.Sprintf("%v", v.PINCodeGUID))
		}
		if v.StatusCode == "204" {
			return "We dont have your Address. So please add your address first and then Update if required.", nil
		}
	}
	finalResp := "Dear Customer, Below is the list of your registered addresses:\n\n" + strings.Join(ans, "") + "\nPlease reply with respective option number to update an address"
	go repo.UpdateAddressMap(key, addresses)
	go repo.UpdateFullAddressMap(key, fulladdress)
	go repo.UpdateAddressPinCodeGUID(key, pinGuids)
	return finalResp, sapmodel.ErrFallBackAddress
}

func UpdateAddress(addressUpdated, addrGuid, pinGuid, custGuid, isPrim string) (string, error) {

	req := sapmodel.AddrReq{AddressLine1: addressUpdated, AddrGUID: addrGuid, PINCodeGUID: pinGuid, CustomerGUID: custGuid, IsPrim: isPrim}
	res, err := updateaddressDetails(&req)

	if err != nil {
		return "Dear Customer, I'm not able to find this particular address.\n \nIn case I am not able to assist you, please reach us at customercare@havells.com for further assistance.\n \nThanks !", err
	}
	if res.StatusCode != "200" {

		return "Dear Customer, Sorry I'm still learning.\nPlease try to update after sometime", nil
	}
	return "Dear Customer, Your address is updated successfully.\n \nYou may go ahead and add a new product or raise a service request if you are facing any problems with Havells Products.\n \nThanks!", nil
}
