package services

import (
	"context"
	"errors"
	"testing"

	storecache "github.com/shellhub-io/shellhub/api/cache"
	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/api/store/mocks"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/shellhub-io/shellhub/pkg/validator"
	"github.com/stretchr/testify/assert"
)

func TestUpdateDataUser(t *testing.T) {
	mock := &mocks.Store{}
	s := NewService(store.Store(mock), privateKey, publicKey, storecache.NewNullCache(), clientMock, nil)

	ctx := context.TODO()

	Err := errors.New("conflict")

	user1 := &models.User{Name: "name", Email: "user1@email.com", Username: "username1", Password: "hash1", ID: "id1"}

	user2 := &models.User{Name: "name", Email: "user2@email.com", Username: "username2", Password: "hash1", ID: "id2"}

	type Expected struct {
		fields []validator.InvalidField
		err    error
	}

	tests := []struct {
		description   string
		user          *models.User
		updateUser    *models.User
		requiredMocks func()
		expected      Expected
	}{
		{
			description: "Fails to find the user by the ID",
			user:        user1,
			updateUser:  &models.User{Name: "name", Email: "user1@email2.com", Username: user2.Username, ID: "id1"},
			requiredMocks: func() {
				mock.On("UserGetByID", ctx, user1.ID, false).Return(nil, 0, Err).Once()
			},
			expected: Expected{nil, Err},
		},
		{
			description: "Fails conflict username",
			user:        &models.User{Name: "name", Email: "user1@email.com", Username: "username1", Password: "hash1", ID: "id1"},
			updateUser:  &models.User{Name: "name", Email: "user1@email2.com", Username: user2.Username, ID: "id1"},
			requiredMocks: func() {
				mock.On("UserGetByID", ctx, user1.ID, false).Return(user1, 0, nil).Once()
				mock.On("UserGetByUsername", ctx, user1.Username).Return(user2, nil).Once()
				mock.On("UserGetByEmail", ctx, user1.Email).Return(user1, nil).Once()
			},
			expected: Expected{[]validator.InvalidField{{"username", "conflict", "", ""}}, Err},
		},
		{
			description: "Fails conflict email and username",
			user:        user1,
			updateUser:  &models.User{Name: "name", Email: "user1@email2.com", Username: user2.Username, ID: "id1"},
			requiredMocks: func() {
				mock.On("UserGetByID", ctx, user1.ID, false).Return(user1, 0, nil).Once()
				mock.On("UserGetByUsername", ctx, user1.Username).Return(user2, nil).Once()
				mock.On("UserGetByEmail", ctx, user1.Email).Return(user2, nil).Once()
			},
			expected: Expected{[]validator.InvalidField{{"username", "conflict", "", ""}, {"email", "conflict", "", ""}}, Err},
		},
		{
			description: "Fails invalid username",
			user:        &models.User{Name: "newname", Email: "user1@email2.com", Username: "invalid_name", ID: "id1"},
			updateUser:  &models.User{Name: "newname", Email: "user1@email2.com", Username: "invalid_name", ID: "id1"},
			requiredMocks: func() {
				mock.On("UserGetByID", ctx, user1.ID, false).Return(user1, 0, nil).Once()
			},
			expected: Expected{[]validator.InvalidField{{"username", "invalid", "alphanum", ""}}, ErrBadRequest},
		},
		{
			description: "Fails invalid email",
			user:        &models.User{Name: "newname", Email: "invalid.email", Username: "newusername", ID: "id1"},
			updateUser:  &models.User{Name: "newname", Email: "invalid.email", Username: "newusername", ID: "id1"},
			requiredMocks: func() {
				mock.On("UserGetByID", ctx, user1.ID, false).Return(user1, 0, nil).Once()
			},
			expected: Expected{[]validator.InvalidField{{"email", "invalid", "email", ""}}, ErrBadRequest},
		},
		{
			description: "Fails invalid email and username",
			user:        &models.User{Name: "newname", Email: "invalid.email", Username: "us", ID: "id1"},
			updateUser:  &models.User{Name: "newname", Email: "invalid.email", Username: "us", ID: "id1"},
			requiredMocks: func() {
				mock.On("UserGetByID", ctx, user1.ID, false).Return(user1, 0, nil).Once()
			},
			expected: Expected{[]validator.InvalidField{{"email", "invalid", "email", ""}, {"username", "invalid", "min", "3"}}, ErrBadRequest},
		},
		{
			description: "Fails empty username",
			user:        &models.User{Name: "", Email: "new@email.com", Username: "newusername", ID: "id1"},
			updateUser:  &models.User{Name: "", Email: "new@email.com", Username: "newusername", ID: "id1"},
			requiredMocks: func() {
				mock.On("UserGetByID", ctx, user1.ID, false).Return(user1, 0, nil).Once()
			},
			expected: Expected{[]validator.InvalidField{{"name", "invalid", "required", ""}}, ErrBadRequest},
		},
		{
			description: "Fails empty email",
			user:        &models.User{Name: "newname", Email: "", Username: "newusername", ID: "id1"},
			updateUser:  &models.User{Name: "newname", Email: "", Username: "newusername", ID: "id1"},
			requiredMocks: func() {
				mock.On("UserGetByID", ctx, user1.ID, false).Return(user1, 0, nil).Once()
			},
			expected: Expected{[]validator.InvalidField{{"email", "invalid", "required", ""}}, ErrBadRequest},
		},
		{
			description: "Successful update user data",
			user:        user1,
			updateUser:  &models.User{Name: "name", Email: "user1@email2.com", Username: user2.Username, ID: "id1"},
			requiredMocks: func() {
				mock.On("UserGetByID", ctx, user1.ID, false).Return(user1, 0, nil).Once()
				mock.On("UserGetByUsername", ctx, user1.Username).Return(nil, Err).Once()
				mock.On("UserGetByEmail", ctx, user1.Email).Return(nil, Err).Once()
				mock.On("UserUpdateData", ctx, user1, user1.ID).Return(nil).Once()
			},
			expected: Expected{nil, nil},
		},
	}

	for _, tc := range tests {
		tc.requiredMocks()
		returnedFields, err := s.UpdateDataUser(ctx, tc.user, tc.updateUser.ID)
		assert.Equal(t, tc.expected, Expected{returnedFields, err})
	}

	mock.AssertExpectations(t)
}

