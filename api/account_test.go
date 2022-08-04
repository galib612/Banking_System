package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/galib612/simplebank/db/mock"
	db "github.com/galib612/simplebank/db/sqlc"
	"github.com/galib612/simplebank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := RandomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Build Stub
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

	//postgres and authdbstore need to be pass
	server := NewServer(store, nil)

	// For testing an HTTP API in Go, we don’t have to start a real HTTP server. Instead,
	// we can just use the recording feature of the httptest package to record the response of the API request
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/get/%d", account.ID)
	// Then we create a new HTTP Request with method GET to that URL.
	// And since it’s a GET request, we can use nil for the request body.
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// Then we call server.router.ServeHTTP() function with the created recorder and request objects.
	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)

}

func RandomAccount() db.Account {
	return db.Account{
		ID:        util.RandomInt(0, 1000),
		Owner:     util.RandomOwner(),
		Balance:   util.RandomMoney(),
		Currency:  util.RandomCurrency(),
		CreatedAt: time.Now(),
	}
}

func TestCreateAccountAPI(t *testing.T) {

	account := RandomAccount()

	arg := db.CreateAccountParams{
		Owner:    account.Owner,
		Balance:  account.Balance,
		Currency: account.Currency,
	}

	jsonData, err := json.Marshal(arg)
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(arg)).Times(1).Return(account, nil)

	server := NewServer(store, nil)
	recorder := httptest.NewRecorder()

	url := "/accounts/create/"

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)

}

func TestListAccountAPI(t *testing.T) {

	account1 := RandomAccount()
	account2 := RandomAccount()

	accounts := []db.Account{account1, account2}

	queryArg := ListAccountRequest{
		PageId:   1,
		PageSize: 2,
	}

	arg := db.ListAccountParams{
		Limit:  queryArg.PageSize,
		Offset: (queryArg.PageId - 1) * queryArg.PageSize,
	}

	jsonData, err := json.Marshal(arg)
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().ListAccount(gomock.Any(), gomock.Eq(arg)).Times(1).Return(accounts, nil)

	url := "/accounts/list/"
	server := NewServer(store, nil)
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(jsonData))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)

}
