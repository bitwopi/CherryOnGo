package sdk

import "time"

type PaymentService interface {
	CreatePayment(CreatePaymentRequest) (*Payment, error)
	GetPayment(paymentID string) (*Payment, error)
	ConfirmPayment(paymentID string) (*Payment, error)
	CancelPayment(paymentID string) (*Payment, error)
	Refund(RefundRequest) (*Refund, error)
}

type Payment struct {
	ID        string
	Status    string
	Amount    Amount
	Income    Amount
	CreatedAt *time.Time
	Provider  string
	Metadata  map[string]interface{}
}

type CreatePaymentRequest struct {
	Amount      *Amount
	Method      string
	Description string
	ReturnURL   string
	Metadata    map[string]interface{}
}

type Amount struct {
	Value    float64
	Currency string
}

type RefundRequest struct {
	PaymentID   string
	Amount      *Amount
	Description string
}

type Refund struct {
	ID        string
	Status    string
	Amount    Amount
	CreatedAt *time.Time
	Provider  string
	Metadata  map[string]interface{}
	PaymentID string
}
