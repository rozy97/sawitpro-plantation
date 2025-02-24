package repository

import (
	"context"
)

func (r *Repository) CreateEstate(ctx context.Context, estate Estate) (id string, err error) {
	err = r.Db.QueryRowContext(ctx, "INSERT INTO estates(length, width) VALUES ($1, $2) RETURNING id", estate.Length, estate.Width).Scan(&id)
	return
}

func (r *Repository) GetEstateByID(ctx context.Context, ID string) (estate Estate, err error) {
	err = r.Db.QueryRowContext(
		ctx,
		"SELECT id, length, width FROM estates WHERE id = $1", ID).Scan(
		&estate.ID,
		&estate.Length,
		&estate.Width,
	)

	return
}

func (r *Repository) CreateTree(ctx context.Context, tree Tree) (id string, err error) {
	err = r.Db.QueryRowContext(ctx, "INSERT INTO trees(estate_id, x, y, height) VALUES ($1, $2, $3, $4) RETURNING id", tree.EstateID, tree.X, tree.Y, tree.Height).Scan(&id)
	return
}

func (r *Repository) GetEstateStats(ctx context.Context, ID string) (stats Stats, err error) {
	err = r.Db.QueryRowContext(
		ctx,
		`SELECT 
			COUNT(*) AS total_trees,
			COALESCE(MAX(height), 0) AS max_height,
			COALESCE(MIN(height), 0) AS min_height,
			COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY height), 0) AS median_height
		FROM trees
		WHERE estate_id = $1`, ID).Scan(
		&stats.TotalTrees,
		&stats.MaxHeight,
		&stats.MinHeight,
		&stats.Median,
	)

	return
}

func (r *Repository) GetEstateTrees(ctx context.Context, ID string) ([]Tree, error) {
	trees := make([]Tree, 0)

	rows, err := r.Db.QueryContext(
		ctx,
		`SELECT id, estate_id, x, y, height
		FROM trees WHERE estate_id = $1
		ORDER BY x, y`,
		ID,
	)
	if err != nil {
		return trees, err
	}

	defer rows.Close()
	for rows.Next() {
		var tree Tree
		err = rows.Scan(
			&tree.ID,
			&tree.EstateID,
			&tree.X,
			&tree.Y,
			&tree.Height,
		)
		if err != nil {
			return trees, err
		}
		trees = append(trees, tree)
	}

	return trees, err
}
