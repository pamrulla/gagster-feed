package user

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/pamrulla/gagster-feed/helpers"
	"github.com/pamrulla/gagster-feed/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type MyTestObject struct {
	mock.Mock
}

func (o *MyTestObject) GetUsers(db *gorm.DB, User *models.Users) (err error) {
	args := o.Called(db, User)
	return args.Error(0)
}
func (u *MyTestObject) CreateUser(db *gorm.DB, User *models.User) (err error) {
	args := u.Called(db, User)
	return args.Error(0)
}
func (u *MyTestObject) GetUser(db *gorm.DB, User *models.User, id string) (err error) {
	args := u.Called(db, User, id)
	return args.Error(0)
}
func (u *MyTestObject) UpdateUser(db *gorm.DB, User *models.User) (err error) {
	args := u.Called(db, User)
	return args.Error(0)
}
func (u *MyTestObject) DeleteUser(db *gorm.DB, User *models.User, id string) (err error) {
	args := u.Called(db, User, id)
	return args.Error(0)
}
func (u *MyTestObject) EmptyUserTable(db *gorm.DB) {
}

type UserTestSuite struct {
	suite.Suite
	ts      *httptest.Server
	router  *chi.Mux
	ur      *UserRepo
	mockObj MyTestObject
}

func (hts *UserTestSuite) SetupTest() {
	hts.mockObj = MyTestObject{}
	hts.ur = &UserRepo{}
	hts.ur.usr = &hts.mockObj
	hts.router = helpers.CreateNewRouter()
	hts.router.Route("/users", func(r chi.Router) {
		r.Get("/", hts.ur.GetUsers)
		r.Post("/", hts.ur.Create)
		r.Route("/{user_id}", func(r chi.Router) {
			r.Get("/", hts.ur.Get)
			r.Put("/", hts.ur.Update)
			r.Delete("/", hts.ur.Delete)
		})
		r.Put("/enable/{user_id}", hts.ur.Enable)
		r.Put("/disable/{user_id}", hts.ur.Disable)
	})
	hts.ts = httptest.NewServer(hts.router)
}

func (hts *UserTestSuite) TearDownTest() {
	hts.ts.Close()
}

