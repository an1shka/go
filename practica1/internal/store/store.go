package store

import (
	"errors"
	"fmt"
	"sync"

	"main/internal/models"
)

// Ошибки хранилища
var (
	ErrUserExists   = errors.New("email уже зарегистрирован")
	ErrUserNotFound = errors.New("пользователь не найден")
	ErrStudentNotFound = errors.New("ученик не найден")
)

// Store — потокобезопасное хранилище данных
type Store struct {
	mu       sync.RWMutex
	users    map[string]models.User    // email -> User
	tokens   map[string]string         // token -> userID
	students map[string]models.Student // id -> Student
	nextID   int
}

// NewStore создает и инициализирует новое хранилище
func NewStore() *Store {
	return &Store{
		users:    make(map[string]models.User),
		tokens:   make(map[string]string),
		students: make(map[string]models.Student),
		nextID:   1,
	}
}

// CreateUser добавляет нового пользователя, проверяя уникальность email
func (s *Store) CreateUser(email, password string) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[email]; exists {
		return models.User{}, ErrUserExists
	}

	userID := fmt.Sprintf("u%d", len(s.users)+1)
	user := models.User{
		ID:       userID,
		Email:    email,
		Password: password, // В реальном проекте хешируется (например, bcrypt)
	}

	s.users[email] = user
	return user, nil
}

// GetUserByEmail ищет пользователя по email
func (s *Store) GetUserByEmail(email string) (models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[email]
	return user, exists
}

// SaveToken связывает токен с ID пользователя
func (s *Store) SaveToken(token, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = userID
}

// GetUserByToken возвращает пользователя по его токену
func (s *Store) GetUserByToken(token string) (models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, exists := s.tokens[token]
	if !exists {
		return models.User{}, false
	}

	for _, u := range s.users {
		if u.ID == userID {
			return u, true
		}
	}
	return models.User{}, false
}

// CreateStudent добавляет нового ученика
func (s *Store) CreateStudent(st models.Student) models.Student {
	s.mu.Lock()
	defer s.mu.Unlock()

	st.ID = fmt.Sprintf("%d", s.nextID)
	s.nextID++
	s.students[st.ID] = st
	return st
}

// GetAllStudents возвращает список всех учеников
func (s *Store) GetAllStudents() []models.Student {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]models.Student, 0, len(s.students))
	for _, st := range s.students {
		list = append(list, st)
	}
	return list
}

// GetStudentByID ищет ученика по ID
func (s *Store) GetStudentByID(id string) (models.Student, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, exists := s.students[id]
	return st, exists
}

// UpdateStudent обновляет существующего ученика
func (s *Store) UpdateStudent(id string, updated models.Student) (models.Student, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.students[id]; !exists {
		return models.Student{}, ErrStudentNotFound
	}

	updated.ID = id
	s.students[id] = updated
	return updated, nil
}

// DeleteStudent удаляет ученика по ID
func (s *Store) DeleteStudent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.students[id]; !exists {
		return ErrStudentNotFound
	}

	delete(s.students, id)
	return nil
}