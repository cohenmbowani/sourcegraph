-- +++
-- parent: 1528395843
-- +++

BEGIN;

ALTER TABLE IF EXISTS batch_spec_executions ADD COLUMN IF NOT EXISTS user_id int REFERENCES users(id) DEFERRABLE;

COMMIT;
