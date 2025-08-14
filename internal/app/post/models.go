package post

import "BlogAPI/pkg/storage"

type PostHandler struct {
	storage *storage.PostgresStorage
}

type ErrorResponce struct {
	Message string `json:"message"`
}

