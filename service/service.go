package service

import (
	"errors"
	"fmt"
	"goproject_SBG-bot_http/datastruct"
	"time"
)

type repo interface {
	OutList(chatID int64) string
	Get(chatID string) *datastruct.Person
	GetPrevious(chatID int64) string
	GetPersonID() map[int64]*datastruct.Person

	Сancel(chatID int64) error
	AddName(chatID int64) error
	AddNameWork(text_vvod string, chatID int64) error
	DeleteName(chatID int64) (string, error)
	DeleteNameWork(text_vvod string, chatID int64) error
	GetPerson(chatID int64) (*datastruct.Person, error)
	GetPersonName(name string) (*datastruct.Person, error)
	AddPersonName(*datastruct.Person)
	AddPersonID(*datastruct.Person)
	WriteDatabase(*datastruct.Person) error
	UpdateNameDatabase(chatID int64) error
}

type Service struct {
	Repo repo
}

func New(r repo) *Service {

	return &Service{
		Repo: r,
	}
}

func (s *Service) OutList(chatID int64) string {
	return s.Repo.OutList(chatID)
}

func (s *Service) AddName(chatID int64) error {
	return s.Repo.AddName(chatID)
}

func (s *Service) AddNameWork(name string, chatID int64) error {
	err := s.Repo.AddNameWork(name, chatID)
	if err == nil {
		err = s.Repo.UpdateNameDatabase(chatID)
	}
	return err
}

func (s *Service) DeleteName(chatID int64) (string, error) {
	return s.Repo.DeleteName(chatID)
}

func (s *Service) DeleteNameWork(text_vvod string, chatID int64) error {
	err := s.Repo.DeleteNameWork(text_vvod, chatID)
	if err == nil {
		err = s.Repo.UpdateNameDatabase(chatID)
	}
	return err
}

func (s *Service) Сancel(chatID int64) error {
	return s.Repo.Сancel(chatID)
}

func (s *Service) Get(chatID string) *datastruct.Person {
	return s.Repo.Get(chatID)
}

func (s *Service) GetWorker() ([][]int64, []string) {

	to_name := make([]string, 0, 10)
	slice_ID := make([][]int64, 0, 10)
	t_now := time.Now()
	year := t_now.Year()
	//	for _, p := range r.Persons_id {
	for _, p := range s.Repo.GetPersonID() {
		if p.Date != "" {
			_, t_birthday, _ := checkDate(p.Date)
			day := t_birthday.Day()
			month := t_birthday.Month()
			//			fmt.Println("name ", p.name)
			t_person := time.Date(year, time.Month(month), day, 20, 59, 0, 0, time.UTC)
			if t_person.After(t_now) == false {
				t_person = t_person.AddDate(1, 0, 0)
			}
			t_hour := t_person.Sub(t_now).Hours()
			if t_hour < 24 {
				//	            fmt.Println("t_hour", t_hour)
				//				fmt.Println("t_person", t_person.Format("2006-01-02 15:04:05"))
				//				fmt.Println("t_now", t_now.Format("2006-01-02 15:04:05"))
				//				fmt.Println("t_person.After(t_now)", t_person.After(t_now))
				to_name = append(to_name, p.Name)
				slice_ID = append(slice_ID, s.FindList(p.Name))
			}
			//			fmt.Println("to_name", to_name)
			//			fmt.Println("slice_ID", slice_ID)
		}
	}

	return slice_ID, to_name
}
func (s *Service) FindList(name string) []int64 {
	list := make([]int64, 0, 10)
	//	for _, p := range r.Persons_id {
	for _, p := range s.Repo.GetPersonID() {
		if p.Subscribers[name] == 1 {
			list = append(list, p.ID)
		}
	}
	return list
}

func (s *Service) GetPrevious(chatID int64) string {
	return s.Repo.GetPrevious(chatID)
}

func (s *Service) EnterName(name string, chatID int64) error {
	p, _ := s.Repo.GetPerson(chatID)
	if p.Name == "" && p.Previous != "name" {
		p.Previous = "name"
		return errors.New("Авторизация введите имя")
	}

	if p.Previous == "name" {

		if checkName(name) == false {
			return errors.New("Авторизация введите имя")
		}
		_, err := s.Repo.GetPersonName(name)
		if err == nil {
			return errors.New("Данное имя уже существует")
		}
		p.Name = name
		p.Previous = ""
		return nil
	}
	return nil
}

func (s *Service) EnterDate(date string, chatID int64) error {
	//	p := r.Persons_id[chatID]
	p, _ := s.Repo.GetPerson(chatID)
	if p.Date == "" && p.Previous != "year" {
		p.Previous = "year"
		return errors.New("Авторизация введите дату рождения")
	}

	if p.Previous == "year" {
		date, _, err := checkDate(date)
		if err != nil {
			return errors.New("Авторизация введите дату рождения")
		}
		p.Date = date
		p.Previous = ""
		//		r.Persons_name[p.Name] = p
		s.Repo.AddPersonName(p)
		p.Subscribers = map[string]int{}
		err = s.Repo.WriteDatabase(p)
		if err != nil {
			return errors.New("Error")
		}
		return nil
	}
	return nil
}

func checkName(name string) bool {

	for _, v := range name {
		if (v >= 65 && v <= 90) || (v >= 97 && v <= 122) || (v >= 1040 && v <= 1071) || (v >= 1072 && v <= 1103) {
		} else {
			return false
		}
	}
	return true
}

func checkDate(date string) (string, time.Time, error) {
	var t time.Time
	var err error

	const shortForm = "2006-01-02"
	t, err = time.Parse(shortForm, date)
	if err != nil {
		return "", t, err
	}
	return date, t, nil
}

func (s *Service) ChekAvtorisation(chatID int64) bool {
	p, err := s.Repo.GetPerson(chatID)
	if err != nil {
		p = &datastruct.Person{}
		p.ID = chatID
		s.Repo.AddPersonID(p)
		return false
	}

	if p.Date == "" {
		return false
	} else {
		return true
	}
}

//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func (s *Service) GetPersons() []*datastruct.Person {
	fmt.Println("--------------------------------")
	fmt.Println(s.Repo.GetPersonID())
	fmt.Println(len(s.Repo.GetPersonID()))
	fmt.Println("--------------------------------")
	records := make([]*datastruct.Person, len(s.Repo.GetPersonID()))
	i := 0
	for _, p := range s.Repo.GetPersonID() {
		records[i] = p
		i++
	}
	return records
}
func (s *Service) GetPersonByName(name string) (*datastruct.Person, error) {
	return s.Repo.GetPersonName(name)
}
func (s *Service) AddPerson(p *datastruct.Person) error {
	s.Repo.AddPersonID(p)
	s.Repo.AddPersonName(p)
	err := s.Repo.WriteDatabase(p)
	return err
}
func (s *Service) GetBDPersonNameNow(uID int64) []string {
	names := make([]string, 0, 1)
	sliceID, toName := s.GetWorker()

	for i := range sliceID {
		for j := range sliceID[i] {
			if sliceID[i][j] == uID {
				names = append(names, toName[i])
			}
		}
	}
	return names
}
