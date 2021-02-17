package token

import (
	"context"
	"testing"

	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/api/store/mocks"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/stretchr/testify/assert"
)

type Data struct {
	Mock      *mocks.Store
	Service   Service
	Context   context.Context
	Namespace *models.Namespace
	Token     *models.Token
	err       error
}

func TestCreateToken(t *testing.T) {
	data := initData(t)

	assert.NoError(t, data.err)
}

func TestListToken(t *testing.T) {
	data := initData(t)

	mock.On("TokenCreateAPIToken", ctx, namespace.TenantID).Return(&token, nil).Once()

	createdToken, err := svc.CreateToken(ctx, namespace.TenantID)
	assert.NoError(t, err)

	Err := errors.New("error")

	type Expected struct {
		userToken *models.Token
		err       error
	}

	tests := []struct {
		description   string
		args          models.Namespace
		requiredMocks func()
		expected      Expected
	}{
		{
			description: "Fails the namespace not found",
			args:        data.Namespace,
			requiredMocks: func() {
				data.Mock.On("TokenListAPIToken", data.Context, data.Namespace.TenantID).Return([]models.Token{}, nil).Once()
			},
		},
		expected: Expected{nil, Err},
		{
			description: "Fails no API Token stored",
			args:        data.Namespace,
			requiredMocks: func() {
				data.Mock.On("TokenListAPIToken", data.Context, data.Namespace.TenantID).Return([]models.Token{}, nil).Once()
			},
		},
		expected: Expected{nil, Err},
		{
			description: "Successful list the tokens",
			args:        data.Namespace,
			requiredMocks: func() {
				data.Mock.On("TokenListAPIToken", data.Context, data.Namespace.TenantID).Return([]models.Token{}, nil).Once()
			},
		},
		expected: Expected{createdToken, nil},
	}

	for _, test := range tests {
		t.Log("PASS:  ", test.description)
		test.requiredMocks()
		apiToken, err := data.Service.ListToken(ctx, test.args.TenantID)
		assert.Equal(t, test.expected, Expected{apiToken, err})
	}

	mock.AssertExpectations(t)
}

func TestGetToken(t *testing.T) {
	data := initData(t)

	mock.On("TokenCreateAPIToken", ctx, namespace.TenantID).Return(&token, nil).Once()

	createdToken, err := svc.CreateToken(ctx, namespace.TenantID)
	assert.NoError(t, err)

	Err := errors.New("error")

	type Expected struct {
		userToken *models.Token
		err       error
	}

	tests := []struct {
		description   string
		args          models.Namespace
		requiredMocks func()
		expected      Expected
	}{
		{
			description: "Fails the namespace not found",
			args:        data.Namespace,
			requiredMocks: func() {
				data.Mock.On("TokenGetAPIToken", data.Context, data.Namespace.TenantID, data.Token.ID).Return(&models.Token{}, nil).Once()
			},
		},
		expected: Expected{nil, Err},
		{
			description: "Fails API Token ID invalid",
			args:        data.Namespace,
			requiredMocks: func() {
				data.Mock.On("TokenGetAPIToken", data.Context, data.Namespace.TenantID, data.Token.ID).Return(&models.Token{}, nil).Once()
			},
		},
		expected: Expected{nil, Err},
		{
			description: "Successful get the API token",
			args:        data.Namespace,
			requiredMocks: func() {
				data.Mock.On("TokenGetAPIToken", data.Context, data.Namespace.TenantID, data.Token.ID).Return(&models.Token{}, nil).Once()
			},
		},
		expected: Expected{createdToken, nil},
	}

	for _, test := range tests {
		t.Log("PASS:  ", test.description)
		test.requiredMocks()
		apiToken, err := data.Service.GetToken(ctx, test.args.TenantID, data.Token.ID)
		assert.Equal(t, test.expected, Expected{apiToken, err})
	}

	mock.AssertExpectations(t)
}

