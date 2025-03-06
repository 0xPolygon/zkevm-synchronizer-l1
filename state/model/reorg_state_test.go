package model_test

import (
	"context"
	"path"
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/model"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/sqlstorage"
	"github.com/stretchr/testify/require"
)

func TestReorg(t *testing.T) {
	storage, err := sqlstorage.NewSqlStorage(sqlstorage.Config{
		DriverName: sqlstorage.SqliteDriverName,
		DataSource: path.Join(t.TempDir(), "TestReorg.sqlite"),
	}, true)
	require.NoError(t, err)
	sut := model.NewReorgState(storage)
	require.NotNil(t, sut)
	tx, err := storage.BeginTransaction(context.TODO())
	require.NoError(t, err)
	ctx := context.TODO()
	storage.AddBlock(ctx, &sqlstorage.L1Block{BlockNumber: 10}, tx)
	storage.AddBlock(ctx, &sqlstorage.L1Block{BlockNumber: 20}, tx)
	storage.AddBlock(ctx, &sqlstorage.L1Block{BlockNumber: 30}, tx)
	storage.AddSequencedBatches(ctx, &sqlstorage.SequencedBatches{
		FromBatchNumber: 1,
		ToBatchNumber:   10,
		L1BlockNumber:   10,
	}, tx)
	storage.AddSequencedBatches(ctx, &sqlstorage.SequencedBatches{
		FromBatchNumber: 11,
		ToBatchNumber:   12,
		L1BlockNumber:   20,
	}, tx)
	storage.AddSequencedBatches(ctx, &sqlstorage.SequencedBatches{
		FromBatchNumber: 13,
		ToBatchNumber:   14,
		L1BlockNumber:   30,
	}, tx)
	err = tx.Commit(ctx)
	require.NoError(t, err)

	tx, err = storage.BeginTransaction(context.TODO())
	require.NoError(t, err)
	res := sut.ExecuteReorg(context.TODO(), model.ReorgRequest{
		FirstL1BlockNumberToKeep: 10,
		ReasonError:              nil,
	}, tx)
	require.NoError(t, res.ExecutionError)
	err = tx.Commit(ctx)
	require.NoError(t, err)

	lastBLock, err := storage.GetLastBlock(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, uint64(10), lastBLock.BlockNumber)

	lastSeq, err := storage.GetLatestSequence(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, uint64(10), lastSeq.ToBatchNumber)
}
