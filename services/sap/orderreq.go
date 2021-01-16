package sap

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/havells/nlp/models"
	"github.com/havells/nlp/models/sapmodel"
	log "github.com/sirupsen/logrus"
)

type ORservice struct {
	Code      string
	SessionID string
	DF        models.DfResp
}

var (
	// OrderRequestURL :
	OrderRequestURL  string
	OrderRequestUser string
	OrderRequestPass string
)

// InitOrderRequest :
func InitOrderRequest(url, user, pswd string) {
	log.Infof("Initializing Order Request ....%v", url)
	OrderRequestURL = url
	OrderRequestUser = user
	OrderRequestPass = pswd
	log.Infof(" Create Order Request initialized ....")
	return
}

func (b *ORservice) GetFFResp(ctx context.Context) (interface{}, error) {

	_, err := sapmodel.NewOrderRequest(b.Code, b.SessionID, b.DF)

	if err != nil {
		inputs := CreateOrderInputs(b.DF)
		if inputs == nil {
			return "HI Dealer, Please enter a valid SKU Code (in capitals), and then place the order with your desired Quantity ! ", err
		}

		switch err {
		case sapmodel.ErrOrderRequest:

			log.Infof("inputs...%v, length of input...%v", inputs, len(inputs))
			if inputs != nil {
				log.Infof("Creating the order request")
				ORresp, _ := makeOrder(b.Code, b.SessionID, inputs)
				log.Infof("ORresp...%v", ORresp)
				log.Infof("err......%v", err)
				return ORresp, sapmodel.ErrOrderConfirmCancel
			}
			return "HI Dealer, Please provide the SKU code (in capitals) and Quantity to place the order !", err
		case sapmodel.ErrOrderConfirm:
			log.Infof("inputs...%v", inputs)
			log.Infof("Confirming the order")
			if inputs != nil {
				ORconfirm, _ := createOrder(b.Code, b.SessionID, inputs)
				log.Infof("ORresp...%v", ORconfirm)
				log.Infof("err......%v", err)
				return ORconfirm, sapmodel.ErrEndSessionOrder
			}
		case sapmodel.ErrOrderCancel:
			log.Infof("inputs...%v", inputs)
			log.Infof("Cancelling the order")
			if inputs != nil {
				return "Hi Dealer, You have opted not to place any order ! Thank you !", sapmodel.ErrEndSessionOrder
			}
			return "HI Dealer, Please enter a valid SKU Code (in capitals), and then place the order with your desired Quantity ! ", err
		default:
			return "Dear Customer, I am still learning. In case I am not able to assist you, please reach us at customercare@havells.com !", sapmodel.ErrEndSessionService
		}
	}
	// inputs := CreateOrderInputs(b.DF)
	// ORresp, err := ORresp(b.Code, b.SessionID, inputs)
	return "", nil
}
func createOrder(code, session string, inputs []sapmodel.OrderInput) (string, error) {

	var txt string
	if len(inputs) == 0 {
		txt = "Incorrect order is placed. Please try again after sometime !"
		return txt, nil
	}
	os := sapmodel.OrderReq{IMKUNNR: "CSA1056", IMORDTYP: "ZNSO", IMINPUTLINES: inputs, IMCONFIRM: "X"}
	r, err := orderreqDetails(&os)
	if err != nil || r == nil {
		txt = "Order not processed yet ! Please try again after sometime !"
		return txt, err
	}
	if r.EXSALESDOCNO == "" {
		txt = "HI Dealer, Please enter a valid SKU Code (in capitals), and then place the order with your desired Quantity ! "
		return txt, nil
	}
	log.Infof("response from order request api %v", r)
	txt = MakeOrderConfirmed(code, session, r.EXSALESDOCNO, r.ETSODETAILS)
	return txt, nil

}

//MakeOrderConfirmed :
func MakeOrderConfirmed(code, session, orderID string, in []sapmodel.ETSODETAILS) string {
	var txt = fmt.Sprintf("Hi Dealer " + "\n\nYour Order ID is :" + " (" + orderID + ")" + "\n")
	for k, v := range in {

		txt = txt + ` ` + "\n\n" +
			strconv.Itoa(k+1) + `.` + "\n" +
			`Material Code. : ` + v.MATNR + "\n" +
			`Material Description : ` + v.MAKTX + "\n" +
			`Quantity : ` + strings.TrimSpace(v.KWMENG) + "\n" +
			`Per Item Price : ₹` + strings.TrimSpace(v.NETPR) + "\n" +
			`Order Price : ₹` + strings.TrimSpace(v.NETWR) + "\n" +
			`Tax : ₹` + strings.TrimSpace(v.MWSBP) + "\n" +
			`Total Price : ₹` + getTotalPrice(v.NETWR, v.MWSBP)
	}
	txt = txt + "\n\nThank you for placing your Order Request. "
	return txt
}

