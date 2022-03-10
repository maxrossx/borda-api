package postgres_test

import (
	"borda/internal/core"
	"borda/internal/core/entity"
	"borda/internal/data/repository/postgres"

	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewTaskRepository(t *testing.T) {
	db := MustConnectAndMigrate(t)

	taskRepository := postgres.NewTaskRepository(db)

	require.Implements(t, (*core.TaskRepository)(nil), taskRepository)
}

func Test_TaskRepository_CreateNewTask(t *testing.T) {
	type testCase struct {
		Name string

		Task entity.Task

		ExpectedInt   int
		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			db := MustConnectAndMigrate(t)
			defer MustCloseDB(t, db)

			taskRepository := postgres.NewTaskRepository(db)

			actualInt, actualError := taskRepository.CreateNewTask(tc.Task)

			require.Equal(t, tc.ExpectedInt, actualInt)
			require.Equal(t, tc.ExpectedError, actualError)
		})
	}

	validate(t, &testCase{
		Name: "OK",
		Task: entity.Task{
			Title: "Test Task",
			Author: entity.Author{
				Name:    "Test Team",
				Contact: "test@mail.com",
			},
		},
		ExpectedInt:   1,
		ExpectedError: nil,
	})
}

func Test_TaskRepository_GetTasks(t *testing.T) {
	type testCase struct {
		Name string

		Filter entity.TaskFilter

		ExpectedSlice []*entity.Task
		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			db := MustConnectAndMigrate(t)
			defer MustCloseDB(t, db)

			for _, i := range tc.ExpectedSlice {
				MustCreateTask(t, db, i)
			}

			taskRepository := postgres.NewTaskRepository(db)

			actualSlice, actualError := taskRepository.GetTasks(tc.Filter)

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	validate(t, &testCase{
		Name: "Select active",
		Filter: entity.TaskFilter{
			IsActive: true,
		},
		ExpectedSlice: []*entity.Task{
			&entity.Task{
				Id:       1,
				Title:    "Test Task 1",
				IsActive: true,
				Author: entity.Author{
					Id:      1,
					Name:    "Test Team",
					Contact: "test@mail.com",
				},
			},
			&entity.Task{
				Id:       2,
				Title:    "Test Task 2",
				IsActive: true,
				Author: entity.Author{
					Id:      1,
					Name:    "Test Team",
					Contact: "test@mail.com",
				},
			},
			&entity.Task{
				Id:       3,
				Title:    "Test Task 3",
				IsActive: true,
				Author: entity.Author{
					Id:      1,
					Name:    "Test Team",
					Contact: "test@mail.com",
				},
			},
		},
		ExpectedError: nil,
	})
}

func Test_TaskRepository_GetTaskById(t *testing.T) {
	type testCase struct {
		Name string

		TaskId int

		ExpectedTask  *entity.Task
		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			db := MustConnectAndMigrate(t)
			defer MustCloseDB(t, db)

			MustCreateTask(t, db, tc.ExpectedTask)

			taskRepository := postgres.NewTaskRepository(db)

			actualTask, actualError := taskRepository.GetTaskById(tc.TaskId)

			assert.Equal(t, tc.ExpectedTask, actualTask)
			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	validate(t, &testCase{
		Name:   "OK",
		TaskId: 1,
		ExpectedTask: &entity.Task{
			Id:          1,
			Title:       "Test Title",
			Description: "Test description",
			Category:    "Test",
			Complexity:  "",
			Points:      100,
			Hint:        "",
			Flag:        "flag{_test_}",
			IsActive:    false,
			IsDisabled:  false,
			Author: entity.Author{
				Id:      1,
				Name:    "Test Name",
				Contact: "test@mail.com",
			},
		},
		ExpectedError: nil,
	})
}

// func Test_TaskRepository_UpdateTask(t *testing.T) {
// 	type testCase struct {
// 		Name string

// 		Task   entity.Task
// 		Update entity.TaskUpdate

// 		ExpectedError error
// 	}

// 	validate := func(t *testing.T, tc *testCase) {
// 		t.Run(tc.Name, func(t *testing.T) {
// 			db := MustConnectAndMigrate(t)
// 			defer MustCloseDB(t, db)

// 			MustCreateTask(t, db, &tc.Task)

// 			taskRepository := postgres.NewTaskRepository(db)

// 			actualError := taskRepository.UpdateTask(tc.Task.Id, tc.Update)
// 			actualTask, _ := taskRepository.GetTaskById(tc.Task.Id)

// 			assert.Equal(t, tc.ExpectedError, actualError)
// 		})
// 	}

// 	validate(t, &testCase{
// 		Name: "",
// 		Task: entity.Task{
// 			Title:       "Test Task",
// 			Description: "",
// 			Category:    "",
// 			Complexity:  "",
// 			Points:      0,
// 			Hint:        "",
// 			Flag:        "",
// 			IsActive:    false,
// 			IsDisabled:  false,
// 			Author:      entity.Author{},
// 		},
// 		Update: entity.TaskUpdate{
// 			Title:         "Test Task After Update",
// 			Description:   "New description",
// 			Category:      "Web",
// 			Complexity:    "Easy",
// 			Points:        400,
// 			Hint:          "Super hint",
// 			Flag:          "flag{flag}",
// 			AuthorName:    "New name",
// 			AuthorContact: "New Contact",
// 		},
// 		ExpectedError: nil,
// 	})
// }

// MustCreateTask creates a task in the database. Fatal on error.
func MustCreateTask(t *testing.T, db *sqlx.DB, task *entity.Task) *entity.Task {
	t.Helper()

	id, err := postgres.NewTaskRepository(db).CreateNewTask(*task)
	if err != nil {
		t.Fatal(err)
	}

	task.Id = id

	return task
}
