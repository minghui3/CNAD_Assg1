package models

type User struct {
    ID           int       `json:"id"`
    Name         string    `json:"name"`
    Email        string    `json:"email"`
    PhoneNumber  string    `json:"phone_number"`
    Membership   string    `json:"membership_tier"`
    PasswordHash string    `json:"-"` // Excluded from JSON for security
}
