package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func createWalletsTable(db migrations.DB) error {
	log.Debug("Creating wallets table...")
	_, err := db.Exec(`
CREATE TABLE wallets (
	id SERIAL PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    public_key TEXT NOT NULL,
    compressed_public_key TEXT,
    store_id TEXT,
    owner_id TEXT,
    active BOOLEAN default true,
	attributes JSONB,
	created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL, 
	updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL
);

CREATE UNIQUE INDEX wallet_unique_address_idx ON wallets (public_key);
CREATE OR REPLACE FUNCTION updated() RETURNS TRIGGER AS 
	$$
	BEGIN
		NEW.updated_at = (now() at time zone 'utc');
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

CREATE TRIGGER wallets_trigger
	BEFORE UPDATE ON wallets
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();
`)
	if err != nil {
		log.WithError(err).Error("Could not create wallets table")
		return err
	}
	log.Info("Created wallets table")

	return nil
}

func dropWalletsTable(db migrations.DB) error {
	log.Debug("Dropping wallets table")
	_, err := db.Exec(`
DROP TRIGGER wallets_trigger ON wallets;

DROP TABLE wallets;
`)
	if err != nil {
		log.WithError(err).Error("Could not drop wallet table")
		return err
	}
	log.Info("Dropped wallets table")

	return nil
}

func init() {
	Collection.MustRegisterTx(createWalletsTable, dropWalletsTable)
}
