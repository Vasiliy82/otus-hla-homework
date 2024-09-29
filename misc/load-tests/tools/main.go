package main

/*
Утилита предназначена для парсинга исходного input.csv (ссылка на него была приложена к ТЗ) и конвертации его в формат, совместимый со схемой БД проекта.
*/

import (
	"crypto/md5"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type User struct {
	FirstName    string
	LastName     string
	BirthDate    string
	City         string
	Biography    string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func hashPassword(password string) string {
	h := md5.New()
	io.WriteString(h, password)
	hash := h.Sum(nil)
	return fmt.Sprintf("%x", hash)
}

func main() {
	// Инициализация для генерации данных
	gofakeit.Seed(0)

	// Чтение CSV-файла
	file, err := os.Open("input.csv")
	if err != nil {
		log.Fatalf("Unable to read input file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Unable to parse file as CSV: %v", err)
	}

	// Подготовка данных для вставки
	var users []User
	usedEmails := make(map[string]bool) // Хранилище для уникальных email

	// Задаем диапазон для случайной даты (последние 3 года)
	startDate := time.Now().AddDate(-3, 0, 0) // 3 года назад
	endDate := time.Now()                     // сегодня

	for _, record := range records {
		// Парсинг данных из CSV
		fullName := record[0]
		birthDate := record[1]
		city := record[2]

		// Разделяем фамилию и имя
		names := strings.SplitN(fullName, " ", 2)
		if len(names) != 2 {
			log.Printf("Skipping record due to invalid name format: %v", record)
			continue
		}
		lastName, firstName := names[0], names[1]

		// Генерация уникального email
		var username string
		for {
			username = gofakeit.Email() // Генерация случайного email
			if !usedEmails[username] {  // Проверяем на уникальность
				usedEmails[username] = true
				break
			}
		}

		// Генерация биографии и случайного пароля
		biography := gofakeit.Paragraph(1, 3, 5, ".")
		passwordHash := hashPassword(gofakeit.Password(true, true, true, false, false, 16))

		// Генерация случайных дат CreatedAt и UpdatedAt
		createdAt := gofakeit.DateRange(startDate, endDate) // Случайная дата за последние 3 года
		updatedAt := gofakeit.DateRange(createdAt, endDate) // UpdatedAt всегда >= CreatedAt

		// Добавляем сгенерированные данные в структуру User
		users = append(users, User{
			FirstName:    firstName,
			LastName:     lastName,
			BirthDate:    birthDate,
			City:         city,
			Biography:    biography,
			Username:     username,
			PasswordHash: passwordHash,
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		})
	}

	// Запись обогащенных данных в новый CSV файл
	outputFile, err := os.Create("output.csv")
	if err != nil {
		log.Fatalf("Unable to create output file: %v", err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// Запись заголовка
	writer.Write([]string{"FirstName", "LastName", "BirthDate", "City", "Biography", "Username", "PasswordHash", "CreatedAt", "UpdatedAt"})

	for _, user := range users {
		writer.Write([]string{
			user.FirstName,
			user.LastName,
			user.BirthDate,
			user.City,
			user.Biography,
			user.Username,
			user.PasswordHash,
			user.CreatedAt.Format(time.RFC3339),
			user.UpdatedAt.Format(time.RFC3339),
		})
	}

	fmt.Println("Data written to output.csv")
}
