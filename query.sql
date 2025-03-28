-- name: CreateUser :one
INSERT INTO users (username, email, password)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ?;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?;

-- name: EditUser :exec
UPDATE users
SET username = ?, email = ?, password = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;