func TestDeleteToken(t *testing.T) {
	data := initData(t)

	mock.On("TokenCreateAPIToken", ctx, namespace.TenantID).Return(&token, nil).Once()

	createdToken, err := svc.CreateToken(ctx, namespace.TenantID)
	assert.NoError(t, err)

	Err := errors.New("error")

	type Expected struct {
		userToken *models.Token
		err       error
	}

	tests := []struct {
		description   string
		namespace     models.Namespace
		token         models.Token
		requiredMocks func()
		expected      Expected
	}{
		{
			description: "Fails the namespace not found",
			namespace:   data.Namespace,
			token:       data.Token,
			requiredMocks: func() {
				data.Mock.On("TokenDeleteAPIToken", data.Context, namespace.TenantID, token.ID).Return(nil).Once()
			},
		},
		expected: Expected{nil, Err},
		{
			description: "Fails API Token ID invalid",
			namespace:   data.Namespace,
			token:       data.Token,
			requiredMocks: func() {
				data.Mock.On("TokenDeleteAPIToken", data.Context, namespace.TenantID, token.ID).Return(nil).Once()
			},
		},
		expected: Expected{nil, Err},
		{
			description: "Successful delete the API token",
			namespace:   data.Namespace,
			token:       data.Token,
			requiredMocks: func() {
				data.Mock.On("TokenDeleteAPIToken", data.Context, namespace.TenantID, token.ID).Return(nil).Once()
			},
		},
		expected: Expected{createdToken, nil},
	}

	for _, test := range tests {
		t.Log("PASS:  ", test.description)
		test.requiredMocks()
		apiToken, err := data.Service.DeleteToken(ctx, test.namespace.TenantID, data.token.ID)
		assert.Equal(t, test.expected, Expected{apiToken, err})
	}

	mock.AssertExpectations(t)
}

func TestUpdateToken(t *testing.T) {
	data := initData(t)

	req := &models.APITokenUpdate{
		TokenFields: models.TokenFields{ReadOnly: false},
	}

	mock.On("TokenCreateAPIToken", ctx, namespace.TenantID).Return(&token, nil).Once()

	createdToken, err := svc.CreateToken(ctx, namespace.TenantID)
	assert.NoError(t, err)

	Err := errors.New("error")

	type Expected struct {
		userToken *models.Token
		err       error
	}

	tests := []struct {
		description   string
		namespace     models.Namespace
		token         models.Token
		requiredMocks func()
		expected      Expected
	}{
		{
			description: "Fails the namespace not found",
			namespace:   data.Namespace,
			token:       data.Token,
			requiredMocks: func() {
				data.Mock.On("TokenUpdateAPIToken", data.Context, data.Namespace.TenantID, data.Token.ID, req).Return(nil, nil).Once()
				data.Mock.On("TokenGetAPIToken", data.Context, data.Namespace.TenantID, data.Token.ID).Return(&models.Token{}, nil).Once()
			},
		},
		expected: Expected{nil, Err},
		{
			description: "Fails API Token ID invalid",
			namespace:   data.Namespace,
			token:       data.Token,
			requiredMocks: func() {
				data.Mock.On("TokenUpdateAPIToken", data.Context, data.Namespace.TenantID, data.Token.ID, req).Return(nil, nil).Once()
				data.Mock.On("TokenGetAPIToken", data.Context, data.Namespace.TenantID, data.Token.ID).Return(&models.Token{}, nil).Once()
			},
		},
		expected: Expected{nil, Err},
		{
			description: "Successful delete the API token",
			namespace:   data.Namespace,
			token:       data.Token,
			requiredMocks: func() {
				data.Mock.On("TokenUpdateAPIToken", data.Context, data.Namespace.TenantID, data.Token.ID, req).Return(nil, nil).Once()
				data.Mock.On("TokenGetAPIToken", data.Context, data.Namespace.TenantID, data.Token.ID).Return(&models.Token{}, nil).Once()
			},
		},
		expected: Expected{createdToken, nil},
	}

	for _, test := range tests {
		t.Log("PASS:  ", test.description)
		test.requiredMocks()
		apiToken, err := data.Service.DeleteToken(ctx, test.namespace.TenantID, data.token.ID)
		assert.Equal(t, test.expected, Expected{apiToken, err})
	}
}

func initData(t *testing.T) (response Data) {
	t.Helper()
	mock := &mocks.Store{}

	ctx := context.TODO()

	svc := NewService(store.Store(mock))

	token := models.Token{
		ID:       "1",
		TenantID: "a736a52b-5777-4f92-b0b8-e359bf484713",
		ReadOnly: true,
	}

	namespace := &models.Namespace{Name: "group1", Owner: "hash1", TenantID: "a736a52b-5777-4f92-b0b8-e359bf484713", APITokens: []models.Token{}}

	mock.On("TokenCreateAPIToken", ctx, namespace.TenantID).Return(&token, nil).Once()

	createdToken, err := svc.CreateToken(ctx, namespace.TenantID)
	assert.NoError(t, err)

	return Data{
		Mock:      mock,
		Service:   svc,
		Context:   ctx,
		Namespace: namespace,
		Token:     createdToken,
		err:       nil,
	}
}
