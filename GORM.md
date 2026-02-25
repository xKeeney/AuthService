# GORM шпаргалка

## 1. Подключение и модель
```go
import (
    "gorm.io/driver/postgres" // или mysql, sqlite и т.д.
    "gorm.io/gorm"
)

// Пример модели
type User struct {
    gorm.Model            // встроенная модель: ID, CreatedAt, UpdatedAt, DeletedAt (мягкое удаление)
    Name         string
    Email        string `gorm:"uniqueIndex"`
    Age          int
    Active       bool
    Profile      Profile   // Один к одному
    Orders       []Order   // Один ко многим
}

type Profile struct {
    gorm.Model
    UserID uint
    Bio    string
}

type Order struct {
    gorm.Model
    UserID uint
    Total  float64
}

// Подключение к БД
dsn := "host=localhost user=postgres password=pass dbname=test port=5432 sslmode=disable"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
    panic("failed to connect database")
}

// Автомиграция (создание/обновление таблиц)
db.AutoMigrate(&User{}, &Profile{}, &Order{})
```

## 2. CREATE (Создание)
Одна запись
```go
user := User{Name: "Alice", Email: "alice@example.com", Age: 30}
result := db.Create(&user) // в user подставится ID, CreatedAt и т.д.
fmt.Println(result.Error)        // ошибка
fmt.Println(result.RowsAffected) // количество вставленных строк (1)
```
Выборочные поля
```go
db.Select("Name", "Email").Create(&user)
// INSERT INTO users (name, email) VALUES ('Alice', 'alice@example.com')
```
Пропуск полей
```go
db.Omit("Age", "Active").Create(&user)
// вставятся все поля кроме Age и Active (они получат значения по умолчанию)
```
Множественная вставка
```go
users := []User{{Name: "Bob"}, {Name: "Charlie"}}
db.Create(&users) // вернутся ID
```
Вставка с игнорированием конфликтов (PostgreSQL, MySQL)
```go
db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)
// Если email уникален, при конфликте ничего не делает
```
Обновление при конфликте (UPSERT)
```go
db.Clauses(clause.OnConflict{
    Columns:   []clause.Column{{Name: "email"}},
    DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),
}).Create(&user)
// при конфликте email обновит name и age
```

## 3. READ (Чтение)
Получение одной записи
```go
var user User

// Первая запись по первичному ключу (сортировка по PK)
db.First(&user, 1)                 // WHERE id = 1
db.First(&user, "email = ?", "alice@example.com") // условия

// Получить одну запись без сортировки
db.Take(&user)                      // LIMIT 1

// Последняя запись (сортировка по PK DESC)
db.Last(&user)
```
Все записи
```go
var users []User
db.Find(&users) // SELECT * FROM users

// с условием
db.Find(&users, "age > ?", 25)
```
Условия WHERE
```go
// Простые
db.Where("name = ?", "Alice").First(&user)
db.Where("name IN ?", []string{"Alice", "Bob"}).Find(&users)
db.Where("age BETWEEN ? AND ?", 20, 30).Find(&users)
db.Where("name LIKE ?", "%lice%").Find(&users)

// AND
db.Where("name = ? AND age >= ?", "Alice", 18).Find(&users)

// OR
db.Where("name = ?", "Alice").Or("age > ?", 30).Find(&users)

// NOT
db.Not("name = ?", "Alice").Find(&users)

// Структура или map в качестве условий
db.Where(&User{Name: "Alice", Active: true}).Find(&users) // только не zero-поля (Active=true учтётся)
db.Where(map[string]interface{}{"name": "Alice", "active": true}).Find(&users)
```
Выбор конкретных полей
```go
db.Select("name", "email").Find(&users)
// или
db.Select([]string{"name", "email"}).Find(&users)
```
Сортировка
```go
db.Order("age desc, name").Find(&users)
```
Лимит и смещение (пагинация)
```go
db.Limit(10).Offset(20).Find(&users) // страница 3 при size=10
```
Подсчёт количества
```go
var count int64
db.Model(&User{}).Where("active = ?", true).Count(&count)
```
Проверка наличия записи
```go
if err := db.Where("email = ?", "test@test.com").First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
    // не найдено
}
```
Подзапросы
```go
subQuery := db.Select("AVG(age)").Where("active = ?", true).Table("users")
db.Select("name, age").Where("age > (?)", subQuery).Find(&users)
```
Сканирование в другую структуру
```go
type Result struct {
    Name  string
    Total int
}
var results []Result
db.Model(&User{}).Select("name, count(*) as total").Group("name").Scan(&results)
```
Продвинутые условия (In, Like, и т.д. через конструкторы)
```go
import "gorm.io/gorm/clause"

db.Where(clause.Eq{Column: "name", Value: "Alice"})
db.Where(clause.Like{Column: "email", Value: "%@example%"})
```

