package model

import (
    "github.com/satori/go.uuid"
)

func generateHash() string {
    return uuid.NewV4().String()
}
