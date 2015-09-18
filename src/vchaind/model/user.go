package model

type User struct {
    Id            int64  `json:"id"`
    Name          string `json:"name"`
    Email         string `json:"email"`
    EmailVerified bool   `json:"email_verified"`
    Password      string `json:"-"`
    CreateTs      int64  `json:"-"`
    LastLoginTs   int64  `json:"-"`
}