## 4. UPDATE (Обновление)
Сохранение всех полей (Save – если первичный ключ есть, обновит все поля)
```go
user.Name = "Alicia"
user.Age = 31
db.Save(&user) // UPDATE users SET name='Alicia', age=31, ... WHERE id=...
```
Обновление одного поля
```go
db.Model(&user).Update("name", "Alice") // user должен содержать ID
// или
db.Model(&User{}).Where("email = ?", "old@test.com").Update("email", "new@test.com")
```
Обновление нескольких полей
```go
db.Model(&user).Updates(User{Name: "Alice", Age: 25}) // обновит только не-zero поля
db.Model(&user).Updates(map[string]interface{}{"name": "Alice", "age": 25}) // обновит все указанные
```
Обновление с условием (массовое)
```go
db.Model(&User{}).Where("active = ?", false).Update("active", true)
```
Обновление с выражением (например, инкремент)
```go
db.Model(&user).Update("age", gorm.Expr("age + ?", 1))
// или
db.Model(&User{}).Where("id > ?", 10).UpdateColumn("views", gorm.Expr("views + ?", 1))
```
Обновление из подзапроса
```go
db.Model(&user).Update("company_id", db.Table("companies").Select("id").Where("name = ?", "Acme"))
```
Пропуск хуков при обновлении
```go
db.Session(&gorm.Session{SkipHooks: true}).Model(&user).Update("name", "Alice")
```

## 5. DELETE (Удаление)
Мягкое удаление (если есть поле DeletedAt)
```go
db.Delete(&user) // устанавливает DeletedAt
db.Where("age < ?", 18).Delete(&User{}) // массовое мягкое удаление
```
Физическое удаление (навсегда)
```go
db.Unscoped().Delete(&user)          // игнорирует мягкое удаление
db.Unscoped().Where("age < ?", 18).Delete(&User{})
```
Удаление по условию (без загрузки)
```go
db.Where("email = ?", "spam@test.com").Delete(&User{})
```
Возврат удаленных данных (PostgreSQL)
```go
var deletedUsers []User
db.Clauses(clause.Returning{}).Where("age < ?", 18).Delete(&deletedUsers)
```

## 6. Работа со связями (Relations)
Предзагрузка (Eager Loading)
```go
var users []User
db.Preload("Profile").Preload("Orders").Find(&users)

// с условиями для предзагрузки
db.Preload("Orders", "total > ?", 100).Find(&users)

// вложенная предзагрузка
db.Preload("Profile").Preload("Orders.Items").Find(&users)
```
Joins (для фильтрации по связанным данным)
```go
db.Joins("Profile").Where("profile.bio LIKE ?", "%golang%").Find(&users)
// INNER JOIN profiles ON users.id = profiles.user_id
```
Обратная связь (принадлежит)
```go
type Profile struct {
    gorm.Model
    UserID uint
    User   User // Belongs To
}
// при запросе профиля можно подгрузить пользователя
var profile Profile
db.Preload("User").First(&profile, 1)
```
Добавление связанных записей (ассоциации)
```go
user := User{Name: "Alice"}
db.Create(&user)

profile := Profile{Bio: "Gopher"}
db.Model(&user).Association("Profile").Append(&profile) // установит user_id

order := Order{Total: 99.9}
db.Model(&user).Association("Orders").Append(&order)
```
Замена/удаление ассоциаций
```go
db.Model(&user).Association("Profile").Replace(&newProfile)
db.Model(&user).Association("Orders").Delete(&order)
db.Model(&user).Association("Orders").Clear()
```
Подсчет ассоциаций
```go
count := db.Model(&user).Association("Orders").Count()
```
Many-to-Many (пример)
```go
type Student struct {
    gorm.Model
    Name    string
    Courses []Course `gorm:"many2many:student_courses;"`
}
type Course struct {
    gorm.Model
    Title    string
    Students []Student `gorm:"many2many:student_courses;"`
}

// Добавление связи
student := Student{Name: "Ivan"}
course := Course{Title: "Math"}
db.Create(&student)
db.Create(&course)
db.Model(&student).Association("Courses").Append(&course)
```

## 7. Транзакции
Автоматическая транзакция (если нужно выполнить несколько операций атомарно)
```go
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err // откат
    }
    if err := tx.Create(&profile).Error; err != nil {
        return err
    }
    return nil // коммит
})
```
Ручное управление
```go
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()
if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}
tx.Commit()
```

## 8. Сырой SQL
```go
type Result struct {
    ID   int
    Name string
}
var results []Result
db.Raw("SELECT id, name FROM users WHERE age > ?", 18).Scan(&results)

// Exec для обновлений/удалений
db.Exec("UPDATE users SET active = ? WHERE age < ?", false, 13)
```

## 9. Миграции (изменение схемы)
```go
// Создание таблицы
db.AutoMigrate(&User{}, &Order{}) // безопасно: добавляет недостающие столбцы/индексы, не удаляет

// Дроп таблицы
db.Migrator().DropTable(&User{})

// Проверка наличия таблицы
db.Migrator().HasTable(&User{})

// Добавление индекса
db.Migrator().CreateIndex(&User{}, "Email")
```

## 10. Хуки (Callbacks)
```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // логика перед созданием
    return nil
}
func (u *User) AfterCreate(tx *gorm.DB) error {
    // после создания
    return nil
}
// Доступны: BeforeSave, AfterSave, BeforeUpdate, AfterUpdate, BeforeDelete, AfterDelete, AfterFind
```

## 11. Скоупы (переиспользуемые условия)
```go
func ActiveScope(db *gorm.DB) *gorm.DB {
    return db.Where("active = ?", true)
}

func AgeGreaterThan(age int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("age > ?", age)
    }
}

// Использование
db.Scopes(ActiveScope, AgeGreaterThan(20)).Find(&users)
```

## 12. Обработка ошибок
```go
err := db.First(&user, 100).Error
if errors.Is(err, gorm.ErrRecordNotFound) {
    fmt.Println("not found")
} else if err != nil {
    fmt.Println("other error:", err)
}
```