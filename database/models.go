package database

type NoteGetAll struct {
        Notes   []NoteGet       `json:"notes"`
}

type NoteGet struct {
        Id      int     `json:"id"`
        UserId  int     `json:"userId"`
        Text    string  `json:"text"`
}

type NoteCreate struct {
        Text    string  `json:"text"`
}
