-- name: CreateBoard :one
INSERT INTO Board (name, color, accountId)
VALUES (?, ?, ?)
RETURNING *;

-- name: UpdateBoardName :exec
UPDATE Board
SET name = ?
WHERE id = ? AND accountId = ?;

-- name: DeleteBoard :exec
DELETE FROM Board
WHERE id = ? AND accountId = ?;

-- name: GetHomeData :many
SELECT * FROM Board
WHERE accountId = ?;

-- name: GetBoard :one
SELECT * FROM Board
WHERE id = ? AND accountId = ?;

-- name: DeleteCard :exec
DELETE FROM Item
WHERE "Item".id = ? AND "Item".boardId = (SELECT id FROM Board WHERE accountId = ? LIMIT 1);

-- name: GetBoardItems :many
SELECT * FROM Item
WHERE boardId = ?;

-- name: GetBoardColumns :many
SELECT * FROM Column
WHERE boardId = ?
ORDER BY "order" ASC;

-- name: CreateColumn :one
INSERT INTO Column (id, boardId, name, "order")
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateColumnName :exec
UPDATE Column
SET name = ?
WHERE "Column".id = ? AND "Column".boardId = (SELECT id FROM Board WHERE accountId = ? LIMIT 1);
