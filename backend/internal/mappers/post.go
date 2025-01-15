package mappers

import (
	"time"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/dto"
)

func ToPostMessage(request *dto.CreateOrUpdatePostRequest) domain.PostMessage {
	return domain.PostMessage(request.Message)
}

func ToPostResponse(post *domain.Post) *dto.GetPostResponse {
	return &dto.GetPostResponse{
		Id:         int64(post.Id),
		Message:    string(post.Message),
		CreatedAt:  post.CreatedAt,
		ModifiedAt: post.ModifiedAt,
	}
}

func ToUpdatedPost(post *domain.Post, request *dto.CreateOrUpdatePostRequest) {
	post.Message = domain.PostMessage(request.Message)
	modifiedAt := time.Now()
	post.ModifiedAt = &modifiedAt
}

func ToPostsResponse(posts []*domain.Post) *dto.GetFeedResponse {
	var lastId domain.PostKey = 0
	if len(posts) > 0 && posts[len(posts)-1] != nil {
		lastId = posts[len(posts)-1].Id
	}
	return &dto.GetFeedResponse{
		Feed:       posts,
		LastPostId: lastId,
	}
}
