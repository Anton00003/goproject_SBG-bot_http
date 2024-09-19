package service

import (
	"goproject_SBG-bot/datastruct"
	"reflect"
	"sort"
	"testing"
	"time"
)

type RepositoryMok struct {
	Persons_id map[int64]*datastruct.Person
	repo
}

func TestFind(t *testing.T) {

	// +++++++++++++++++ data for tests ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	time := time.Now().Format("2006-01-02")
	p1 := &datastruct.Person{"Tt", 544223, time, "", map[string]int{"Rr": 1}}
	p2 := &datastruct.Person{"Mm", 544793, "2010-12-30", "", map[string]int{"Tt": 1}}
	p3 := &datastruct.Person{"Rr", 691223, time, "", map[string]int{"Mm": 1, "Tt": 1}}

	r1 := &RepositoryMok{Persons_id: map[int64]*datastruct.Person{544223: p1, 544793: p2, 691223: p3}} //		Persons_name: map[string]*datastruct.Person{"Tt": p1, "Mm": p2, "Rr": p3},

	//	p4 := &datastruct.Person{"Qq", 344183, "2004-01-14", "", map[string]int{}}
	r2 := &RepositoryMok{Persons_id: map[int64]*datastruct.Person{}} //		Persons_name: map[string]*datastruct.Person{},

	p5 := &datastruct.Person{"Qq", 344183, time, "", map[string]int{"Qq": 1}}
	r3 := &RepositoryMok{Persons_id: map[int64]*datastruct.Person{344183: p5}} //		Persons_name: map[string]*datastruct.Person{"Qq": p5},

	// +++++++++++++++++ data for tests ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//	r := &RepositoryMok{}

	tests := []struct {
		name      string
		arg       *RepositoryMok
		res1_want [][]int64
		res2_want []string
		want      bool
	}{
		{
			name:      "success1",
			arg:       r1,
			res1_want: [][]int64{{544793, 691223}, {544223}},
			res2_want: []string{"Tt", "Rr"},
			want:      true,
		},
		{
			name:      "success2",
			arg:       r2,
			res1_want: [][]int64{},
			res2_want: []string{},
			want:      true,
		},
		{
			name:      "success3",
			arg:       r3,
			res1_want: [][]int64{{344183}},
			res2_want: []string{"Qq"},
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//			got, _ := r.Chek_avtorisation(111)
			//			res, _ := ReadFile(r)
			//r := &RepositoryMok{}.Init(read)
			//			r = r.Init(tt.arg)
			res1_got, res2_got := (&Service{tt.arg}).GetWorker()
			var got bool
			got = check_summ(res1_got, res2_got, tt.res1_want, tt.res2_want)
			if got != tt.want {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func check_summ(res1_got [][]int64, res2_got []string, res1_want [][]int64, res2_want []string) bool {
	if check1(res1_got, res1_want) == true && check2(res2_got, res2_want) == true {
		return true
	}
	return false
}
func check1(s1 [][]int64, s2 [][]int64) bool {
	for _, v := range s1 {
		//		sort.Ints(v)
		sort.Slice(v, func(i, j int) bool { return v[i] < v[j] })
	}
	sort.Slice(s1, func(i, j int) bool { return s1[i][0] < s1[j][0] })
	//	fmt.Println(s1)

	for _, v := range s2 {
		//		sort.Ints(v)
		sort.Slice(v, func(i, j int) bool { return v[i] < v[j] })
	}
	sort.Slice(s2, func(i, j int) bool { return s2[i][0] < s2[j][0] })
	//	fmt.Println(s2)

	return reflect.DeepEqual(s1, s2)
}

func check2(s1 []string, s2 []string) bool {
	sort.Strings(s1)
	sort.Strings(s2)
	//	fmt.Println(s1)
	//	fmt.Println(s2)
	return reflect.DeepEqual(s1, s2)
}

func check(res1_got [][]int64, res2_got []string, res1_want [][]int64, res2_want []string) bool {
	if reflect.DeepEqual(res1_got, res1_want) == true && reflect.DeepEqual(res2_got, res2_want) == true {
		return true
	}
	return false
}

// ++++++++++++++++ FUNC ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (r *RepositoryMok) GetPersonID() map[int64]*datastruct.Person {
	return r.Persons_id
}

//++++++++++++++++ FUNC ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
