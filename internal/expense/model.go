package expense

import "time"

type Expense struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"createdAt"`
}
