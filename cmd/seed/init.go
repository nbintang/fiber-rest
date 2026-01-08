package main

import (
	"fmt"
	"math/rand"
	"rest-fiber/internal/category"
	"rest-fiber/internal/enums"
	"rest-fiber/internal/post"
	"rest-fiber/internal/user"
	"rest-fiber/pkg"
	"strings"

	"github.com/go-faker/faker/v4"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Options struct {
	Count int
}

func InitSeeds(db *gorm.DB, opt Options) {
	domain := "example.com"

	if err := resetTables(db); err != nil {
		panic(err)
	}

	categoryPool := []string{
		"Technology", "Programming", "Lifestyle", "Backend", "DevOps",
		"Database", "Security", "Testing",
	}
	hashed, err := pkg.HashPassword("Password123")
	if err != nil {
		return
	}
	for i := 0; i < opt.Count; i++ {
		name := faker.Name()

		userSeed := user.User{
			Name:            name,
			Email:           realisticEmail(name, i, domain),
			AvatarURL:       realisticAvatar(name),
			Password:        hashed,
			Role:            evenOddRole(i),
			IsEmailVerified: rand.Intn(2) == 1,
		}

		for p := 0; p < 3; p++ {
			catName := categoryPool[rand.Intn(len(categoryPool))]

			userSeed.Posts = append(userSeed.Posts, post.Post{
				Title:    realisticPostTitle(),
				Body:     realisticPostBody(),
				ImageURL: faker.URL(),
				Status:   evenOddStatus(i + p),
				Category: category.Category{Name: catName},
			})
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "email"}},
			DoNothing: true,
		}).Create(&userSeed).Error; err != nil {
			fmt.Printf("Error when create user: %s\n", userSeed.Name)
			return
		}
	}
}

func evenOddRole(i int) user.Role {
	if i%2 == 0 {
		return user.Role(enums.Member)
	}
	return user.Role(enums.Admin)
}
func evenOddStatus(i int) post.Status {
	if i%2 == 0 {
		return post.Status(enums.Published)
	}
	return post.Status(enums.Draft)
}

func realisticEmail(fullName string, i int, domain string) string {
	// domain tanpa "@"
	domain = strings.TrimPrefix(domain, "@")
	user := slug.Make(fullName)
	return fmt.Sprintf("%s.%d@%s", user, i+1, domain)
}

func realisticAvatar(fullName string) string {
	// gampang & rapi, avatar dari initials
	// contoh: https://ui-avatars.com/api/?name=John+Doe&background=random
	name := strings.ReplaceAll(fullName, " ", "+")
	return fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=random", name)
}

func realisticPostTitle() string {
	// bikin judul yang lebih “blog-like”
	templates := []string{
		"How to %s in Go (Fiber + GORM)",
		"Guide: %s for REST API",
		"%s: Best Practices & Examples",
		"Building %s with Fiber",
		"Debugging %s in Production",
	}
	topics := []string{
		"pagination", "JWT auth", "database migrations", "validation",
		"error handling", "repository pattern", "dependency injection",
		"transaction handling", "soft delete", "logging",
	}
	return fmt.Sprintf(templates[rand.Intn(len(templates))], topics[rand.Intn(len(topics))])
}

func realisticPostBody() string {
	// faker.Paragraph() biasanya lebih enak dari Sentence()
	// kalau faker versi kamu tidak punya Paragraph, ganti jadi faker.Sentence() beberapa kali
	return faker.Paragraph()
}

func resetTables(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("TRUNCATE TABLE posts RESTART IDENTITY CASCADE").Error; err != nil {
			return err
		}
		if err := tx.Exec("TRUNCATE TABLE categories RESTART IDENTITY CASCADE").Error; err != nil {
			return err
		}
		if err := tx.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error; err != nil {
			return err
		}
		return nil
	})
}
