package storage

type Note struct {
        Id      int     `json:"id"`
        UserId  int     `json:"userId"`
        Text    string  `json:"text"`
}
