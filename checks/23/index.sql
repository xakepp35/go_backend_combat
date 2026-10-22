DO $$ BEGIN IF NOT EXISTS(SELECT 1 FROM pg_indexes WHERE indexname='tasks_open_priority_idx' AND schemaname=current_schema()) THEN RAISE EXCEPTION 'index missing'; END IF; END $$;
