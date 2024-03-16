DELETE FROM Item;
DELETE FROM Column;
DELETE FROM Board;
DELETE FROM Password;
DELETE FROM Account;

INSERT INTO Account (id, email)
VALUES ('1', 'boarduser@email.com');

INSERT INTO Password (id, salt, hash, accountId)
VALUES ('1', 'RandomizedSaltValue', 'PlaceHolder', '1');

INSERT INTO Board (id, name, color, accountId)
VALUES (1, 'Personal tasks', '#e0e0e0', '1');

INSERT INTO Column (id, name, "order", boardId)
VALUES
('1', 'Todo', 1, 1),
('2', 'In-progress', 2, 1),
('3', 'Done', 3, 1);

INSERT INTO Item (id, title, content, "order", columnId, boardId)
VALUES
('1', 'Clean the house', 'Clean the kitchen, bathroom and living room.', 1, '1', 1),
('2', 'Buy groceries', 'Milk, Bread, Peanut Butter, Jelly', 2, '1', 1),
('3', 'Make dinner', 'Prepare the chicken curry for dinner.', 3, '1', 1);
