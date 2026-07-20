-- =============================================================================
-- Migration: 003_positions_per_tenant_unique
-- Phase 2 follow-up (flagged in 002_add_organization_id.sql FOLLOW-UPS #1).
--
-- positions.name currently has a GLOBAL unique index (idx_positions_name),
-- left over from single-tenant. That means once any org creates a position
-- named e.g. "Manager", no other org can ever create a position with that
-- name — the insert fails with 23505 duplicate key. Confirmed in prod data:
-- org 55959b43-08d1-4bb3-900a-bf5ff9182d88 already owns "Manager", which
-- blocks every other org from using that name.
--
-- Fix: drop the global unique index and replace the existing plain
-- (organization_id, name) lookup index with a composite UNIQUE index, so
-- uniqueness is enforced per-org instead of globally.
-- =============================================================================

BEGIN;

DROP INDEX IF EXISTS idx_positions_name;
DROP INDEX IF EXISTS idx_positions_org_name;

CREATE UNIQUE INDEX IF NOT EXISTS idx_positions_org_name ON positions (organization_id, name);

COMMIT;

-- =============================================================================
-- ROLLBACK (manual):
--
-- BEGIN;
-- DROP INDEX IF EXISTS idx_positions_org_name;
-- CREATE INDEX IF NOT EXISTS idx_positions_org_name ON positions (organization_id, name);
-- CREATE UNIQUE INDEX IF NOT EXISTS idx_positions_name ON positions (name);
-- COMMIT;
-- =============================================================================
