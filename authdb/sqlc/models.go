// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package authdb

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Passwd    string    `json:"passwd"`
	CreatedAt time.Time `json:"created_at"`
}
