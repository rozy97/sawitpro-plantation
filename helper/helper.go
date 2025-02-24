package helper

import (
	"math"

	"github.com/SawitProRecruitment/UserService/repository"
)

type Stats struct {
	Estate              repository.Estate
	Trees               Trees
	Rest                Rest
	Distance            int
	CurrentHeight       int
	CountFirstRest      bool
	IsFirstRestResolved bool
	MaxDistance         int
}

type Trees []repository.Tree
type Rest struct {
	X int
	Y int
}

func (t *Trees) GetTreeByCoordinate(x, y int) repository.Tree {
	for _, tree := range *t {
		if tree.X == x && tree.Y == y {
			return tree
		}
	}

	return repository.Tree{}
}

func (s *Stats) CalculateDistance(x, y int) {
	tree := s.Trees.GetTreeByCoordinate(x, y)

	s.Distance += int(math.Abs(float64(s.CurrentHeight) - float64(tree.Height+1)))
	s.CurrentHeight = tree.Height + 1

	if s.CountFirstRest && !s.IsFirstRestResolved {
		if s.Distance+s.CurrentHeight <= s.MaxDistance {
			s.Rest.X = x
			s.Rest.Y = y
		} else {
			s.IsFirstRestResolved = true
		}
	}

	if x < s.Estate.Length || y < s.Estate.Width {
		s.Distance += 10
	}

}

func (s *Stats) CalculateTotalDistance() {
	s.CurrentHeight = 1
	s.Distance = 1
	for width := 1; width <= s.Estate.Width; width++ {
		if width%2 == 1 {
			for length := 1; length <= s.Estate.Length; length++ {
				s.CalculateDistance(length, width)
			}
		} else {
			for length := s.Estate.Length; length >= 1; length-- {
				s.CalculateDistance(length, width)
			}
		}
	}
	s.Distance += 1
}
