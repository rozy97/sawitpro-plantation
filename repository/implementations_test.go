package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_CreateEstate(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	repo := Repository{Db: mockDB}

	t.Run("failed test case: database error", func(t *testing.T) {
		estate := Estate{Length: 15, Width: 25}

		mock.ExpectQuery(`INSERT INTO estates\(length, width\) VALUES \(\$1, \$2\) RETURNING id`).
			WithArgs(estate.Length, estate.Width).
			WillReturnError(sql.ErrConnDone)

		id, err := repo.CreateEstate(context.Background(), estate)
		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("success test case", func(t *testing.T) {
		estate := Estate{Length: 10, Width: 20}
		estateID := "some-uuid"

		mock.ExpectQuery(`INSERT INTO estates\(length, width\) VALUES \(\$1, \$2\) RETURNING id`).
			WithArgs(estate.Length, estate.Width).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(estateID))

		id, err := repo.CreateEstate(context.Background(), estate)
		assert.NoError(t, err)
		assert.Equal(t, estateID, id)
	})

}

func Test_CreateTree(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := Repository{Db: db}

	t.Run("failed test case: database error", func(t *testing.T) {
		tree := Tree{
			EstateID: "some-estate-id",
			X:        5,
			Y:        10,
			Height:   15,
		}

		mock.ExpectQuery(`INSERT INTO trees\(estate_id, x, y, height\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`).
			WithArgs(tree.EstateID, tree.X, tree.Y, tree.Height).
			WillReturnError(assert.AnError)

		id, err := repo.CreateTree(context.Background(), tree)
		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("success test case", func(t *testing.T) {
		tree := Tree{
			EstateID: "some-estate-id",
			X:        5,
			Y:        10,
			Height:   15,
		}
		treeID := "some-tree-id"

		mock.ExpectQuery(`INSERT INTO trees\(estate_id, x, y, height\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`).
			WithArgs(tree.EstateID, tree.X, tree.Y, tree.Height).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(treeID))

		id, err := repo.CreateTree(context.Background(), tree)
		assert.NoError(t, err)
		assert.Equal(t, treeID, id)
	})
}

func Test_GetEstateStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := Repository{Db: db}
	estateID := "some-uuid"

	t.Run("success case", func(t *testing.T) {
		expectedStats := Stats{
			TotalTrees: 10,
			MaxHeight:  15,
			MinHeight:  5,
			Median:     10,
		}

		mock.ExpectQuery(`SELECT COUNT\(\*\) AS total_trees, COALESCE\(MAX\(height\), 0\) AS max_height, COALESCE\(MIN\(height\), 0\) AS min_height, COALESCE\(PERCENTILE_CONT\(0.5\) WITHIN GROUP \(ORDER BY height\), 0\) AS median_height FROM trees WHERE estate_id = \$1`).
			WithArgs(estateID).
			WillReturnRows(sqlmock.NewRows([]string{"total_trees", "max_height", "min_height", "median_height"}).
				AddRow(expectedStats.TotalTrees, expectedStats.MaxHeight, expectedStats.MinHeight, expectedStats.Median))

		stats, err := repo.GetEstateStats(context.Background(), estateID)
		assert.NoError(t, err)
		assert.Equal(t, expectedStats, stats)
	})

	t.Run("failed test case: db error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT COUNT\(\*\).*WHERE estate_id = \$1`).
			WithArgs(estateID).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.GetEstateStats(context.Background(), estateID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
	})

	t.Run("success test case", func(t *testing.T) {
		expectedStats := Stats{
			TotalTrees: 10,
			MaxHeight:  15,
			MinHeight:  5,
			Median:     10,
		}

		mock.ExpectQuery(`SELECT COUNT\(\*\) AS total_trees, COALESCE\(MAX\(height\), 0\) AS max_height, COALESCE\(MIN\(height\), 0\) AS min_height, COALESCE\(PERCENTILE_CONT\(0.5\) WITHIN GROUP \(ORDER BY height\), 0\) AS median_height FROM trees WHERE estate_id = \$1`).
			WithArgs(estateID).
			WillReturnRows(sqlmock.NewRows([]string{"total_trees", "max_height", "min_height", "median_height"}).
				AddRow(expectedStats.TotalTrees, expectedStats.MaxHeight, expectedStats.MinHeight, expectedStats.Median))

		stats, err := repo.GetEstateStats(context.Background(), estateID)
		assert.NoError(t, err)
		assert.Equal(t, expectedStats, stats)
	})
}

func Test_GetEstateTrees(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := Repository{Db: db}
	estateID := "some-uuid"

	t.Run("failed case: db error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, estate_id, x, y, height FROM trees WHERE estate_id = \$1 ORDER BY x, y`).
			WithArgs(estateID).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.GetEstateTrees(context.Background(), estateID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
	})

	t.Run("failed case: row scan error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, estate_id, x, y, height FROM trees WHERE estate_id = \$1 ORDER BY x, y`).
			WithArgs(estateID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "estate_id", "x", "y", "height"}).
				AddRow(nil, nil, nil, nil, nil)) // Simulating a scan error

		_, err := repo.GetEstateTrees(context.Background(), estateID)
		assert.Error(t, err)
	})

	t.Run("success test case", func(t *testing.T) {
		expectedTrees := []Tree{
			{ID: "tree-1", EstateID: estateID, X: 1, Y: 2, Height: 10},
			{ID: "tree-2", EstateID: estateID, X: 2, Y: 3, Height: 15},
		}

		mock.ExpectQuery(`SELECT id, estate_id, x, y, height FROM trees WHERE estate_id = \$1 ORDER BY x, y`).
			WithArgs(estateID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "estate_id", "x", "y", "height"}).
				AddRow(expectedTrees[0].ID, expectedTrees[0].EstateID, expectedTrees[0].X, expectedTrees[0].Y, expectedTrees[0].Height).
				AddRow(expectedTrees[1].ID, expectedTrees[1].EstateID, expectedTrees[1].X, expectedTrees[1].Y, expectedTrees[1].Height))

		trees, err := repo.GetEstateTrees(context.Background(), estateID)
		assert.NoError(t, err)
		assert.Equal(t, expectedTrees, trees)
	})
}
