package datagenerator_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/postgresqldb"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/repository"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/services"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/testutils/datagenerator"
	"github.com/go-faker/faker/v4"
)

type user struct {
	First       string `faker:"-"`
	Last        string `faker:"-"`
	FirstFemale string `faker:"russian_first_name_female"`
	LastFemale  string `faker:"russian_last_name_female"`
	FirstMale   string `faker:"russian_first_name_male"`
	LastMale    string `faker:"russian_last_name_male"`
	Sex         byte   `faker:"oneof: 0,1"`
	Birthdate   string `faker:"date"`
	Biography   string `faker:"paragraph"`
	City        string `faker:"word"`
	Username    string `faker:"email"`
	Password    string `faker:"password"`
}

func randomUser() *domain.User {

	u := user{}

	faker.FakeData(&u)
	if u.Sex == byte(domain.Male) {
		u.First = u.FirstMale
		u.Last = u.LastMale
	} else {
		u.First = u.FirstFemale
		u.Last = u.LastFemale
	}

	parsedDate, _ := time.Parse("2006-01-02", u.Birthdate)

	// Создание пользователя с согласованным именем и полом
	return &domain.User{
		FirstName:    u.First,
		LastName:     u.Last,
		Sex:          domain.Sex(u.Sex),
		Birthdate:    parsedDate,
		Biography:    u.Biography,
		City:         u.City,
		Username:     u.Username,
		PasswordHash: domain.HashPassword(u.Password),
	}
}

func generateTestData(generator datagenerator.DataGenerator, numUsers, numFriends, numPosts int) error {
	users := make([]domain.UserKey, numUsers)
	sem := make(chan struct{}, 25) // Ограничиваем до 50 горутин одновременно

	// Создание пользователей
	var userErr error
	var wg sync.WaitGroup
	for i := 0; i < numUsers; i++ {
		sem <- struct{}{} // Блокируем, если 50 горутин уже запущено
		wg.Add(1)
		go func(i int) {
			defer func() { <-sem; wg.Done() }()
			u := randomUser()
			userId, err := generator.CreateUser(u)
			if err != nil {
				userErr = fmt.Errorf("failed to create user: %v\n", err)
				return
			}
			users[i] = userId
			if i%1000 == 0 {
				fmt.Printf("Creating users: %d of %d\n", i, numUsers)
			}
		}(i)
	}
	wg.Wait()
	if userErr != nil {
		return userErr
	}

	// Добавление друзей
	for i := 0; i < numFriends; i++ {
		userId := users[rand.Intn(len(users))]
		friendId := users[rand.Intn(len(users))]
		for friendId == userId {
			friendId = users[rand.Intn(len(users))]
		}
		sem <- struct{}{}
		wg.Add(1)
		go func(userId, friendId domain.UserKey, i int) {
			defer func() { <-sem; wg.Done() }()
			if err := generator.AddFriend(userId, friendId); err != nil {
				fmt.Printf("failed to invite user to friends %v: %v\n", userId, err)
			}
			if err := generator.AddFriend(friendId, userId); err != nil {
				fmt.Printf("failed to invite user to friends %v: %v\n", userId, err)
			}
			if i%1000 == 0 {
				fmt.Printf("Adding friends: %d of %d\n", i, numFriends)
			}
		}(userId, friendId, i)
	}
	wg.Wait()

	// Создание постов
	for i := 0; i < numPosts; i++ {
		userId := users[rand.Intn(len(users))]
		sem <- struct{}{}
		wg.Add(1)
		go func(userId domain.UserKey, i int) {
			defer func() { <-sem; wg.Done() }()
			_, err := generator.CreatePost(userId, domain.PostMessage(fmt.Sprintf("Test post %d. %v", i, faker.Paragraph())))
			if err != nil {
				fmt.Printf("failed to create post %v: %v\n", userId, err)
			}
			if i%1000 == 0 {
				fmt.Printf("Adding posts: %d of %d\n", i, numPosts)
			}
		}(userId, i)
	}
	wg.Wait()

	return nil
}

func TestGenerateTestData(t *testing.T) {

	dir, _ := os.Getwd()
	fmt.Printf("%s", dir)

	ctx := context.Background()

	cfg, err := config.LoadConfig("./../../../app-local.yaml")
	if err != nil {
		t.Fatalf("failed to open config: %v", err)
	}
	db, err := postgresqldb.InitDBCluster(ctx, cfg.SQLServer)
	if err != nil {
		t.Fatalf("failed to connect DB: %v", err)
	}
	bl := repository.NewBlacklistRepository(ctx, db)
	ur := repository.NewUserRepository(ctx, db)
	pr := repository.NewPostRepository(ctx, db)
	jwts, err := services.NewJWTService(cfg.JWT, bl)
	if err != nil {
		t.Fatalf("failed to create JWT Service: %v", err)
	}
	us := services.NewSocialNetworkService(cfg, ur, pr, nil, jwts, nil)

	gen := datagenerator.NewServiceDataGenerator(us)
	err = generateTestData(gen, 1000000, 10000000, 10000000)
	if err != nil {
		t.Fatalf("failed to generate test data: %v", err)
	}

	//u := randomUser()

	//fmt.Printf("u = %v", u)

}
