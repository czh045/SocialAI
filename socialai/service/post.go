package service

import (
	"mime/multipart"
	"reflect"

	"socialai/backend"
	"socialai/constants"
	"socialai/model"

	"github.com/olivere/elastic/v7"
)

func SearchPostsByUser(user string) ([]model.Post, error) {
	query := elastic.NewTermQuery("user", user)
	result, err := backend.ESBackend.ReadFromES(query, constants.PostIndex)
	if err != nil {
		return nil, err
	}
	return postsFromSearchResult(result), nil
}

func SearchPostsByKeywords(keywords string) ([]model.Post, error) {
	query := elastic.NewMatchQuery("message", keywords).Operator("AND")
	if keywords == "" {
		query.ZeroTermsQuery("all")
	}

	result, err := backend.ESBackend.ReadFromES(query, constants.PostIndex)
	if err != nil {
		return nil, err
	}
	return postsFromSearchResult(result), nil
}

func SavePost(post *model.Post, file multipart.File, extension string) error {
	mediaURL, err := backend.MediaBackend.Save(file, post.ID+extension)
	if err != nil {
		return err
	}

	post.URL = mediaURL
	return backend.ESBackend.SaveToES(post, constants.PostIndex, post.ID)
}

func DeletePost(id, user string) error {
	query := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("id", id)).
		Must(elastic.NewTermQuery("user", user))
	result, err := backend.ESBackend.ReadFromES(query, constants.PostIndex)
	if err != nil {
		return err
	}

	posts := postsFromSearchResult(result)
	if len(posts) == 0 {
		return nil
	}

	// The media object is removed as part of deletion, so local test uploads
	// and production GCS objects do not accumulate after a post disappears.
	if err := backend.MediaBackend.Delete(posts[0].URL); err != nil {
		return err
	}
	return backend.ESBackend.DeleteFromES(query, constants.PostIndex)
}

func postsFromSearchResult(result *elastic.SearchResult) []model.Post {
	var posts []model.Post
	var postType model.Post
	for _, item := range result.Each(reflect.TypeOf(postType)) {
		posts = append(posts, item.(model.Post))
	}
	return posts
}
