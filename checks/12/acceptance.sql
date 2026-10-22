INSERT INTO tasks(id,title,priority) VALUES('a','чай',1);
DO $$ BEGIN
 BEGIN INSERT INTO tasks(id,title,priority) VALUES('a','кофе',1); RAISE EXCEPTION 'duplicate accepted'; EXCEPTION WHEN unique_violation THEN NULL; END;
 BEGIN INSERT INTO tasks(id,title,priority) VALUES('b','',1); RAISE EXCEPTION 'empty accepted'; EXCEPTION WHEN check_violation THEN NULL; END;
 BEGIN INSERT INTO tasks(id,title,priority) VALUES('b','x',4); RAISE EXCEPTION 'priority accepted'; EXCEPTION WHEN check_violation THEN NULL; END;
END $$;
