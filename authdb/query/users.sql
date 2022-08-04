-- name: CreateUser :exec
INSERT INTO users (
  UserName, Passwd
) VALUES (
  ?, ?
);

-- name: GetUser :exec
SELECT * FROM users
WHERE UserName = ? AND Passwd = ? 
LIMIT 1;

-- name: ListUser :many
SELECT * FROM users
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: UpdateUser :exec
UPDATE users
SET Passwd = ?
WHERE UserName = ?;

-- name: DeleteUser :execrows
DELETE FROM users WHERE UserName = ?;