func makeOrder(code, session string, inputs []sapmodel.OrderInput) (string, error) {

	var txt string
	if len(inputs) == 0 {
		txt = "Incorrect Order"
		return txt, nil
	}
	os := sapmodel.OrderReq{IMKUNNR: "CSA1056", IMORDTYP: "ZNSO", IMINPUTLINES: inputs}
	r, err := orderreqDetails(&os)
	if err != nil || r == nil {
		//	log.Infof("Order confirm error : %v, Message: %v", err, r)
		txt = "Order not processed yet ! Please try again after sometime !"
		return txt, err
	}
	log.Infof("Order confirm error : %v, Message: %v", err, r)
	if r.EXSALESDOCNO == "" {
		txt = "HI Dealer, Please enter a valid SKU Code (in capitals), and then place the order with your desired Quantity !"
		return txt, nil
	}
	log.Infof("Response from order request api  %v", r)
	txt = MakeOrderReqCofirmPrompt(code, session, r.ETSODETAILS)
	return txt, sapmodel.ErrOrderConfirm

}

//MakeOrderReqCofirmPrompt :
func MakeOrderReqCofirmPrompt(code, session string, resp []sapmodel.ETSODETAILS) string {
	var txt = fmt.Sprintf("Hi Dealer " + ", You have requested for below items :\n")
	var itemList []sapmodel.OrderInput

	for k, v := range resp {
		txt = txt + ` ` + "\n\n" +
			strconv.Itoa(k+1) + `.` + "\n" +
			`Material No. : ` + v.MATNR + "\n" +
			`Quantity : ` + strings.TrimSpace(v.KWMENG) + "\n" +
			`Order Price : ₹` + strings.TrimSpace(v.NETWR) + "\n" +
			`Tax : ₹` + strings.TrimSpace(v.MWSBP) + "\n" +
			`Total Price : ₹` + getTotalPrice(v.NETWR, v.MWSBP)
		itemList = append(itemList, sapmodel.OrderInput{MATNR: v.MATNR, KWMENG: v.KWMENG})
	}

	txt = txt + "\n\nPlease check and confirm the above details by typing yes !"
	return txt
}

//CreateOrderInputs :
func CreateOrderInputs(df models.DfResp) []sapmodel.OrderInput {

	material := df.Fields["skucode"].GetListValue()

	log.Infof("materials.... %v", material)
	quantity := df.Fields["quantity"].GetListValue()
	log.Infof("quantities.... %v", quantity)

	if len(material.GetValues()) != len(quantity.GetValues()) {
		return nil
	}
	var resp []sapmodel.OrderInput
	for i, v := range material.GetValues() {
		log.Infof("material(v)...%v", v.GetStringValue())

		q := quantity.GetValues()[i]
		log.Infof("q........", q)
		log.Infof("quantity(q)...%v", q.GetNumberValue())

		if t := getInput(v.GetStringValue(), q.GetNumberValue()); t != nil {
			resp = append(resp, *t)
			log.Infof("resp...%v", resp)
		}
	}
	log.Infof("resp...%v", resp)
	return resp

}

func getInput(material string, quantity float64) *sapmodel.OrderInput {

	quant := strconv.FormatFloat(quantity, 'f', 0, 64)
	if _, err := strconv.Atoi(quant); err != nil {
		log.Errorf("Not a valid quantity %v ", quant)
		return nil
	}
	log.Infof("quant.....", quant)
	return &sapmodel.OrderInput{MATNR: material, KWMENG: quant}
}

// func updateOrderMap(userCode string,order *OrderConfirmDetail){
// 	if order.Done {
// 		OrderReqMap.Delete(userCode)
// 		return
// 	}
// 	OrderReqMap.Store(userCode,*order)
// 	return
// }

//GetConfirmedOrders :
// func GetConfirmedOrders(userCode string)[]sapmodel.OrderInput{
// 	tmp := []sapmodel.OrderInput{}
// 	if d, ok := OrderReqMap.Load(userCode); ok {
// 		if orders, ok := d.(OrderConfirmDetail); ok {
// 			tmp = orders.Items
// 		}
// 	}
// 	return tmp
// }

func getTotalPrice(net, tax string) string {
	n, _ := strconv.ParseFloat(strings.TrimSpace(net), 64)
	t, _ := strconv.ParseFloat(strings.TrimSpace(tax), 64)
	return fmt.Sprintf("%.3f", n+t)
}
