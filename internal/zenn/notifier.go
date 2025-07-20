package zenn

import (
	"context"
	"fmt"
	"neko-bot/internal/errors"
	"neko-bot/redis"
	"sync"
	"time"
)

type WatchUser struct {
	UserId   string
	Articles []Article
}

type ArticleNotifier struct {
	WatchUsers     []*WatchUser
	NewArticleChan chan []Article
	client         Repository
	mu             sync.Mutex
}

func NewArticleNotifier() *ArticleNotifier {
	client := redis.Client()
	users, err := client.SMembers(context.Background(), redis.WatchedZennUsers).Result()
	if err != nil {
		errors.Catch(err, "redis client fetch err: watched_zenn_users")
	}

	watchUsers := make([]*WatchUser, 0, len(users))
	for _, user := range users {
		watchUsers = append(watchUsers, &WatchUser{
			UserId:   user,
			Articles: make([]Article, 0),
		})
	}

	return &ArticleNotifier{
		WatchUsers:     watchUsers,
		NewArticleChan: make(chan []Article),
		client:         NewRepository(),
	}
}

func (n *ArticleNotifier) Init() {
	for _, user := range n.WatchUsers {
		articles, err := n.client.FetchArticlesByUsername(context.Background(), user.UserId)
		if err != nil {
			errors.Catch(err, fmt.Sprintf("fetch articles err: %s", user.UserId))
			return
		}

		if len(articles) == 0 {
			continue
		}

		user.Articles = make([]Article, len(articles))
		copy(user.Articles, articles)
	}
}

func (n *ArticleNotifier) Start() {
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:
			n.mu.Lock()
			for _, users := range n.WatchUsers {
				newArticles, err := n.fetchNewArticles(users)

				if len(newArticles) == 0 || err != nil {
					continue
				}

				n.NewArticleChan <- newArticles
			}
			n.mu.Unlock()
		}
	}
}

func (n *ArticleNotifier) fetchNewArticles(user *WatchUser) ([]Article, error) {
	articles, err := n.client.FetchArticlesByUsername(context.Background(), user.UserId)
	if err != nil {
		errors.Catch(err, fmt.Sprintf("fetch articles err: %s", user.UserId))
		return nil, err
	}

	if len(articles) == 0 {
		return nil, nil
	}

	newArticles := make([]Article, 0)

	if len(user.Articles) == 0 {
		newArticles = append(newArticles, articles[0])
		user.Articles = articles
		return newArticles, nil
	}

	latestExistingTime := user.Articles[0].PublishedAt
	// 新しく取得した記事の中で、既存の最新記事より新しいものを探す
	for _, article := range articles {
		if article.PublishedAt.After(latestExistingTime) {
			newArticles = append(newArticles, article)
		} else {
			// これより古い記事なので終了
			break
		}
	}

	// 記事リストを更新（新しい記事を保持）
	if len(newArticles) > 0 {
		user.Articles = articles
	}

	return newArticles, nil
}

func (n *ArticleNotifier) AddUser(userId string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	err := redis.Client().SAdd(redis.Context(), redis.WatchedZennUsers, userId).Err()
	if err != nil {
		errors.Catch(err, fmt.Sprintf("add user err: %s", userId))
		return err
	}

	articles, err := n.client.FetchArticlesByUsername(context.Background(), userId)
	if err != nil {
		errors.Catch(err, fmt.Sprintf("fetch articles err: %s", userId))
		return err
	}

	newUser := &WatchUser{
		UserId:   userId,
		Articles: make([]Article, len(articles)),
	}
	copy(newUser.Articles, articles)

	n.WatchUsers = append(n.WatchUsers, newUser)
	return nil
}

func (n *ArticleNotifier) RemoveUser(userId string) (bool, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	err := redis.Client().SRem(redis.Context(), redis.WatchedZennUsers, userId).Err()
	if err != nil {
		errors.Catch(err, fmt.Sprintf("remove user err: %s", userId))
		return false, err
	}

	foundIndex := -1
	for i, user := range n.WatchUsers {
		if user.UserId == userId {
			foundIndex = i
			break
		}
	}
	if foundIndex != -1 {
		n.WatchUsers = append(n.WatchUsers[:foundIndex], n.WatchUsers[foundIndex+1:]...)
		return true, nil
	}

	return false, nil
}
