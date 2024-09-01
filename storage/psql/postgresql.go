package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Storage struct {
	db *pgxpool.Pool
}

type User struct {
    ID       int    `json:"id"`
    Login    string `json:"login"`
    Password string `json:"password"`
}

type Note struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    AuthorID    int    `json:"author_id"`
}

var (
	ErrUserNotFound = errors.New("url not found")
	ErrUserExists   = errors.New("url exists")
)

func NewStorage(storageParams string) (*Storage, error) {
	const op = "storage.postgresql.NewStorage"
	var err error
	pool, err := pgxpool.Connect(context.Background(), storageParams)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS users(
			id SERIAL PRIMARY KEY,
			login TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			username TEXT NOT NULL);
        CREATE TABLE IF NOT EXISTS notes(
            id SERIAL PRIMARY KEY,
            title TEXT NOT NULL,
            author_id INTEGER REFERENCES users(id),
        	description TEXT NOT NULL);
    `)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return &Storage{db: pool}, nil
}

func (s *Storage) CreateNote(title,description string, user int) (int,error) {
	const op = "storage.psql.CreateNote"

	sql := `INSERT INTO notes (title, description, author_id) VALUES ($1, $2, $3) RETURNING id`

    var noteID int
    err := s.db.QueryRow(context.Background(), sql, title, description, user).Scan(&noteID)
    if err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    return noteID, nil
}

func (s *Storage) GetNotes(user int) ([]Note, error) {
	const op = "storage.psql.GetNotes"

	sql := `SELECT id, title, description, author_id FROM notes WHERE author_id = $1`

    rows, err := s.db.Query(context.Background(), sql, user)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
    defer rows.Close()

    notes := []Note{}

    for rows.Next() {
        var note Note
        err = rows.Scan(&note.ID, &note.Title, &note.Description, &note.AuthorID)
        if err != nil {
            return nil, fmt.Errorf("%s: %w", op, err)
        }
        notes = append(notes, note)
    }

	if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return notes, nil
}

func (s *Storage) SignUp(login,password,username string) (int, error) {
	const op  = "storage.psql.SignUp"

	sql := `INSERT INTO users (login, password, username) VALUES ($1, $2, $3) ON CONFLICT (login) DO NOTHING RETURNING id`
	var userID int

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return 0,fmt.Errorf("%s: %w", op, err)
    }

	err = s.db.QueryRow(context.Background(), sql, login,string(hashedPassword),username).Scan(&userID)
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
        }
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    return userID, nil
}

func (s *Storage) SignIn(login, password string) (*User,error) {
	const op  = "storage.psql.SignIn"

	var user User
    err := s.db.QueryRow(context.Background(), "SELECT id, login, password FROM users WHERE login = $1", login).Scan(&user.ID, &user.Login, &user.Password)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return &user, nil
}

func (s *Storage) Close() error {
    s.db.Close()
    return nil
}