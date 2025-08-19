package post

import "BlogAPI/pkg/storage"

type PostHandler struct {
	Storage *storage.PostgresStorage
}

type ErrorResponce struct {
	Message string `json:"message"`
}

type EditPostRequest struct {
	Title *string `json:"title"`
	Content *string `json:"content"`
}