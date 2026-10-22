INSERT INTO tasks(id,title,priority) VALUES('a','чай',1),('b','хлеб',2);
INSERT INTO task_events(task_id,action) VALUES('a','created'),('a','changed');
DO $$ DECLARE n integer; BEGIN
 BEGIN INSERT INTO task_events(task_id,action) VALUES('missing','created'); RAISE EXCEPTION 'orphan accepted'; EXCEPTION WHEN foreign_key_violation THEN NULL; END;
 SELECT count(*) INTO n FROM tasks t JOIN task_events e ON e.task_id=t.id; IF n<>2 THEN RAISE EXCEPTION 'JOIN count'; END IF;
 SELECT count(*) INTO n FROM tasks t LEFT JOIN task_events e ON e.task_id=t.id; IF n<>3 THEN RAISE EXCEPTION 'LEFT JOIN count'; END IF;
END $$;
