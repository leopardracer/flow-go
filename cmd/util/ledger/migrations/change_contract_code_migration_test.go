package migrations_test

import (
	"context"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/flow-go/cmd/util/ledger/migrations"
	"github.com/onflow/flow-go/ledger"
	"github.com/onflow/flow-go/ledger/common/convert"
	"github.com/onflow/flow-go/model/flow"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChangeContractCodeMigration(t *testing.T) {

	address1, err := common.HexToAddress("0x1")
	address2, err := common.HexToAddress("0x2")
	require.NoError(t, err)

	ctx := context.Background()

	contractPayload := func(address common.Address, contractName string, contract []byte) *ledger.Payload {
		return ledger.NewPayload(
			convert.RegisterIDToLedgerKey(
				flow.ContractRegisterID(flow.ConvertAddress(address), contractName),
			),
			contract,
		)
	}

	t.Run("no contracts", func(t *testing.T) {
		migration := migrations.ChangeContractCodeMigration{}
		log := zerolog.New(zerolog.NewTestWriter(t))
		err := migration.InitMigration(log, nil, 0)
		require.NoError(t, err)

		_, err = migration.MigrateAccount(ctx, address1,
			[]*ledger.Payload{},
		)

		require.NoError(t, err)

		err = migration.Close()
		require.NoError(t, err)
	})

	t.Run("1 contract - dont migrate", func(t *testing.T) {
		migration := migrations.ChangeContractCodeMigration{}
		log := zerolog.New(zerolog.NewTestWriter(t))
		err := migration.InitMigration(log, nil, 0)
		require.NoError(t, err)

		payloads, err := migration.MigrateAccount(ctx, address1,
			[]*ledger.Payload{
				contractPayload(address1, "A", []byte("A")),
			},
		)

		require.NoError(t, err)
		require.Len(t, payloads, 1)
		require.Equal(t, []byte("A"), []byte(payloads[0].Value()))

		err = migration.Close()
		require.NoError(t, err)
	})

	t.Run("1 contract - migrate", func(t *testing.T) {
		migration := migrations.ChangeContractCodeMigration{}
		log := zerolog.New(zerolog.NewTestWriter(t))
		err := migration.InitMigration(log, nil, 0)
		require.NoError(t, err)

		migration.ChangeContract(address1, "A", "B")

		payloads, err := migration.MigrateAccount(ctx, address1,
			[]*ledger.Payload{
				contractPayload(address1, "A", []byte("A")),
			},
		)

		require.NoError(t, err)
		require.Len(t, payloads, 1)
		require.Equal(t, []byte("B"), []byte(payloads[0].Value()))

		err = migration.Close()
		require.NoError(t, err)
	})

	t.Run("2 contracts - migrate 1", func(t *testing.T) {
		migration := migrations.ChangeContractCodeMigration{}
		log := zerolog.New(zerolog.NewTestWriter(t))
		err := migration.InitMigration(log, nil, 0)
		require.NoError(t, err)

		migration.ChangeContract(address1, "A", "B")

		payloads, err := migration.MigrateAccount(ctx, address1,
			[]*ledger.Payload{
				contractPayload(address1, "A", []byte("A")),
				contractPayload(address1, "B", []byte("A")),
			},
		)

		require.NoError(t, err)
		require.Len(t, payloads, 2)
		require.Equal(t, []byte("B"), []byte(payloads[0].Value()))
		require.Equal(t, []byte("A"), []byte(payloads[1].Value()))

		err = migration.Close()
		require.NoError(t, err)
	})

	t.Run("2 contracts - migrate 2", func(t *testing.T) {
		migration := migrations.ChangeContractCodeMigration{}
		log := zerolog.New(zerolog.NewTestWriter(t))
		err := migration.InitMigration(log, nil, 0)
		require.NoError(t, err)

		migration.ChangeContract(address1, "A", "B")
		migration.ChangeContract(address1, "B", "B")

		payloads, err := migration.MigrateAccount(ctx, address1,
			[]*ledger.Payload{
				contractPayload(address1, "A", []byte("A")),
				contractPayload(address1, "B", []byte("A")),
			},
		)

		require.NoError(t, err)
		require.Len(t, payloads, 2)
		require.Equal(t, []byte("B"), []byte(payloads[0].Value()))
		require.Equal(t, []byte("B"), []byte(payloads[1].Value()))

		err = migration.Close()
		require.NoError(t, err)
	})

	t.Run("2 contracts on different accounts - migrate 1", func(t *testing.T) {
		migration := migrations.ChangeContractCodeMigration{}
		log := zerolog.New(zerolog.NewTestWriter(t))
		err := migration.InitMigration(log, nil, 0)
		require.NoError(t, err)

		migration.ChangeContract(address1, "A", "B")

		payloads, err := migration.MigrateAccount(ctx, address1,
			[]*ledger.Payload{
				contractPayload(address1, "A", []byte("A")),
				contractPayload(address2, "A", []byte("A")),
			},
		)

		require.NoError(t, err)
		require.Len(t, payloads, 2)
		require.Equal(t, []byte("B"), []byte(payloads[0].Value()))
		require.Equal(t, []byte("A"), []byte(payloads[1].Value()))

		err = migration.Close()
		require.NoError(t, err)
	})

	t.Run("not all contracts on one account migrated", func(t *testing.T) {
		migration := migrations.ChangeContractCodeMigration{}
		log := zerolog.New(zerolog.NewTestWriter(t))
		err := migration.InitMigration(log, nil, 0)
		require.NoError(t, err)

		migration.ChangeContract(address1, "A", "B")
		migration.ChangeContract(address1, "B", "B")

		_, err = migration.MigrateAccount(ctx, address1,
			[]*ledger.Payload{
				contractPayload(address1, "A", []byte("A")),
			},
		)

		require.Error(t, err)
	})

	t.Run("not all accounts migrated", func(t *testing.T) {
		migration := migrations.ChangeContractCodeMigration{}
		log := zerolog.New(zerolog.NewTestWriter(t))
		err := migration.InitMigration(log, nil, 0)
		require.NoError(t, err)

		migration.ChangeContract(address2, "A", "B")

		_, err = migration.MigrateAccount(ctx, address1,
			[]*ledger.Payload{
				contractPayload(address1, "A", []byte("A")),
			},
		)

		require.NoError(t, err)

		err = migration.Close()
		require.Error(t, err)
	})
}
