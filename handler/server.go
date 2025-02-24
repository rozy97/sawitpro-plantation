package handler

import "github.com/SawitProRecruitment/UserService/repository"

type Server struct {
	Repository repository.RepositoryInterface
}

func NewServer(r repository.RepositoryInterface) *Server {
	return &Server{Repository: r}
}
