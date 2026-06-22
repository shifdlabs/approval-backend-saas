package seeder

import "gorm.io/gorm"

func Run(db *gorm.DB) error {

	// ── Tier 1: Position (no FK) ─────────────────
	if err := SeedPositions(db); err != nil {
		return err
	}

	// ── Tier 2: User (needs PositionID) ──────────
	if err := SeedUsers(db); err != nil {
		return err
	}

	return nil
}