func TestUpdatePasswordUser(t *testing.T) {
	mock := &mocks.Store{}
	s := NewService(store.Store(mock), privateKey, publicKey, storecache.NewNullCache(), clientMock, nil)

	ctx := context.TODO()

	user1 := &models.User{Name: "name", Email: "user1@email.com", Username: "username1", Password: "hash1", ID: "id1"}

	type updatePassword struct {
		currentPassword string
		newPassword     string
		expected        error
	}

	tests := []updatePassword{
		{
			"hiadoshioasc",
			"hashnew",
			ErrUnauthorized,
		},
		{
			"pass123",
			"hashnew",
			ErrUnauthorized,
		},
		{
			"askdhkasd",
			"hashnew",
			ErrUnauthorized,
		},
		{
			"pass890",
			"hashnew",
			ErrUnauthorized,
		},
		{
			"hash1",
			"hashnew",
			nil,
		},
	}

	for _, test := range tests {
		mock.On("UserGetByID", ctx, user1.ID, false).Return(user1, 0, nil).Once()
		if test.expected == nil {
			mock.On("UserUpdatePassword", ctx, test.newPassword, user1.ID).Return(nil).Once()
		}
		err := s.UpdatePasswordUser(ctx, test.currentPassword, test.newPassword, user1.ID)
		assert.Equal(t, err, test.expected)
	}

	mock.AssertExpectations(t)
}
