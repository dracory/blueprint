package settingsonce

import (
	"context"
	"errors"
	"fmt"
	"project/internal/app"

	"github.com/dracory/base/cfmt"
	"github.com/dracory/neat/database/migrator"
	"github.com/dromara/carbon/v2"
)

var _ migrator.MigrationInterface = (*InstalledAtSeed)(nil)

// KeyInstalledAt is the setting key under which the install timestamp is
// recorded.
const KeyInstalledAt = "app.installed_at"

// InstalledAtSeed records the install timestamp in the setting store. It
// runs exactly once — applied seeds are recorded in the shared
// migration_tracker table under the seed_ signature, and the migrator
// skips them on subsequent boots.
type InstalledAtSeed struct {
	migrator.BaseMigration
	app app.AppInterface
}

// NewInstalledAtSeed creates a new InstalledAtSeed instance.
func NewInstalledAtSeed(a app.AppInterface) *InstalledAtSeed {
	return &InstalledAtSeed{app: a}
}

func (s *InstalledAtSeed) Signature() string {
	return "seed_2026_10_04_0001_settings_installed_at"
}

func (s *InstalledAtSeed) Description() string {
	return "Record the install timestamp in the setting store"
}

func (s *InstalledAtSeed) Up() error {
	return seedInstalledAt(context.Background(), s.app)
}

// Down is a no-op: the installed-at row may have been read or referenced
// since the seed ran, so it is left in place.
func (s *InstalledAtSeed) Down() error {
	return nil
}

func seedInstalledAt(ctx context.Context, app app.AppInterface) error {
	if app == nil {
		return errors.New("app is nil")
	}
	if app.IsDisabledSettingStore() {
		return errors.New("setting store is not initialized")
	}

	store := app.GetSettingStore()

	exists, err := store.Has(ctx, KeyInstalledAt)
	if err != nil {
		return fmt.Errorf("check %s: %w", KeyInstalledAt, err)
	}
	if exists {
		return nil
	}

	if err := store.Set(ctx, KeyInstalledAt,
		carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)); err != nil {
		return fmt.Errorf("set %s: %w", KeyInstalledAt, err)
	}

	cfmt.Successln("Seeded install timestamp:", KeyInstalledAt)

	return nil
}
