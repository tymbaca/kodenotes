package storage

// PostgresStorage implements Storage interface
type PostgresStorage struct {
        
}

func (s *PostgresStorage) GetNotes() []Note {}

func (s *PostgresStorage) PostNote(Note) {}
