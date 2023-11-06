package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/quible-backend/auth-service/domain"
)

const (
	passwordHashCost = 15
)

type Database struct {
	DB *pgxpool.Pool
}

func NewRepository(pgx *pgxpool.Pool) Database {
	return Database{DB: pgx}
}

func (pool Database) Create(user domain.User) (int64, error) {
	q := `
  INSERT INTO users (username, email, hashed_password, full_name, phone, image, is_oauth)
  VALUES ($1,$2,$3,$4,$5,$6,$7)
  RETURNING id;
  `

	row := pool.DB.QueryRow(
		context.Background(),
		q,
		user.Username,
		user.Email,
		user.HashedPassword,
		user.FullName,
		user.Phone,
		user.Image,
		user.IsOauth,
	)

	u := new(domain.User)

	err := row.Scan(&u.ID)
	if err != nil {
		return 0, err
	}

	return u.ID, nil
}

func (pool Database) Gets(id int64) (*domain.UserResponse, error) {
	query := `
  SELECT id,username,email,full_name,phone,is_oauth 
  FROM users 
  WHERE id = $1
  `

	row := pool.DB.QueryRow(
		context.Background(),
		query,
		id,
	)

	u := new(domain.UserResponse)

	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.FullName,
		&u.Phone,
		&u.IsOauth,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (pool Database) GetUserByEmail(email string) (*domain.UserResponse, error) {
	query := `
  SELECT id,username,email,full_name,phone,is_oauth,created_at,updated_at
  FROM users 
  WHERE email = $1
  `

	row := pool.DB.QueryRow(
		context.Background(),
		query,
		email,
	)

	u := new(domain.UserResponse)

	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.FullName,
		&u.Phone,
		&u.IsOauth,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (pool Database) GetByUsername(username string) (*domain.UserResponse, error) {
	query := `
  SELECT id,username,email,full_name,phone,is_oauth,created_at,updated_at
  FROM users 
  WHERE username = $1
  `

	row := pool.DB.QueryRow(
		context.Background(),
		query,
		username,
	)

	u := new(domain.UserResponse)

	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.FullName,
		&u.Phone,
		&u.IsOauth,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (pool Database) GetLoginCredential(email string) (*domain.UserLoginResponse, error) {
	query := `
  SELECT id,username,email,hashed_password
  FROM users 
  WHERE email = $1 or username = $1
  `

	row := pool.DB.QueryRow(
		context.Background(),
		query,
		email,
	)

	u := new(domain.UserLoginResponse)

	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.HashedPassword,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (pool Database) Update(user domain.User) (int64, error) {
	q := `
  UPDATE users SET
  username=$2,email=$3,full_name=$4,phone=$5,image=$6,is_oauth=$7,updated_at=$8
  WHERE id = $1
  RETURNING id;
  `
	row := pool.DB.QueryRow(
		context.Background(),
		q,
		user.ID,
		user.Username,
		user.Email,
		user.FullName,
		user.Phone,
		user.Image,
		user.IsOauth,
		user.UpdatedAt,
	)

	u := new(domain.User)

	if err := row.Scan(&u.ID); err != nil {
		return 0, err
	}

	return u.ID, nil
}

func (pool Database) Delete(id int64) (int64, error) {
	q := `
  DELETE FROM users
  WHERE id = $1
  RETURNING id;
  `

	row := pool.DB.QueryRow(context.Background(), q, id)

	u := new(domain.User)

	if err := row.Scan(&u.ID); err != nil {
		return 0, err
	}

	return u.ID, nil
}
