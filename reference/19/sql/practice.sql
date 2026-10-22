INSERT INTO tasks(id,title,priority) VALUES('a','Купить чай',1);
SELECT id,title,priority FROM tasks WHERE done=false ORDER BY priority DESC,id;
UPDATE tasks SET title='Купить зелёный чай' WHERE id='a';
DELETE FROM tasks WHERE id='a';
