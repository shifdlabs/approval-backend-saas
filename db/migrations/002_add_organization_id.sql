-- =============================================================================
-- Migration: 002_add_organization_id
-- Phase 2 — Step 1: Multi-tenant schema changes (ADDITIVE ONLY)
--
-- Adds a nullable organization_id to every tenant-scoped table, seeds existing
-- rows with a default organization, enforces NOT NULL, then creates composite
-- (organization_id, <primary lookup column>) indexes.
--
-- This migration intentionally does NOT:
--   * drop or rename any existing column, and
--   * alter any existing UNIQUE constraint (see "FOLLOW-UPS" at the bottom).
--
-- It is wrapped in a single transaction (Postgres DDL is transactional) and is
-- idempotent / safe to re-run (IF NOT EXISTS + WHERE ... IS NULL + ON CONFLICT).
--
-- Tables that get a direct organization_id (per CLAUDE.md):
--   users, positions, documents, numbering_groups,
--   letter_templates, app_settings, user_logs
-- Child tables (document_sequences, document_histories, recipients, bookmarks,
-- signatures, delegators, ...) inherit the org through their parent and are
-- intentionally untouched here.
--
-- Prerequisite: the base tables already exist (created by GORM AutoMigrate on
-- app startup). Run this migration AFTER the app has created those tables.
-- =============================================================================

BEGIN;

