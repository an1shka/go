package main

import (
	"sync"
	"time"
)

var (
	studentsMutex sync.RWMutex
	students      []Student
	nextStudentID = 1

	teachersMutex sync.RWMutex
	teachers      []Teacher
	nextTeacherID = 1

	subjectsMutex sync.RWMutex
	subjects      []Subject
	nextSubjectID = 1

	gradesMutex sync.RWMutex
	grades      []Grade
	nextGradeID = 1

	homeworksMutex sync.RWMutex
	homeworks      []Homework
	nextHomeworkID = 1
)

func seedData() {
	students = []Student{
		{ID: 1, FullName: "Анна Петрова", Email: "anna@example.com", ClassName: "10A"},
		{ID: 2, FullName: "Борис Смирнов", Email: "boris@example.com", ClassName: "10A"},
		{ID: 3, FullName: "Вера Иванова", Email: "vera@example.com", ClassName: "10A"},
		{ID: 4, FullName: "Глеб Ким", Email: "gleb@example.com", ClassName: "10B"},
		{ID: 5, FullName: "Дарья Алиева", Email: "darya@example.com", ClassName: "10B"},
	}
	teachers = []Teacher{
		{ID: 1, FullName: "Ольга Соколова", Email: "sokolova@schoolhub.local"},
		{ID: 2, FullName: "Илья Морозов", Email: "morozov@schoolhub.local"},
		{ID: 3, FullName: "Мария Волкова", Email: "volkova@schoolhub.local"},
	}
	subjects = []Subject{
		{ID: 1, Title: "Математика", TeacherID: 1},
		{ID: 2, Title: "Русский язык", TeacherID: 2},
		{ID: 3, Title: "Информатика", TeacherID: 3},
		{ID: 4, Title: "Физика", TeacherID: 1},
	}

	now := time.Now().UTC()
	grades = []Grade{
		{ID: 1, StudentID: 1, SubjectID: 1, Value: 5, Date: now.AddDate(0, 0, -6), Comment: "Отличная работа"},
		{ID: 2, StudentID: 1, SubjectID: 3, Value: 4, Date: now.AddDate(0, 0, -5), Comment: "Хороший результат"},
		{ID: 3, StudentID: 2, SubjectID: 1, Value: 4, Date: now.AddDate(0, 0, -4), Comment: "Есть небольшие ошибки"},
		{ID: 4, StudentID: 2, SubjectID: 2, Value: 5, Date: now.AddDate(0, 0, -3), Comment: "Без ошибок"},
		{ID: 5, StudentID: 3, SubjectID: 2, Value: 3, Date: now.AddDate(0, 0, -2), Comment: "Нужно повторить правило"},
		{ID: 6, StudentID: 3, SubjectID: 4, Value: 4, Date: now.AddDate(0, 0, -1), Comment: "Зачтено"},
		{ID: 7, StudentID: 4, SubjectID: 3, Value: 5, Date: now, Comment: "Проект принят"},
		{ID: 8, StudentID: 5, SubjectID: 4, Value: 4, Date: now, Comment: "Хорошо"},
	}
	homeworks = []Homework{
		{ID: 1, SubjectID: 1, Title: "Квадратные уравнения", Description: "Решить номера 10-15", IssuedAt: now.AddDate(0, 0, -10), Deadline: now.AddDate(0, 0, -3)},
		{ID: 2, SubjectID: 2, Title: "Сочинение", Description: "Написать сочинение о любимой книге", IssuedAt: now.AddDate(0, 0, -5), Deadline: now.AddDate(0, 0, 2)},
		{ID: 3, SubjectID: 3, Title: "HTTP API", Description: "Составить список REST-методов", IssuedAt: now.AddDate(0, 0, -4), Deadline: now.AddDate(0, 0, 1)},
		{ID: 4, SubjectID: 4, Title: "Законы Ньютона", Description: "Подготовить краткий конспект", IssuedAt: now.AddDate(0, 0, -2), Deadline: now.AddDate(0, 0, 4)},
		{ID: 5, SubjectID: 1, Title: "Графики функций", Description: "Построить графики в тетради", IssuedAt: now, Deadline: now.AddDate(0, 0, 7)},
	}
	nextStudentID = 6
	nextTeacherID = 4
	nextSubjectID = 5
	nextGradeID = 9
	nextHomeworkID = 6
}
