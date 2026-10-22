-- Run in jobflow_lab after creating tasks and task_events from migrations01/02.
-- The final ROLLBACK restores rows, but consumed IDENTITY values are not reused.
BEGIN;
DELETE FROM tasks WHERE id IN ('join-a','join-b','join-c');
INSERT INTO tasks(id,title,priority) VALUES
 ('join-a','Чай',1),('join-b','Прогулка',2),('join-c','Повторить Go',3);
INSERT INTO task_events(task_id,action) VALUES
 ('join-a','created'),('join-a','changed'),('join-b','created');

-- 3 rows: a-created, a-changed, b-created.
SELECT t.id,t.title,e.action FROM tasks t
JOIN task_events e ON e.task_id=t.id
WHERE t.id IN ('join-a','join-b','join-c') ORDER BY t.id,e.id;

-- 4 rows: the preceding 3, plus c-NULL.
SELECT t.id,t.title,e.action FROM tasks t
LEFT JOIN task_events e ON e.task_id=t.id
WHERE t.id IN ('join-a','join-b','join-c') ORDER BY t.id,e.id;

-- All 3 tasks; unmatched event values are NULL.
SELECT t.id,e.action FROM tasks t
LEFT JOIN task_events e ON e.task_id=t.id AND e.action='changed'
WHERE t.id IN ('join-a','join-b','join-c') ORDER BY t.id,e.id;

-- Only a remains: WHERE removes the NULL right-hand side.
SELECT t.id,e.action FROM tasks t
LEFT JOIN task_events e ON e.task_id=t.id
WHERE t.id IN ('join-a','join-b','join-c') AND e.action='changed'
ORDER BY t.id,e.id;

-- a:2/2, b:1/1, c:0/1 (events / joined rows).
SELECT t.id,count(e.id) AS events,count(*) AS joined_rows FROM tasks t
LEFT JOIN task_events e ON e.task_id=t.id
WHERE t.id IN ('join-a','join-b','join-c')
GROUP BY t.id ORDER BY t.id;

DELETE FROM tasks WHERE id='join-b' RETURNING id;
SELECT task_id,action FROM task_events WHERE task_id='join-b'; -- no rows
ROLLBACK;
