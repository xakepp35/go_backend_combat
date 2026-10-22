-- Run only in a dedicated teaching database.
INSERT INTO tasks(id,title,priority,done) SELECT 'bench-'||n,'Задача '||n,1+(n%3),n%10<>0 FROM generate_series(1,50000) n ON CONFLICT(id) DO NOTHING;
ANALYZE tasks;
EXPLAIN (ANALYZE, BUFFERS) SELECT id,title FROM tasks WHERE done=false ORDER BY priority DESC,id LIMIT 20;