-- Required for uuid_generate_v4(); the app already enables this on startup, but
-- we declare it here so the migration is self-contained.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- -----------------------------------------------------------------------------
-- 0. Parent table + default organization
--
--    The INSERT below targets `organizations`, which does NOT yet exist in the
--    Go models / AutoMigrate. It is created here so the migration is runnable
--    standalone. If `organizations` is owned by another service (e.g. SIS) or
--    by an earlier migration, delete this CREATE TABLE block and keep only the
--    INSERT. The column set (id/name/slug/timestamps) mirrors GORM conventions
--    so a future Organization model can AutoMigrate cleanly onto it.
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS organizations (
    id         uuid        PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       varchar(255) NOT NULL,
    slug       varchar(255) NOT NULL UNIQUE,
    created_at timestamptz  NOT NULL DEFAULT now(),
    updated_at timestamptz  NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

INSERT INTO organizations (id, name, slug)
VALUES ('00000000-0000-0000-0000-000000000001', 'Default Org', 'default-org')
ON CONFLICT (id) DO NOTHING;

-- -----------------------------------------------------------------------------
-- 1. Add organization_id (NULLABLE first — additive, requires no code changes).
-- -----------------------------------------------------------------------------
ALTER TABLE users            ADD COLUMN IF NOT EXISTS organization_id uuid;
ALTER TABLE positions        ADD COLUMN IF NOT EXISTS organization_id uuid;
ALTER TABLE documents        ADD COLUMN IF NOT EXISTS organization_id uuid;
ALTER TABLE numbering_groups ADD COLUMN IF NOT EXISTS organization_id uuid;
ALTER TABLE letter_templates ADD COLUMN IF NOT EXISTS organization_id uuid;
ALTER TABLE app_settings     ADD COLUMN IF NOT EXISTS organization_id uuid;
ALTER TABLE user_logs        ADD COLUMN IF NOT EXISTS organization_id uuid;

-- -----------------------------------------------------------------------------
-- 2. Seed existing rows with the default organization.
--
--    NOTE: we intentionally do NOT filter on deleted_at. These tables use GORM
--    soft deletes, so soft-deleted rows still physically exist and would violate
--    the NOT NULL constraint in step 3 if left unseeded. "WHERE organization_id
--    IS NULL" keeps the statement idempotent.
-- -----------------------------------------------------------------------------
UPDATE users            SET organization_id = '00000000-0000-0000-0000-000000000001' WHERE organization_id IS NULL;
UPDATE positions        SET organization_id = '00000000-0000-0000-0000-000000000001' WHERE organization_id IS NULL;
UPDATE documents        SET organization_id = '00000000-0000-0000-0000-000000000001' WHERE organization_id IS NULL;
UPDATE numbering_groups SET organization_id = '00000000-0000-0000-0000-000000000001' WHERE organization_id IS NULL;
UPDATE letter_templates SET organization_id = '00000000-0000-0000-0000-000000000001' WHERE organization_id IS NULL;
UPDATE app_settings     SET organization_id = '00000000-0000-0000-0000-000000000001' WHERE organization_id IS NULL;
UPDATE user_logs        SET organization_id = '00000000-0000-0000-0000-000000000001' WHERE organization_id IS NULL;

-- -----------------------------------------------------------------------------
-- 3. Enforce NOT NULL now that every row (including soft-deleted) is seeded.
--    Re-running SET NOT NULL on an already-NOT NULL column is a no-op.
-- -----------------------------------------------------------------------------
ALTER TABLE users            ALTER COLUMN organization_id SET NOT NULL;
ALTER TABLE positions        ALTER COLUMN organization_id SET NOT NULL;
ALTER TABLE documents        ALTER COLUMN organization_id SET NOT NULL;
ALTER TABLE numbering_groups ALTER COLUMN organization_id SET NOT NULL;
ALTER TABLE letter_templates ALTER COLUMN organization_id SET NOT NULL;
ALTER TABLE app_settings     ALTER COLUMN organization_id SET NOT NULL;
ALTER TABLE user_logs        ALTER COLUMN organization_id SET NOT NULL;

-- -----------------------------------------------------------------------------
-- 4. Composite indexes: (organization_id, <primary lookup column>).
--
--    The leading organization_id column also serves the plain
--    "WHERE organization_id = ?" listing queries that Step 4 will introduce.
--    The second column is the dominant per-table filter observed in the repos:
--      users            -> email      (UserRepository.GetByEmail)
--      positions        -> name       (PositionRepository.FindByName)
--      documents        -> author_id  (filter used across ~10 query methods)
--      numbering_groups -> name       (unique business key)
--      letter_templates -> name       (human lookup)
--      app_settings     -> key        (AppSettingsRepository.GetByKey)
--      user_logs        -> user_id    (per-user log listing / joins)
-- -----------------------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_users_org_email           ON users            (organization_id, email);
CREATE INDEX IF NOT EXISTS idx_positions_org_name        ON positions        (organization_id, name);
CREATE INDEX IF NOT EXISTS idx_documents_org_author      ON documents        (organization_id, author_id);
CREATE INDEX IF NOT EXISTS idx_numbering_groups_org_name ON numbering_groups (organization_id, name);
CREATE INDEX IF NOT EXISTS idx_letter_templates_org_name ON letter_templates (organization_id, name);
CREATE INDEX IF NOT EXISTS idx_app_settings_org_key      ON app_settings     (organization_id, "key");
CREATE INDEX IF NOT EXISTS idx_user_logs_org_user        ON user_logs        (organization_id, user_id);

COMMIT;

-- =============================================================================
-- FOLLOW-UPS (NOT done in this migration — tracked for later Phase 2 steps):
--
--  1. Per-tenant uniqueness. These columns currently have GLOBAL unique indexes
--     (created by GORM `uniqueIndex`): users.email, positions.name,
--     numbering_groups.name. For real isolation they should become unique PER
--     org, e.g. UNIQUE (organization_id, name). That is a drop+recreate of a
--     constraint, which is out of scope for this additive Step 1.
--
--  2. users.email index. idx_users_org_email is built on email, which CLAUDE.md
--     Step 5 will remove from this table (auth fields move to SIS). When email
--     is dropped this index is dropped automatically; add a replacement
--     (organization_id, id) composite at that time.
--
--  3. Foreign keys. No FK from <table>.organization_id -> organizations(id) is
--     added here (CLAUDE.md Step 1 lists only column + seed + NOT NULL + index).
--     Adding FKs is a reasonable hardening step once org provisioning is settled.
-- =============================================================================

-- =============================================================================
-- ROLLBACK (manual — apply only if you need to undo this migration):
--
-- BEGIN;
-- DROP INDEX IF EXISTS idx_users_org_email;
-- DROP INDEX IF EXISTS idx_positions_org_name;
-- DROP INDEX IF EXISTS idx_documents_org_author;
-- DROP INDEX IF EXISTS idx_numbering_groups_org_name;
-- DROP INDEX IF EXISTS idx_letter_templates_org_name;
-- DROP INDEX IF EXISTS idx_app_settings_org_key;
-- DROP INDEX IF EXISTS idx_user_logs_org_user;
-- ALTER TABLE users            DROP COLUMN IF EXISTS organization_id;
-- ALTER TABLE positions        DROP COLUMN IF EXISTS organization_id;
-- ALTER TABLE documents        DROP COLUMN IF EXISTS organization_id;
-- ALTER TABLE numbering_groups DROP COLUMN IF EXISTS organization_id;
-- ALTER TABLE letter_templates DROP COLUMN IF EXISTS organization_id;
-- ALTER TABLE app_settings     DROP COLUMN IF EXISTS organization_id;
-- ALTER TABLE user_logs        DROP COLUMN IF EXISTS organization_id;
-- -- DROP TABLE IF EXISTS organizations;  -- only if this migration created it
-- COMMIT;
-- =============================================================================
