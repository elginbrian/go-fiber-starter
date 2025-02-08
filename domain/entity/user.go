package domain

import "time"

type User struct {
    ID           string    `json:"id"`          
    Name         string    `json:"name"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`              
    ImageURL     string    `json:"image_url"` 
    Bio          string    `json:"bio"`   
    CreatedAt    time.Time `json:"created"`
    UpdatedAt    time.Time `json:"updated"`
}