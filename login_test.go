package movielite

import (
	"net/http"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/ms900ft/movielite/models"
	log "github.com/sirupsen/logrus"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginSuccess(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Post("/login").
			JSON(`{"username":"test", "password": "test123"}`).
			Expect(t).
			Assert(jsonpath.Matches(`$.token`, `^ey.+$`)).
			Assert(jsonpath.Contains(`$.user_name`, "test")).
			Status(http.StatusOK).
			End()
}

func TestLoginFailed(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Post("/login").
			JSON(`{"username":"test", "password": "test"}`).
			Expect(t).
			Assert(jsonpath.Contains(`$.error`, InvalidCredentials.Error())).
			Status(http.StatusBadRequest).
			End()
}
func TestLoginInvalid(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Post("/login").
			FormData("username", "test").
			FormData("password", "test").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
}
func TestService_FindOne(t *testing.T) {
	//c := Config{Secret: "test123"}
	dbc := models.DBConfig{DBName: ":memory:"}
	db := models.ConnectDataBase(dbc)
	pass, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	user := models.User{UserName: "login", Password: string(pass)}
	if err := db.Create(&user).Error; err != nil {
		log.Fatal(err)
	}
	type fields struct {
		Config *Config
		DB     *gorm.DB
	}
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    token
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				Config: &Config{Secret: "test123"}, DB: db},
			args:    args{username: "login", password: "test123"},
			want:    token{UserName: "login"},
			wantErr: nil,
		},
		{
			name: "wrong password",
			fields: fields{
				Config: &Config{Secret: "test123"}, DB: db},
			args: args{username: "login", password: "test456"},
			//want:    token{UserName: "login"},
			wantErr: InvalidCredentials,
		},
		{
			name: "user not found",
			fields: fields{
				Config: &Config{Secret: "test123"}, DB: db},
			args: args{username: "login123", password: "test456"},
			//want:    token{UserName: "login"},
			wantErr: InvalidCredentials,
		},
		{
			name: "no secret",
			fields: fields{
				Config: &Config{}, DB: db},
			args: args{username: "login", password: "test123"},
			//want:    token{UserName: "login"},
			wantErr: NoSecretFound,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Config: tt.fields.Config,
				DB:     tt.fields.DB,
			}
			got, err := s.FindOne(tt.args.username, tt.args.password)
			if err != nil {
				if err != tt.wantErr {
					t.Errorf("Service.FindOne() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if got.UserName != tt.want.UserName {
				t.Errorf("Service.FindOne() = %v, want %v", got, tt.want)
			}
		})
	}
}
