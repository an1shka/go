package main

import "sync"

var (
	teachersMutex sync.Mutex
	teachers      = []Teacher{}
	nextTeacherID = 1

	subjectsMutex sync.Mutex
	subjects      = []Subject{}
	nextSubjectID = 1

	gradesMutex sync.Mutex
	grades      = []Grade{}
	nextGradeID = 1

	homeworksMutex sync.Mutex
	homeworks      = []Homework{}
	nextHomeworkID = 1
)