package api

import (
	//	"goproject_SBG-bot/datastruct"
	//	"reflect"
	"testing"
	//	"time"
	//	"log"
	"errors"
)

type serviceMok struct {
	repos string
	service
}

func TestFind(t *testing.T) {

	a := &Api{&serviceMok{}}
	//s:= &serviceMok{}

	tests := []struct {
		name     string
		arg1     string
		arg2     int64
		res_want string
		want     bool
	}{
		{
			name:     "success1",
			arg1:     "Egor",
			arg2:     5345,
			res_want: "Авторизация завершена" + "\n" + "\n" + Menu(),
			want:     true,
		},
		{
			name:     "success2",
			arg1:     "f2sf4fs6",
			arg2:     4729,
			res_want: "Авторизация\n\nИмя пользователя не введео или введено неверно\nВведите имя пользователя\nили отмена действия /cancel",
			want:     true,
		},
		{
			name:     "success3",
			arg1:     "200-04-32",
			arg2:     1394,
			res_want: "Авторизация\n\nДата рождения пользователя не введена или введена неверно\nВведите дату рождения пользователя в формате:\nГГГГ-ММ-ДД\nнапример:\n2003-10-28\nили отмена действия /cancel",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res_got := a.Autorisation(tt.arg1, tt.arg2)
			var got bool
			if tt.res_want == res_got {
				got = true
			} else {
				got = false
			}
			if got != tt.want {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// ++++++++++++++++ FUNC ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (s *serviceMok) EnterName(name string, chatID int64) error {
	if name == "f2sf4fs6" {
		return errors.New("Ошибка")
	}
	return nil
}

func (s *serviceMok) EnterDate(date string, chatID int64) error {
	if date == "200-04-32" {
		return errors.New("Ошибка")
	}
	return nil
}

//++++++++++++++++ FUNC ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
