// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	CreateEstate(ctx context.Context, estate Estate) (id string, err error)
	GetEstateByID(ctx context.Context, ID string) (estate Estate, err error)
	CreateTree(ctx context.Context, tree Tree) (id string, err error)
	GetEstateStats(ctx context.Context, ID string) (stats Stats, err error)
	GetEstateTrees(ctx context.Context, ID string) (trees []Tree, err error)
}
