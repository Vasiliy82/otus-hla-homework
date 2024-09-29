# ДЗ №2. Отчет о проделанной работе
## 1. Постановка задачи

Описание/Пошаговая инструкция выполнения домашнего задания:
	1. Сгенерировать любым способ 1,000,000 анкет. Имена и Фамилии должны быть реальными, чтобы учитывать селективность индекса. Так же можно воспользоваться [уже готовым списком](https://raw.githubusercontent.com/OtusTeam/highload/master/homework/people.v2.csv "уже готовым списком") как основой.
	2. Реализовать функционал поиска анкет по префиксу имени и фамилии (одновременно) в вашей социальной сети (реализовать метод [/user/search из спецификации](https://github.com/OtusTeam/highload/blob/master/homework/openapi.json#L165 "/user/search из спецификации")) (запрос в форме firstName LIKE ? and secondName LIKE ?). Сортировать вывод по id анкеты.
	3. Провести нагрузочные тесты этого метода. Поиграть с количеством одновременных запросов. 1/10/100/1000.
	4. Построить графики и сохранить их в отчет
	5. Сделать подходящий индекс.
	6. Повторить пункт 3 и 4.
	7. В качестве результата предоставить отчет в котором должны быть:
		- графики latency до индекса;
		- графики throughput до индекса;
		- графики latency после индекса;
		- графики throughput после индекса;
		- запрос добавления индекса;
		- explain запросов после индекса;
		- объяснение почему индекс именно такой.
	ДЗ принимается в виде отчета по выполненной работе.

## 2. Генерация тестовых данных

Т.к. критерии генерации фамилий и имен не были заданы, принято решение взять за основу [готовый список](https://raw.githubusercontent.com/OtusTeam/highload/master/homework/people.v2.csv "готовый список") анкет с реальными именами и фамилиями для улучшения селективности индекса. Для упрощения загрузки данных в БД потребовалось скорректировать формат выгрузки. Для этого разработана утилита на Go. Код представлен в листинге.
```Go
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
```

## 3. Функционал поиска анкет по префиксу имени и фамилии

Метод поиска реализован через endpoint `/api/user/search`, который принимает входные параметры `first_name` и `last_name`. В основе метода лежит запрос SQL с префиксным поиском с использованием оператора `LIKE`:
```sql
SELECT * FROM users WHERE first_name LIKE ? AND last_name LIKE ? ORDER BY id;
```
### 3.1 Новый метод в слое представления
`~/internal/rest/rest.go`
```Go
func (h *userHandler) Search(c echo.Context) error {
	// Извлечение query параметров first_name и last_name
	firstName := c.QueryParam("first_name")
	lastName := c.QueryParam("last_name")

	// Валидация параметров
	if !isValidName(firstName) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат имени"})
	}
	if !isValidName(lastName) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат фамилии"})
	}

	users, err := h.userService.Search(firstName, lastName)
	if err != nil {
		var apperr *apperrors.AppError
		if errors.As(err, &apperr) {
			return c.JSON(apperr.Code, map[string]string{"error": apperr.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}
```
### 3.2 Новый метод в слое бизнес-логики
`~/internal/services/user.go`
``` Go
func (s *userService) Search(firstName, lastName string) ([]*domain.User, error) {
	users, err := s.userRepo.Search(firstName, lastName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFoundError("User not found")
		}
		return nil, apperrors.NewInternalServerError("UserService.Login: s.userRepo.GetByUserName returned unknown error", err)
	}

	return users, nil
}
```
### 3.3 Новый метод в слое хранения
`~/internal/repository/user.go`
```Go
func (r *userRepository) Search(firstName, lastName string) ([]*domain.User, error) {
	var users []*domain.User

	ptnFirstName := fmt.Sprintf("%s%%", firstName)
	ptnLastName := fmt.Sprintf("%s%%", lastName)

	q, err := r.db.Query("SELECT id, first_name, last_name, birthdate, biography, city, username, password_hash, created_at, updated_at FROM users WHERE first_name LIKE $1 AND last_name LIKE $2 ORDER BY id", ptnFirstName, ptnLastName)
	if err != nil {
		return nil, fmt.Errorf("userRepository.Search: r.db.Query returned error %w", err)
	}
	defer q.Close()

	for q.Next() {
		user := domain.User{}
		err := q.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Birthdate,
			&user.Biography, &user.City, &user.Username, &user.PasswordHash,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("userRepository.Search: q.Scan returned error %w", err)
		}
		users = append(users, &user)
	}

	return users, nil
}
```
### 4. Конфигурация и методология тестирования

Тестирование проводилось на виртуальной машине с параметрами:
- **Процессор**: 6 ядер CPU (1.3 Ггц, до 4.6 ГГц в режиме Turbo).
- **Оперативная память**: 4096 MB.
- **Накопитель**: SSD.

Инструмент для тестирования — Apache JMeter 5.6.3. 
Проводились запросы к `/api/user/search` в 1, 10, 100 и 1000 одновременных запросов. Значения параметров `first_name` и `last_name` генерировались по формуле `${__RandomString(1,АБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ)}` для одного символа и `${__RandomString(1,АБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ)}${__RandomString(1,абвгдеёжзийклмнопрстуфхцчшщъыьэюя)}` для двух символов.

Для мониторинга использовался Listener Aggregate Report, который фиксировал следующие показатели:
- **Среднее время отклика (Average)**.
- **Перцентиль 90, 95, 99**.
- **Throughput (пропускная способность)**.
- **Error Rate (процент ошибок)**.

### 5. Результаты тестирования до индексации

![[lat-before-1.png]] 

![[thr-before-1.png]]

![[lat-before-2.png]]

![[thr-before-2.png]]



**Результаты (ключевые метрики):**
- **Average** (время отклика) было значительно выше, особенно при высоком количестве потоков.
- **Throughput** для 1000 потоков составил ~144 запросов/сек.
- **Error %** был высоким — около 90% запросов завершались с ошибкой.

### 6. Оптимизация

#### 6.1 До оптимизации
В качестве примера был рассмотрен следующий запрос:
```sql
EXPLAIN ANALYZE
SELECT id, first_name, last_name, birthdate, biography, city, username, password_hash, created_at, updated_at
FROM users
WHERE first_name LIKE 'А%' AND last_name LIKE 'Б%'
ORDER BY id;
```

Результаты выполнения `EPLAIN ANALYZE`:

```
Gather Merge  (cost=75106.39..76371.84 rows=10846 width=237) (actual time=791.416..795.397 rows=11832 loops=1)
  Workers Planned: 2
  Workers Launched: 2
  ->  Sort  (cost=74106.36..74119.92 rows=5423 width=237) (actual time=786.584..787.042 rows=3944 loops=3)
        Sort Key: id
        Sort Method: quicksort  Memory: 1128kB
        Worker 0:  Sort Method: quicksort  Memory: 1114kB
        Worker 1:  Sort Method: quicksort  Memory: 1145kB
        ->  Parallel Seq Scan on users  (cost=0.00..73770.01 rows=5423 width=237) (actual time=401.433..784.665 rows=3944 loops=3)
              Filter: (((first_name)::text ~~ 'А%'::text) AND ((last_name)::text ~~ 'Б%'::text))
              Rows Removed by Filter: 329390
Planning Time: 3.887 ms
Execution Time: 795.799 ms
```

**Основные проблемы:**
- Запрос выполнялся через **Parallel Seq Scan**, что означало последовательное сканирование всей таблицы для поиска соответствующих записей. Это приводит к значительным затратам на чтение данных с диска, особенно в больших таблицах с миллионами строк.
- **Filter** (фильтрация) применялась постфактум, что означает, что все строки, не удовлетворяющие условиям, удалялись только после полного сканирования.
- Несмотря на использование параллелизма, время выполнения запроса оставалось высоким.

#### 6.2 После оптимизации

Для оптимизации был создан комбинированный индекс `ix_users_first_last` с использованием оператора `text_pattern_ops` для ускорения префиксного поиска:

```sql
CREATE INDEX ix_users_first_last ON users(first_name text_pattern_ops, last_name text_pattern_ops);
```

Результаты после добавления индексов:

```
Sort  (cost=37629.63..37662.17 rows=13015 width=237) (actual time=8.257..9.523 rows=11832 loops=1)
  Sort Key: id
  ->  Bitmap Heap Scan on users  (cost=4239.21..36740.20 rows=13015 width=237) (actual time=2.380..5.580 rows=11832 loops=1)
        Filter: (((first_name)::text ~~ 'А%'::text) AND ((last_name)::text ~~ 'Б%'::text))
        Heap Blocks: exact=1868
        ->  Bitmap Index Scan on ix_users_first_last  (cost=0.00..4235.96 rows=12876 width=0) (actual time=2.192..2.192 rows=11832 loops=1)
              Index Cond: (((first_name)::text ~>=~ 'А'::text) AND ((first_name)::text ~<~ 'Б'::text) AND ((last_name)::text ~>=~ 'Б'::text) AND ((last_name)::text ~<~ 'В'::text))
Planning Time: 0.364 ms
Execution Time: 9.775 ms
```

**Что изменилось:**
**Bitmap Index Scan (Bitmap-сканирование индекса)**:
   - После добавления комбинированного индекса на поля `first_name` и `last_name`, оптимизатор начал использовать более эффективную стратегию — **Bitmap Index Scan**. 
   - Суть в том, что теперь запрос не сканирует всю таблицу, а использует индекс, который позволяет мгновенно найти только те строки, которые соответствуют условиям фильтрации (`first_name LIKE 'А%'` и `last_name LIKE 'Б%'`). Это сокращает объем данных, которые нужно просканировать.
   - Время выполнения индексации — **2.192 мс** — значительно меньше, чем время полного сканирования таблицы.
**Bitmap Heap Scan**:
   - После использования индекса происходит **Bitmap Heap Scan**, который читает данные только из тех блоков таблицы, где находятся нужные строки (1868 блоков). Это значительно сокращает количество прочитываемых данных.
   - Время выполнения чтения из хранилища — **2.380 - 5.580 мс**, что гораздо быстрее по сравнению с предыдущим подходом с полным сканированием таблицы.
**Сортировка (Sort)**:
   - Сортировка по полю `id` теперь выполняется быстрее, так как обрабатывается значительно меньшее количество данных после применения индексов.
   - Время сортировки составило **8.257 - 9.523 мс**, что тоже значительно быстрее, чем до индексации (более 780 мс).
**Отсутствие параллелизма**:
   - Запрос больше не требует параллелизма, так как использование индекса и оптимизированных операций значительно уменьшило объем работы. В данном случае весь запрос выполняется на одном потоке за меньшее время, чем при параллельном выполнении.

#### 6.3 Почему индекс ускоряет поиск?

Индексы с `text_pattern_ops` оптимизированы для префиксных поисков по текстовым полям. Когда запрос использует конструкцию `LIKE 'А%'`, индекс позволяет быстро найти все записи, начинающиеся на букву "А", не прибегая к полному сканированию таблицы. Это значительно сокращает объем данных, которые нужно проверять и сортировать. Добавление этих индексов позволило:
- Уменьшить количество обрабатываемых строк.
- Ускорить время выполнения запроса более чем в 50 раз (с ~795 мс до ~14 мс).
- Снизить нагрузку на процессор и диск за счет уменьшения объема данных для обработки.

### 7. Результаты тестирования после индексации

![[lat-after-1.png]] 

![[thr-after-1.png]]

![[lat-after-2.png]]

![[thr-after-2.png]]


После добавления индексов результаты значительно улучшились:
- **Average** время отклика упало до 1-10 мс при малом количестве потоков.
- **Throughput** значительно увеличился и составил 568-954 запросов/сек в зависимости от нагрузки.
- **Error %** существенно снизился и стал практически нулевым, хотя при 1000 потоков остались единичные ошибки (0.25% ошибок при 1000 потоках).

### 8. Наблюдения и выводы

#### 8.1 Наблюдения

В процессе тестирования было сделано неожиданное наблюдение: при высокой нагрузке кроме 500 Internal Server Error наблюдались также ошибки 401 Unauthorized. Ошибки 401 возвращались при передаче заведомо валидного токена, что требует дальнейшего исследования.

#### 8.2 Выводы

Добавление индексов дало существенный прирост производительности:
- **Throughput** вырос более чем в 4 раза.
- **Average latency** упала до минимальных значений (1-10 мс) при малом числе потоков.
- Ошибки снизились до минимальных значений.

Индексы на поля `first_name` и `last_name` с использованием `text_pattern_ops` оказались оптимальными для решения задачи префиксного поиска. Они позволяют эффективно работать с запросами `LIKE` и минимизируют количество полных сканирований таблицы.

#### 8.3 Будущие шаги

Для устранения проблемы с ошибками 401 Unauthorized потребуется детальное исследование механизма авторизации и обработки токенов. Возможно, проблема кроется в механизмах проверки токенов при высокой нагрузке.