func (hts *UserTestSuite) verifyUsers(exp models.User, act models.User) {
	assert.Equal(hts.T(), exp.First_Name, act.First_Name)
	assert.Equal(hts.T(), exp.Last_Name, act.Last_Name)
	assert.Equal(hts.T(), exp.Email, act.Email)
	assert.Equal(hts.T(), exp.Password, act.Password)
	assert.Equal(hts.T(), exp.IsEnabled, act.IsEnabled)
}
func (hts *UserTestSuite) TestGetAllUsers() {
	// Arrange
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd"}
	hts.mockObj.On("GetUsers", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*models.Users)
		*arg = append(*arg, user)
	})

	// Act
	resp, err := helpers.RunRequest("GET", hts.ts, "/users", nil)

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 200, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	require.Nil(hts.T(), err)

	var actUser models.Users
	err = json.Unmarshal(data, &actUser)
	require.Nil(hts.T(), err)
	require.NotNil(hts.T(), 1, len(actUser))
	hts.verifyUsers(user, actUser[0])
}
func (hts *UserTestSuite) TestGetAllUsers_WhenUnknownErrorOccurred() {
	// Arrange
	hts.mockObj.On("GetUsers", mock.Anything, mock.Anything).Return(gorm.ErrInvaildDB)

	// Act
	resp, err := helpers.RunRequest("GET", hts.ts, "/users", nil)

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 500, resp.StatusCode)
}
func (hts *UserTestSuite) TestGetAUser() {
	// Arrange
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd"}
	hts.mockObj.On("GetUser", mock.Anything, mock.Anything, "1").Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*models.User)
		*arg = user
	})

	// Act
	resp, err := helpers.RunRequest("GET", hts.ts, "/users/1", nil)

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 200, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	require.Nil(hts.T(), err)

	var actUser models.User
	err = json.Unmarshal(data, &actUser)
	require.Nil(hts.T(), err)
	hts.verifyUsers(user, actUser)
}
func (hts *UserTestSuite) TestGetAUser_WhenUnknownErrorOccured() {
	// Arrange
	hts.mockObj.On("GetUser", mock.Anything, mock.Anything, "1").Return(gorm.ErrInvaildDB)

	// Act
	resp, err := helpers.RunRequest("GET", hts.ts, "/users/1", nil)

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 500, resp.StatusCode)
}
func (hts *UserTestSuite) TestGetAUser_WhenUserNotFound() {
	// Arrange
	hts.mockObj.On("GetUser", mock.Anything, mock.Anything, "1").Return(gorm.ErrRecordNotFound)

	// Act
	resp, err := helpers.RunRequest("GET", hts.ts, "/users/1", nil)

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 404, resp.StatusCode)
}
func (hts *UserTestSuite) TestCreateAUser() {
	// Arrange
	var recievedUser models.User
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd"}
	hts.mockObj.On("CreateUser", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*models.User)
		recievedUser = *arg
	})
	payload, _ := json.Marshal(user)

	// Act
	resp, err := helpers.RunRequest("POST", hts.ts, "/users", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 200, resp.StatusCode)
	hts.verifyUsers(user, recievedUser)
}
func (hts *UserTestSuite) TestCreateAUser_WhenUnknownErrorOccured() {
	// Arrange
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd"}
	hts.mockObj.On("CreateUser", mock.Anything, mock.Anything).Return(gorm.ErrInvaildDB)
	payload, _ := json.Marshal(user)

	// Act
	resp, err := helpers.RunRequest("POST", hts.ts, "/users", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 500, resp.StatusCode)
}
func (hts *UserTestSuite) TestCreateAUser_WhenInvalidDataSent() {
	// Arrange
	payload, _ := json.Marshal([]byte(`{"name":what?}`))

	// Act
	resp, err := helpers.RunRequest("POST", hts.ts, "/users", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), http.StatusBadRequest, resp.StatusCode)
}
func (hts *UserTestSuite) TestCreateAUser_WhenNoDataSent() {
	// Act
	resp, err := helpers.RunRequest("POST", hts.ts, "/users", nil)

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), http.StatusBadRequest, resp.StatusCode)
}
func (hts *UserTestSuite) TestUpdateAUser() {
	// Arrange
	var recievedUser models.User
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd"}
	hts.mockObj.On("UpdateUser", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*models.User)
		recievedUser = *arg
	})
	payload, _ := json.Marshal(user)

	// Act
	resp, err := helpers.RunRequest("PUT", hts.ts, "/users/1", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 200, resp.StatusCode)
	hts.verifyUsers(user, recievedUser)
}
func (hts *UserTestSuite) TestUpdateAUser_WhenUnknownErrorOccured() {
	// Arrange
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd"}
	hts.mockObj.On("UpdateUser", mock.Anything, mock.Anything).Return(gorm.ErrInvaildDB)
	payload, _ := json.Marshal(user)

	// Act
	resp, err := helpers.RunRequest("PUT", hts.ts, "/users/1", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 500, resp.StatusCode)
}
func (hts *UserTestSuite) TestUpdateAUser_WhenUserNotFound() {
	// Arrange
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd"}
	hts.mockObj.On("UpdateUser", mock.Anything, mock.Anything).Return(gorm.ErrRecordNotFound)
	payload, _ := json.Marshal(user)

	// Act
	resp, err := helpers.RunRequest("PUT", hts.ts, "/users/1", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 404, resp.StatusCode)
}
func (hts *UserTestSuite) TestUpdateAUser_WhenInvalidDataSent() {
	// Arrange
	payload, _ := json.Marshal([]byte(`{"name":what?}`))

	// Act
	resp, err := helpers.RunRequest("PUT", hts.ts, "/users/1", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), http.StatusBadRequest, resp.StatusCode)
}
func (hts *UserTestSuite) TestUpdateAUser_WhenNoDataSent() {
	// Act
	resp, err := helpers.RunRequest("PUT", hts.ts, "/users/1", nil)

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), http.StatusBadRequest, resp.StatusCode)
}
func (hts *UserTestSuite) TestDeleteAUser() {
	// Arrange
	user_id := "1"
	hts.mockObj.On("DeleteUser", mock.Anything, mock.Anything, user_id).Return(nil)

	// Act
	resp, err := helpers.RunRequest("DELETE", hts.ts, "/users/"+user_id, nil)

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 200, resp.StatusCode)
}
func (hts *UserTestSuite) TestDeleteAUser_WhenUnknownErrorOccurs() {
	// Arrange
	user_id := "1"
	hts.mockObj.On("DeleteUser", mock.Anything, mock.Anything, user_id).Return(gorm.ErrInvaildDB)

	// Act
	resp, err := helpers.RunRequest("DELETE", hts.ts, "/users/"+user_id, nil)

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 500, resp.StatusCode)
}
func (hts *UserTestSuite) TestDeleteAUser_WhenUserNotFound() {
	// Arrange
	user_id := "1"
	hts.mockObj.On("DeleteUser", mock.Anything, mock.Anything, user_id).Return(gorm.ErrRecordNotFound)

	// Act
	resp, err := helpers.RunRequest("DELETE", hts.ts, "/users/"+user_id, nil)

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 404, resp.StatusCode)
}
func (hts *UserTestSuite) TestEnableAUser() {
	// Arrange
	var recievedUser models.User
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd", IsEnabled: false}
	hts.mockObj.On("GetUser", mock.Anything, mock.Anything, "1").Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*models.User)
		*arg = user
	})
	hts.mockObj.On("UpdateUser", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*models.User)
		recievedUser = *arg
	})
	payload, _ := json.Marshal(user)

	// Act
	resp, err := helpers.RunRequest("PUT", hts.ts, "/users/enable/1", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 200, resp.StatusCode)
	assert.True(hts.T(), recievedUser.IsEnabled)
}
func (hts *UserTestSuite) TestEnableAUser_UserNotFoundWhileUpdatingUser() {
	// Arrange
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd", IsEnabled: true}
	hts.mockObj.On("GetUser", mock.Anything, mock.Anything, "1").Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*models.User)
		*arg = user
	})
	hts.mockObj.On("UpdateUser", mock.Anything, mock.Anything).Return(gorm.ErrRecordNotFound)
	payload, _ := json.Marshal(user)

	// Act
	resp, err := helpers.RunRequest("PUT", hts.ts, "/users/enable/1", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 404, resp.StatusCode)
}
func (hts *UserTestSuite) TestEnableAUser_FailedToUpdateUser() {
	// Arrange
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd", IsEnabled: true}
	hts.mockObj.On("GetUser", mock.Anything, mock.Anything, "1").Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*models.User)
		*arg = user
	})
	hts.mockObj.On("UpdateUser", mock.Anything, mock.Anything).Return(gorm.ErrNotImplemented)
	payload, _ := json.Marshal(user)

	// Act
	resp, err := helpers.RunRequest("PUT", hts.ts, "/users/enable/1", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 500, resp.StatusCode)
}
func (hts *UserTestSuite) TestEnableAUser_FailedToGetUser() {
	// Arrange
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd", IsEnabled: true}
	hts.mockObj.On("GetUser", mock.Anything, mock.Anything, "1").Return(gorm.ErrNotImplemented)
	payload, _ := json.Marshal(user)

	// Act
	resp, err := helpers.RunRequest("PUT", hts.ts, "/users/enable/1", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 500, resp.StatusCode)
}
func (hts *UserTestSuite) TestEnableAUser_UserNotFoundToGetUser() {
	// Arrange
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd", IsEnabled: true}
	hts.mockObj.On("GetUser", mock.Anything, mock.Anything, "1").Return(gorm.ErrRecordNotFound)
	payload, _ := json.Marshal(user)

	// Act
	resp, err := helpers.RunRequest("PUT", hts.ts, "/users/enable/1", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 404, resp.StatusCode)
}
func (hts *UserTestSuite) TestDisableAUser() {
	// Arrange
	var recievedUser models.User
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd", IsEnabled: true}
	hts.mockObj.On("GetUser", mock.Anything, mock.Anything, "1").Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*models.User)
		*arg = user
	})
	hts.mockObj.On("UpdateUser", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*models.User)
		recievedUser = *arg
	})
	payload, _ := json.Marshal(user)

	// Act
	resp, err := helpers.RunRequest("PUT", hts.ts, "/users/disable/1", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 200, resp.StatusCode)
	assert.False(hts.T(), recievedUser.IsEnabled)
}
func (hts *UserTestSuite) TestDisableAUser_UserNotFoundWhileUpdatingUser() {
	// Arrange
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd", IsEnabled: true}
	hts.mockObj.On("GetUser", mock.Anything, mock.Anything, "1").Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*models.User)
		*arg = user
	})
	hts.mockObj.On("UpdateUser", mock.Anything, mock.Anything).Return(gorm.ErrRecordNotFound)
	payload, _ := json.Marshal(user)

	// Act
	resp, err := helpers.RunRequest("PUT", hts.ts, "/users/disable/1", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 404, resp.StatusCode)
}
func (hts *UserTestSuite) TestDisableAUser_FailedToUpdateUser() {
	// Arrange
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd", IsEnabled: true}
	hts.mockObj.On("GetUser", mock.Anything, mock.Anything, "1").Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*models.User)
		*arg = user
	})
	hts.mockObj.On("UpdateUser", mock.Anything, mock.Anything).Return(gorm.ErrNotImplemented)
	payload, _ := json.Marshal(user)

	// Act
	resp, err := helpers.RunRequest("PUT", hts.ts, "/users/disable/1", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 500, resp.StatusCode)
}
func (hts *UserTestSuite) TestDisableAUser_FailedToGetUser() {
	// Arrange
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd", IsEnabled: true}
	hts.mockObj.On("GetUser", mock.Anything, mock.Anything, "1").Return(gorm.ErrNotImplemented)
	payload, _ := json.Marshal(user)

	// Act
	resp, err := helpers.RunRequest("PUT", hts.ts, "/users/disable/1", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 500, resp.StatusCode)
}
func (hts *UserTestSuite) TestDisableAUser_UserNotFoundToGetUser() {
	// Arrange
	user := models.User{First_Name: "test first name", Last_Name: "test last name", Email: "testemail@ag.com", Password: "asfd", IsEnabled: true}
	hts.mockObj.On("GetUser", mock.Anything, mock.Anything, "1").Return(gorm.ErrRecordNotFound)
	payload, _ := json.Marshal(user)

	// Act
	resp, err := helpers.RunRequest("PUT", hts.ts, "/users/disable/1", bytes.NewReader(payload))

	// Assert
	require.Nil(hts.T(), err)
	hts.mockObj.AssertExpectations(hts.T())
	assert.Equal(hts.T(), 404, resp.StatusCode)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
