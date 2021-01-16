package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/havells/nlp/models"
	"github.com/havells/nlp/models/sapmodel"
	"github.com/havells/nlp/services"
	"github.com/havells/nlp/services/google"
	log "github.com/sirupsen/logrus"
)

//App :
type App struct {
	Service services.Service //services.go -> Service interface
}

//VoiceReq : request handler for incoming nlp request
func (a *App) VoiceReq(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*60)
	defer cancel()
	msg := models.AppReq{}
	if err := c.BindJSON(&msg); err != nil {
		log.Errorf("app.go, ERROR UNMARSHALLING APP REQUEST: %v", err)
		c.JSON(http.StatusBadRequest, "Bad request")
		return
	}

	log.Infof("app.go,incoming nlp request: %v", msg)
	enTxt := translate(ctx, msg.Msg, msg.LangCode, "en-IN", a.Service.GetTranslateService())
	// get intent responce from dialogflow service
	dfResp, err := a.Service.GetDFService().GetIntent(c, msg.SessionID, enTxt)
	if err != nil || dfResp.Intent == "" {
		tmp := "Sorry I'm still learning.\nI would like you to PLEASE REPHRASE so that I can understand it better.\nPlease type EXIT if you want to start over again ! "
		c.JSON(http.StatusOK, &models.AppResp{Data: translate(ctx, tmp, msg.LangCode, msg.LangCode,
			a.Service.GetTranslateService()), SessionID: msg.SessionID})
		return
	}

	srv := a.Service.GetFFService(getQueryParams("Param1", *dfResp), msg.MobileNumber, msg.SessionID, msg.UserCode, *dfResp)
	log.Infof("app.go, srv: %v", srv)
	if srv == nil && dfResp.Intent == "Fallback" {
		tmp := "Iyris is still learning !\n \nI would like you to PLEASE REPHRASE!\n \nI can assist you with placing an order request.\n \nPlease let me know how can I help you !"
		c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
			a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
		return
	}
	if srv == nil && dfResp.Intent == "EndIntent" {
		tmp := "Iyris is happy to help you !\n \nIt is my pleasure to have you on this chat today!\n \nI can assist you with placing an order request.\n \nPlease let me know how can I help you\n \nThank you and have a nice day!"
		c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
			a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
		return
	}
	if srv == nil && dfResp.Intent == "HowAreYouIntent" {
		tmp := "Hi, I am doing good. !\n \nIt is my pleasure to have you on this chat today!\n \nI can assist you with placing an order request.\n \nPlease let me know how can I help you."
		c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
			a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
		return
	}
	if srv == nil {
		log.Infof("could not resolve ff service : %v", dfResp)
		tmp := "Dear Customer, Sorry I'm still learning.\n \nIn case I am not able to assist you, please reach us at customercare@havells.com for further assistance.\n\nThanks !"
		c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
			a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
		return
	}

	// get fulfillment service response
	r, err := srv.GetFFResp(ctx)
	log.Infof("app.go, r: %v ", r)
	log.Infof("app.go, err %v:", err)

	if err != nil {
		// if err == sapmodel.ErrFallBack {
		// 	c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
		// 		msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
		// 	return
		// }
		if r == "" {
			tmp := "Dear User, I am facing problem in fetching the desired information currently. Please try after sometime or reach us at customercare@havells.com for further assistance.\n\nThanks !"
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
				a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
			return
		}
		if err == sapmodel.ErrCustomerName {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		if err == sapmodel.ErrSerial {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		if err == sapmodel.ErrCustomerEmail {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		if err == sapmodel.ErrFallBackProduct {
			if r == "I'm sorry but you have no product registered !\nPlease type EXIT if you want to start over again ! " {
				c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
					msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
				return
			}
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		if err == sapmodel.ErrFallBackProductType {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		if err == sapmodel.ErrFallBackRegistration {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}

		if err == sapmodel.ErrFallBackIssue {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		if err == sapmodel.ErrFallBackComplaint {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		if err == sapmodel.ErrFallBackAddress {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		if err == sapmodel.ErrFallBackService {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		if err == sapmodel.ErrJobID {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		if err == sapmodel.ErrAddressLines {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		if err == sapmodel.ErrPin {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		if err == sapmodel.ErrOrderRequest {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		if err == sapmodel.ErrCreditRequest {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		if err == sapmodel.ErrWalletRequest {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		if err == sapmodel.ErrOrderConfirmCancel {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
				msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		// if err == sapmodel.ErrOrde{
		// 	c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode,
		// 		msg.LangCode, a.Service.GetTranslateService()), SessionID: msg.SessionID})
		// 	return
		// }
		// if err == sapmodel.ErrEndSession {
		// 	tmp := "You have opted not to apply for any leave, Thanks !!!"
		// 	c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
		// 		a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
		// 	return
		// }
		if err == sapmodel.ErrEndSessionCustomer {
			tmp := "You have opted not to register yourself, Thanks !!!"
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
				a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
			return
		}
		if err == sapmodel.ErrEndSessionProduct {
			tmp := "You have opted not to provide any Product Serial Code, Thanks !!!"
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
				a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
			return
		}
		if err == sapmodel.ErrEndSessionCreateSR {
			tmp := "You have opted not to raise any Service Request, Thanks !!!"
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
				a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
			return
		}
		if err == sapmodel.ErrWrongAddress {
			tmp := "Dear Customer, We are not able to raise a service request for the address option selected. Please try to raise the service request again !"
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
				a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
			return
		}
		if err == sapmodel.ErrEndSessionAddressAdd {
			tmp := "You have opted not to add any new address, Thanks !!!"
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
				a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
			return
		}
		if err == sapmodel.ErrEndSessionAddressUpdate {
			tmp := "You have opted not to update any new address, Thanks !!!"
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
				a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
			return
		}
		if err == sapmodel.ErrEndJobStatus {
			tmp := "You have opted out of Job Status, Please visit again. Thanks !!!"
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
				a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
			return
		}
		if err == sapmodel.ErrEndRephrase {
			tmp := "Sorry I'm still learning.\nI would like you to PLEASE REPHRASE and provide all the information !\nPlease type EXIT if you want to start over again ! "
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
				a.Service.GetTranslateService()), SessionID: msg.SessionID})
			return
		}
		if err == sapmodel.ErrEndSessionService {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode, msg.LangCode,
				a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
			return
		}
		if err == sapmodel.ErrEndSessionOrder {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode, msg.LangCode,
				a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
			return
		}
		if err == sapmodel.ErrEndSessionCredit {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode, msg.LangCode,
				a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
			return
		}
		if err == sapmodel.ErrEndSessionWallet {
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode, msg.LangCode,
				a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
			return
		}
		if err == sapmodel.ErrFallBackIndex {
			tmp := "You've entered INVALID INDEX NUMBER. Please try to raise the Service Request again."
			c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
				a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
			return
		}

		c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r), msg.LangCode, msg.LangCode,
			a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
		return
	}
	if r == nil || r == "" {
		tmp := "Dear Customer, Sorry I'm still learning.\n\nIn case I am not able to assist you, please reach us at customercare@havells.com for further assistance.\n\nThanks !"
		c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
			a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
		return
	}
	if r == "" {
		tmp := "Dear Customer, Sorry I'm still learning.\n\nIn case I am not able to assist you, please reach us at Customercare@havells.com for further assistance.\n\nThanks !"
		c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, tmp, msg.LangCode, msg.LangCode,
			a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
		return
	}

	c.JSON(http.StatusOK, &models.AppResp{Data: translate(c, fmt.Sprintf("%v", r),
		msg.LangCode, msg.LangCode, a.Service.GetTranslateService()), SessionID: getSessionID(msg.UserCode)})
	return
}

func translate(c context.Context, msg, langCode, toLangCode string, ts google.TranslateService) string {
	if langCode == "en-IN" {
		return msg
	}
	r, err := ts.Translate(c, toLangCode, msg)
	if err != nil {
		return msg
	}
	return r
}

func getQueryParams(params string, dfResp models.DfResp) string {
	return dfResp.Fields[params].GetStringValue()
}

func getSessionID(userID string) string {
	return userID + fmt.Sprintf("%d", time.Now().UnixNano())
}
