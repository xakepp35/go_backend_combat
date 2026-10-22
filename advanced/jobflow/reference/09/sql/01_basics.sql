INSERT INTO drafts(id,doc) VALUES('example','{"kind":"echo"}');
SELECT id,doc FROM drafts ORDER BY id;
UPDATE drafts SET doc='{"kind":"sum"}' WHERE id='example';
DELETE FROM drafts WHERE id='example';
