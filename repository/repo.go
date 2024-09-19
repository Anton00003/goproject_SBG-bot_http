package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"goproject_SBG-bot_http/datastruct"
	"strconv"
)

const (
// file = "tmp/Birthbay.csv"
)

//type read interface {
//	ReadFile_2() ([][]string, error)
//	WriteFile_2([][]string) error
//	WriteFile([][]string) error
//}

type databaseRead interface {
	ReadDatabase() ([][]string, [][]byte, error)
	WriteDatabase([]string, []byte) error
	UpdateDatabase(nameUpdate string, record []byte) error
}

type Repository struct {
	Persons_id   map[int64]*datastruct.Person
	Persons_name map[string]*datastruct.Person
	Reader       databaseRead
	cnt          int
}

func New(reader databaseRead) *Repository {
	return (&Repository{}).Init(reader)
}

func (r *Repository) Init(reader databaseRead) *Repository {
	r.Reader = reader
	//record, err := reader.ReadDatabase()
	record, b, err := reader.ReadDatabase()
	if err != nil {
		fmt.Println(err) // может return nil, err
	}

	map_person_n := map[string]*datastruct.Person{}

	for i, v := range record {
		if len(v) > 1 {
			p := &datastruct.Person{}
			p.Name = v[0]
			p.ID, _ = strconv.ParseInt(v[1], 10, 64)
			//			p.Date, _, _ = checkDate(v[2])
			p.Date = v[2]

			m := map[string]int{}
			err = json.Unmarshal(b[i], &m)
			if err != nil {
				fmt.Println(err)
			}
			p.Subscribers = m
			map_person_n[p.Name] = p
		}
	}
	fmt.Println(map_person_n)
	//++++++++++++++++++++++++++++++++++++++++++++++++++

	map_person_id := map[int64]*datastruct.Person{}
	for _, v := range map_person_n {
		if v.ID != 0 {
			map_person_id[v.ID] = v
			fmt.Println(*v)
		}
	}
	//	fmt.Println(map_person_n["Gdgrd"])
	r.Persons_id = map_person_id
	r.Persons_name = map_person_n

	return r
}

func (r *Repository) GetPerson(chatID int64) (*datastruct.Person, error) {
	p, ok := r.Persons_id[chatID]
	if !ok {
		return nil, errors.New("Error")
	}
	return p, nil
}
func (r *Repository) GetPersonName(name string) (*datastruct.Person, error) {
	p, ok := r.Persons_name[name]
	if !ok {
		return nil, errors.New("Error")
	}
	return p, nil
}

func (r *Repository) AddPersonName(p *datastruct.Person) {
	r.Persons_name[p.Name] = p
}

func (r *Repository) AddPersonID(p *datastruct.Person) {
	r.Persons_id[p.ID] = p
}

func (r *Repository) OutList(chatID int64) string {

	textout := "Зарегестрированы следующие сотрудники:"
	for n := range r.Persons_name {
		textout = textout + "\n" + n
	}

	return textout
}

func (r *Repository) Get(chatID string) *datastruct.Person {
	return r.Persons_name[chatID]
}

func (r *Repository) GetPrevious(chatID int64) string {
	return r.Persons_id[chatID].Previous
}

func (r *Repository) AddName(chatID int64) error {

	_, ok := r.Persons_id[chatID]
	if ok {
		r.Persons_id[chatID].Previous = "Add"
		return nil

	} else {
		return errors.New("Error")
	}

	return nil

}

func (r *Repository) AddNameWork(text_vvod string, chatID int64) error {

	_, q := r.Persons_name[text_vvod]
	if q == true {
		if r.Persons_id[chatID].Subscribers[text_vvod] == 1 {
			return errors.New("Error")
		}
		r.Persons_id[chatID].Subscribers[text_vvod] = 1
		r.Persons_id[chatID].Previous = ""
		return nil
	} else {
		return errors.New("Error")
	}

	return nil
}

func (r *Repository) DeleteName(chatID int64) (string, error) {
	var listName string

	p := r.Persons_id[chatID]
	//	n := 0
	for p_name := range p.Subscribers {
		//		n++
		listName = listName + "\n" + p_name
	}

	if listName == "" {
		return listName, errors.New("Error")
	}
	p.Previous = "Delete"

	return listName, nil
}

func (r *Repository) DeleteNameWork(text_vvod string, chatID int64) error {

	p := r.Persons_id[chatID]
	_, q := r.Persons_name[text_vvod]
	if q == false {
		return errors.New("Error")
	} else {
		_, ok := p.Subscribers[text_vvod]
		if ok {
			delete(p.Subscribers, text_vvod)
			p.Previous = ""
			return nil

		} else {
			return errors.New("Error")
		}
	}

	return nil
}

func (r *Repository) Сancel(chatID int64) error {
	//	var textout string
	_, ok := r.Persons_id[chatID]
	if ok {
		r.Persons_id[chatID].Previous = ""
		//		textout = "Отменено"
		return nil
	}
	return errors.New("Error")
}

func (r *Repository) GetPersonID() map[int64]*datastruct.Person {
	return r.Persons_id
}

//func ReadFile(reader read) ([][]string, error) {
//
//	return (reader).ReadFile_2()
//}

func ReadDatabase(reader databaseRead) ([][]string, [][]byte, error) {

	return reader.ReadDatabase()
}

func (r *Repository) WriteDatabase(p *datastruct.Person) error {

	record := make([]string, 3)
	record[0] = p.Name
	record[1] = strconv.Itoa(int(p.ID))
	record[2] = p.Date
	//		record[3] = p.Previous

	b, err := json.Marshal(p.Subscribers)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = r.Reader.WriteDatabase(record, b)

	return err
}

func (r *Repository) UpdateNameDatabase(chatID int64) error {
	p := r.Persons_id[chatID]
	record, err := json.Marshal(p.Subscribers)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = r.Reader.UpdateDatabase(p.Name, record)
	if err != nil {
		return err
	}

	return nil
}
