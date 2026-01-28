package sdk

import (
	"fmt"
	"strconv"

	"github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
	yoorefund "github.com/rvinnie/yookassa-sdk-go/yookassa/refund"
)

type Yookassa struct {
	PaymentService
	client *yookassa.Client
}

func NewYookassa(shopID string, secretKey string) *Yookassa {
	return &Yookassa{
		client: yookassa.NewClient(shopID, secretKey),
	}
}

func (y *Yookassa) CreatePayment(req CreatePaymentRequest) (*Payment, error) {
	handler := yookassa.NewPaymentHandler(y.client)

	payment, err := handler.CreatePayment(
		&yoopayment.Payment{
			Amount: &yoocommon.Amount{
				Value:    fmt.Sprintf("%.2f", req.Amount.Value),
				Currency: "RUB",
			},
			PaymentMethod: yoopayment.PaymentMethodType(req.Method),
			Confirmation: yoopayment.Redirect{
				Type:      "redirect",
				ReturnURL: req.ReturnURL,
			},
			Description: req.Description,
			Metadata:    req.Metadata,
		})
	if err != nil {
		return nil, err
	}

	return convertPayment(payment)
}
func (y *Yookassa) ConfirmPayment(paymentID string) (*Payment, error) {
	handler := yookassa.NewPaymentHandler(y.client)

	payment, err := handler.FindPayment(paymentID)
	if err != nil {
		return nil, err
	}

	updPayment, err := handler.CapturePayment(payment)
	if err != nil {
		return nil, err
	}

	return convertPayment(updPayment)
}

func (y *Yookassa) CancelPayment(paymentID string) (*Payment, error) {
	handler := yookassa.NewPaymentHandler(y.client)

	updPayment, err := handler.CancelPayment(paymentID)
	if err != nil {
		return nil, err
	}

	return convertPayment(updPayment)
}

func (y *Yookassa) GetPayment(paymentID string) (*Payment, error) {
	handler := yookassa.NewPaymentHandler(y.client)

	payment, err := handler.FindPayment(paymentID)
	if err != nil {
		return nil, err
	}

	return convertPayment(payment)
}

func (y *Yookassa) Refund(req RefundRequest) (*Refund, error) {
	handler := yookassa.NewRefundHandler(y.client)

	refund, err := handler.CreateRefund(&yoorefund.Refund{
		PaymentId: req.PaymentID,
		Amount: &yoocommon.Amount{
			Value:    fmt.Sprintf("%.2f", req.Amount.Value),
			Currency: "RUB",
		},
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	value, err := strconv.ParseFloat(refund.Amount.Value, 64)
	if err != nil {
		return nil, err
	}

	return &Refund{
		ID:        refund.Id,
		PaymentID: refund.PaymentId,
		Amount: Amount{
			Value:    value,
			Currency: refund.Amount.Currency,
		},
		Status:    string(refund.Status),
		CreatedAt: refund.CreatedAt,
		Provider:  "YooKassa",
	}, nil
}

func convertPayment(payment *yoopayment.Payment) (*Payment, error) {
	value, err := strconv.ParseFloat(payment.Amount.Value, 64)
	if err != nil {
		return nil, err
	}
	incValue, err := strconv.ParseFloat(payment.IncomeAmount.Value, 64)
	if err != nil {
		return nil, err
	}
	md, ok := payment.Metadata.(map[string]interface{})
	if !ok {
		md = nil
	}
	return &Payment{
		ID: payment.ID,
		Amount: Amount{
			Value:    value,
			Currency: payment.Amount.Currency,
		},
		Income: Amount{
			Value:    incValue,
			Currency: payment.IncomeAmount.Currency,
		},
		Status:    string(payment.Status),
		CreatedAt: payment.CreatedAt,
		Provider:  "YooKassa",
		Metadata:  md,
	}, nil
}
