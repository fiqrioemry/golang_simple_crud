// /pkg/models/job.go
package models

type Job struct {
   	ID			uint	`json:"id" gorm:"primaryKey"`
   	Title		string	`json:"title"`
	Description	string	`json:"description"`
	Company 	string	`json:"company"`
	Salary		string	`json:"salary"`
	Location	string	`json:"location"`
	Category	string	`json:"category"`
}
