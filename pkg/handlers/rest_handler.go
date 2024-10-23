package handler

import (
	srv "github.com/mccor2000/cinema/pkg/service"
	str "github.com/mccor2000/cinema/pkg/storage"
)

type RestHandler struct {	
	service *srv.Service
	store *str.Storage
}
