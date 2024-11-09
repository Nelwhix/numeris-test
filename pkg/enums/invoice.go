package enums

import "errors"

type InvoiceState int

const (
	Draft InvoiceState = iota
	Overdue
	Unpaid
	Paid
	PaymentPending
)

func (s InvoiceState) String() string {
	switch s {
	case Draft:
		return "Draft"
	case Overdue:
		return "Overdue"
	case Unpaid:
		return "Unpaid"
	case Paid:
		return "Paid"
	case PaymentPending:
		return "PaymentPending"
	default:
		return "Unknown"
	}
}

func (s InvoiceState) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *InvoiceState) UnmarshalText(text []byte) error {
	switch string(text) {
	case "Draft":
		*s = Draft
	case "Overdue":
		*s = Overdue
	case "Unpaid":
		*s = Unpaid
	case "Paid":
		*s = Paid
	case "PaymentPending":
		*s = PaymentPending
	default:
		return errors.New("invalid InvoiceState value")
	}
	return nil
}

func ParseInvoiceState(value int) (InvoiceState, error) {
	switch value {
	case int(Draft):
		return Draft, nil
	case int(Overdue):
		return Overdue, nil
	case int(Unpaid):
		return Unpaid, nil
	case int(Paid):
		return Paid, nil
	case int(PaymentPending):
		return PaymentPending, nil
	default:
		return -1, errors.New("invalid InvoiceState value")
	}
}
