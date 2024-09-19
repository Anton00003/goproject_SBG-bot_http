package repository

import (
	"goproject_SBG-bot/datastruct"
	"reflect"
	"testing"
)

type Reader struct {
	filePath string
	databaseRead
}

func TestFind(t *testing.T) {

	// +++++++++++++++++ data for tests ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	p1 := &datastruct.Person{"Tt", 544223, "2000-10-24", "", map[string]int{}}
	p2 := &datastruct.Person{"Mm", 544793, "2010-12-30", "", map[string]int{"Tt": 1}}
	p3 := &datastruct.Person{"Rr", 691223, "1980-06-05", "", map[string]int{"Mm": 1, "Tt": 1}}

	r1 := &Repository{Persons_id: map[int64]*datastruct.Person{544223: p1, 544793: p2, 691223: p3},
		Persons_name: map[string]*datastruct.Person{"Tt": p1, "Mm": p2, "Rr": p3},
		Reader:       &Reader{filePath: "WWW"},
	}
	//	p4 := &datastruct.Person{"Qq", 344183, "2004-01-14", "", map[string]int{}}
	r2 := &Repository{Persons_id: map[int64]*datastruct.Person{},
		Persons_name: map[string]*datastruct.Person{},
		Reader:       &Reader{filePath: "RRR"},
	}

	p5 := &datastruct.Person{"Qq", 344183, "2004-01-14", "", map[string]int{}}
	r3 := &Repository{Persons_id: map[int64]*datastruct.Person{344183: p5},
		Persons_name: map[string]*datastruct.Person{"Qq": p5},
		Reader:       &Reader{filePath: "YYY"},
	}
	// +++++++++++++++++ data for tests ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	r := &Repository{}

	tests := []struct {
		name string
		arg  string
		res  *Repository
		want bool
	}{
		{
			name: "success",
			arg:  "WWW",
			res:  r1,
			want: true,
		},
		{
			name: "error1",
			arg:  "RRR",
			res:  r2,
			want: true,
		},
		{
			name: "error2",
			arg:  "YYY",
			res:  r3,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//			got, _ := r.Chek_avtorisation(111)
			//			res, _ := ReadFile(r)
			//r := &Repository{}.Init(read)
			r = r.Init(&Reader{filePath: tt.arg})

			var got bool
			got = reflect.DeepEqual(r, tt.res)
			if got != tt.want {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func (read *Reader) ReadDatabase() ([][]string, [][]byte, error) {

	if read.filePath == "WWW" {
		s := [][]string{{"Tt", "544223", "2000-10-24"}, {"Mm", "544793", "2010-12-30"}, {"Rr", "691223", "1980-06-05"}}
		return s, [][]byte{{123, 125}, {123, 34, 84, 116, 34, 58, 49, 125}, {123, 34, 77, 109, 34, 58, 49, 44, 34, 84, 116, 34, 58, 49, 125}}, nil
	}

	if read.filePath == "RRR" {
		s := [][]string{}
		return s, [][]byte{}, nil
	}

	if read.filePath == "YYY" {
		s := [][]string{{"Qq", "344183", "2004-01-14"}}
		return s, [][]byte{{123, 125}}, nil
	}
	return [][]string{}, [][]byte{}, nil
}
