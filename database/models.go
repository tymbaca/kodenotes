package database

import "github.com/google/uuid"

type UserSecureCredentials struct {
        Username string
        Password string
}

type User struct {
        Id       uuid.UUID
        Username string
        Password string
}

type NoteModel struct {
        Id      uuid.UUID
        UserId  uuid.UUID
        Text    string
}

type NoteGetAll struct {
        Notes   []NoteGet       `json:"notes"`
}

type NoteGet struct {
        Id      string  `json:"id"`
        UserId  string  `json:"userId"`
        Text    string  `json:"text"`
}

type NoteCreate struct {
        Text    string  `json:"text"`
}
