package moviesRepository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gormModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/gorm"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRepository_CreateMovieWithCastList(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	r := &Postgres{DB: gormDB}

	type mockBehavior func(movie gormModels.Movie, castList []uint64, id uint64)

	tests := []struct {
		name            string
		mockBehavior    mockBehavior
		movie           gormModels.Movie
		castList        []uint64
		movieID         uint64
		isResponseError bool
	}{
		{
			name: "ok",
			mockBehavior: func(movie gormModels.Movie, castList []uint64, id uint64) {
				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO movies").
					WithArgs(movie.Title, movie.Description, movie.ReleaseDate, movie.Rating).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

				for _, v := range castList {
					mock.ExpectQuery("INSERT INTO actor_movie_relations").
						WithArgs(id, v).
						WillReturnRows()
				}

				mock.ExpectCommit()
			},
			movie:           gormModels.Movie{},
			castList:        []uint64{1},
			movieID:         1,
			isResponseError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.movie, test.castList, test.movieID)

			_, err := r.CreateMovieWithCastList(test.movie, test.castList)
			if test.isResponseError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
