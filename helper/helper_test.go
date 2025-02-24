package helper

import (
	"testing"

	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/stretchr/testify/assert"
)

func Test_CalculateTotalDistance(t *testing.T) {
	t.Run("calculate total distance with max_param", func(t *testing.T) {
		estateID := "estate-123"
		estate := repository.Estate{
			ID:     estateID,
			Length: 5,
			Width:  4,
		}

		trees := Trees{
			repository.Tree{
				ID:       "tree-1",
				EstateID: estateID,
				X:        2,
				Y:        1,
				Height:   5,
			},
			repository.Tree{
				ID:       "tree-2",
				EstateID: estateID,
				X:        3,
				Y:        1,
				Height:   3,
			},
			repository.Tree{
				ID:       "tree-3",
				EstateID: estateID,
				X:        4,
				Y:        1,
				Height:   4,
			},
		}

		stats := Stats{
			Estate:         estate,
			Trees:          trees,
			CountFirstRest: true,
			MaxDistance:    100,
		}

		stats.CalculateTotalDistance()

		assert.Equal(t, 204, stats.Distance)
		assert.Equal(t, 2, stats.Rest.X)
		assert.Equal(t, 2, stats.Rest.Y)
	})

}
