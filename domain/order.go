package domain

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyOrder       = errors.New("cannot pay empty order")
	ErrAlreadyPaid      = errors.New("order is already paid")
	ErrCannotModifyPaid = errors.New("cannot modify paid order")
	ErrInvalidTotal     = errors.New("total does not match sum of lines")
)

type OrderID string

type OrderStatus string

const (
	OrderStatusPending OrderStatus = "pending"
	OrderStatusPaid    OrderStatus = "paid"
)

type Money struct {
	amount   int64
	currency string
}

func NewMoney(amount int64, currency string) Money {
	return Money{amount: amount, currency: currency}
}

func (m Money) Amount() int64 {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency == "" && m.amount == 0 {
		return NewMoney(other.amount, other.currency), nil
	}

	if other.currency == "" && other.amount == 0 {
		return NewMoney(m.amount, m.currency), nil
	}

	if m.currency != other.currency {
		return Money{}, fmt.Errorf("currency mismatch: %s != %s", m.currency, other.currency)
	}
	return NewMoney(m.amount+other.amount, m.currency), nil
}

func (m Money) String() string {
	return fmt.Sprintf("%d %s", m.amount/100, m.currency)
}

type OrderLine struct {
	ProductID OrderID
	Quantity  int
	Price     Money
}

func (l OrderLine) Total() Money {
	return NewMoney(l.Price.amount*int64(l.Quantity), l.Price.currency)
}

type Order struct {
	id     OrderID
	status OrderStatus
	lines  []OrderLine
	total  Money
}

func NewOrder(id OrderID, lines []OrderLine) (*Order, error) {
	if len(lines) == 0 {
		return nil, ErrEmptyOrder
	}

	total := lines[0].Total()
	for i := 1; i < len(lines); i++ {
		var err error
		total, err = total.Add(lines[i].Total())
		if err != nil {
			return nil, err
		}
	}

	order := &Order{
		id:     id,
		status: OrderStatusPending,
		lines:  lines,
		total:  total,
	}

	if err := order.ValidateTotal(); err != nil {
		return nil, err
	}

	return order, nil
}

func (o *Order) ID() OrderID {
	return o.id
}

func (o *Order) Status() OrderStatus {
	return o.status
}

func (o *Order) Lines() []OrderLine {
	lines := make([]OrderLine, len(o.lines))
	copy(lines, o.lines)
	return lines
}

func (o *Order) Total() Money {
	return o.total
}

func (o *Order) Pay() error {
	if len(o.lines) == 0 {
		return ErrEmptyOrder
	}
	if o.status == OrderStatusPaid {
		return ErrAlreadyPaid
	}
	o.status = OrderStatusPaid
	return nil
}

func (o *Order) AddLine(line OrderLine) error {
	if o.status == OrderStatusPaid {
		return ErrCannotModifyPaid
	}
	o.lines = append(o.lines, line)
	lineTotal, err := o.total.Add(line.Total())
	if err != nil {
		return err
	}
	o.total = lineTotal
	return nil
}

func (o *Order) ValidateTotal() error {
	var sum Money
	for _, line := range o.lines {
		lineTotal, err := sum.Add(line.Total())
		if err != nil {
			return err
		}
		sum = lineTotal
	}
	if sum.amount != o.total.amount {
		return ErrInvalidTotal
	}
	return nil
}
