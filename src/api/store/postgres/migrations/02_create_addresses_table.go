package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func createAddressesTable(db migrations.DB) error {
	log.Debug("Creating addresses table...")
	_, err := db.Exec(`
CREATE TABLE addresses (
	id SERIAL PRIMARY KEY,
    wallet_id    SERIAL   NOT NULL,
    wallet_type    varchar(20)   NOT NULL,
    address    varchar(100)   NOT NULL,
	created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL, 
	updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL
);

ALTER TABLE "addresses"
    ADD FOREIGN KEY ("wallet_id") REFERENCES wallets ("id") ON DELETE CASCADE;
CREATE UNIQUE INDEX address_unique_address_idx ON addresses (wallet_id, wallet_type);

CREATE TRIGGER addresses_trigger
	BEFORE UPDATE ON addresses
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();
`)
	if err != nil {
		log.WithError(err).Error("Could not create accounts table")
		return err
	}
	log.Info("Created accounts table")

	return nil
}

func dropAddressesTable(db migrations.DB) error {
	log.Debug("Dropping addresses table")
	_, err := db.Exec(`
DROP TRIGGER addresses_trigger ON addresses;

DROP TABLE addresses;
`)
	if err != nil {
		log.WithError(err).Error("Could not drop address table")
		return err
	}
	log.Info("Dropped address table")

	return nil
}

func init() {
	Collection.MustRegisterTx(createAddressesTable, dropAddressesTable)
}
