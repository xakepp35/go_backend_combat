SELECT t.id,t.title,e.action FROM tasks t LEFT JOIN task_events e ON e.task_id=t.id ORDER BY t.id,e.id;
