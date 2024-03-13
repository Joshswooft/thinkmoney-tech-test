package quantity

import "fmt"

// A product quantity is the positive count of the item and therefore cant be initialized below zero
type Quantity struct {
	value int
}

// creates a new non-negative product quantity
func New(qty int) *Quantity {
	if qty < 0 {
		return &Quantity{value: 0}
	}
	return &Quantity{value: qty}
}

func (q *Quantity) Value() int {
	return q.value
}

// Adds a given amount to the quantity
// we ensure that the quantity cant go below zero
func (q *Quantity) Add(amount int) int {

	qty := q.value

	newValue := qty + amount

	if newValue < 0 {
		newValue = 0
	}

	q.value = newValue

	return newValue

}

func (q *Quantity) String() string {
	return fmt.Sprintf("%d", q.value)
}
