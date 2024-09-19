package api_http

import (
	//	"encoding/csv"
	//	"database/sql"
	//	"encoding/json"
	//	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"goproject_SBG-bot_http/datastruct"

	_ "github.com/lib/pq"
)

type User struct {
	Person        *datastruct.Person
	BDSubscribers []string
}

type Api struct {
	Service service
}

type service interface {
	ChekAvtorisation(chatID int64) bool
	OutList(chatID int64) string
	AddName(chatID int64) error
	AddNameWork(text_vvod string, chatID int64) error
	DeleteName(chatID int64) (string, error)
	DeleteNameWork(text_vvod string, chatID int64) error
	Сancel(chatID int64) error
	GetWorker() ([][]int64, []string)
	GetPrevious(chatID int64) string
	EnterName(text string, chatID int64) error
	EnterDate(text string, chatID int64) error
	//+++++++++++++++++
	GetPersons() []*datastruct.Person
	GetPersonByName(string) (*datastruct.Person, error)
	AddPerson(*datastruct.Person) error
	GetBDPersonNameNow(int64) []string
}

func New(s service) *Api {
	return &Api{Service: s}
}

func (a *Api) Run() {
	a.handleRequest()
}
func (a *Api) handleRequest() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/", a.main_page).Methods("GET")
	rtr.HandleFunc("/contacts/", a.contacts_page).Methods("GET")
	rtr.HandleFunc("/enter/", a.enter_page).Methods("GET", "POST")
	rtr.HandleFunc("/registration/", a.registration_page).Methods("GET", "POST")
	rtr.HandleFunc("/enter_user/", a.enter_user)
	rtr.HandleFunc("/save_user/", a.save_user)
	rtr.HandleFunc("/users/{name}/subscribe/", a.subscribe)
	rtr.HandleFunc("/users/{name}/unsubscribe/", a.unsubscribe)

	rtr.HandleFunc("/users/{name}/", a.user_page).Methods("GET", "POST")
	rtr.HandleFunc("/users/{name}/subscribers/", a.user_subscribers_page).Methods("GET", "POST")
	rtr.HandleFunc("/error/{err}/", a.error_page).Methods("GET", "POST")

	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil) //        http://localhost:8080/
}

func (a *Api) main_page(w http.ResponseWriter, r *http.Request) {
	//	bob := User{"Bob", 25, []string{"football", "skate"}}
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "index", nil)
}
func (a *Api) error_page(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars["err"])
	//	user := GetUserByName(vars["name"])

	//fmt.Fprintf(w, "ID: %v\n", user.Name)

	t, err := template.ParseFiles("templates/error.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "error", vars["err"])
}
func (a *Api) contacts_page(w http.ResponseWriter, r *http.Request) {
	users := a.Service.GetPersons()
	fmt.Println("++++++++++++++++++++++++++++++")
	fmt.Println(users)
	fmt.Println("++++++++++++++++++++++++++++++")

	t, err := template.ParseFiles("templates/contacts.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "contacts", users)
}
func (a *Api) enter_page(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/enter.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "enter", nil)
}
func (a *Api) registration_page(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/reg.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "reg", nil)
}
func (a *Api) user_page(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars["name"])
	user, _ := a.Service.GetPersonByName(vars["name"])

	//fmt.Fprintf(w, "ID: %v\n", user.Name)

	t, err := template.ParseFiles("templates/user.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	u := User{user, a.Service.GetBDPersonNameNow(user.ID)}

	//	t.ExecuteTemplate(w, "user", vars["name"])
	t.ExecuteTemplate(w, "user", u)
}
func (a *Api) user_subscribers_page(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, _ := a.Service.GetPersonByName(vars["name"])
	t, err := template.ParseFiles("templates/subscribers.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "subscribers", user)
}
func (a *Api) enter_user(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	birthday := r.FormValue("birthday")

	fmt.Println(name)
	fmt.Println(birthday)

	//	b, _ := json.Marshal(map[string]int{})
	//	err := WriteDatabase(name, birthday, b)
	u, err := a.Service.GetPersonByName(name)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/error/error_name_notexist/", http.StatusSeeOther)
		return
	}
	if u.Date != birthday {
		http.Redirect(w, r, "/error/error_date_notexist/", http.StatusSeeOther)
		return
	}

	// http.Redirect(w, r, "/users/{{name}}/", http.StatusSeeOther)
	//	fmt.Println(fmt.Sprintf("/users/%s/", name))
	http.Redirect(w, r, fmt.Sprintf("/users/%s/", name), http.StatusSeeOther)
}
func (a *Api) save_user(w http.ResponseWriter, r *http.Request) {
	var idMax int64
	name := r.FormValue("name")
	birthday := r.FormValue("birthday")

	fmt.Println(name)
	fmt.Println(birthday)

	_, err_u := a.Service.GetPersonByName(name)
	if err_u == nil {
		fmt.Println("Пользователь с таким именем существует")
		http.Redirect(w, r, "/error/error_name_exist/", http.StatusSeeOther)
		return
	}
	idMax = 1
	users := a.Service.GetPersons()
	for _, p := range users {
		if p.ID < 1000 && p.ID >= idMax {
			idMax = p.ID + 1
		}
	}
	p := &datastruct.Person{Name: name, ID: idMax, Date: birthday, Subscribers: map[string]int{}}
	err := a.Service.AddPerson(p)
	//	b, _ := json.Marshal(map[string]int{})
	//	err := WriteDatabase(name, birthday, b)
	if err != nil {
		fmt.Println(err)
	}

	// http.Redirect(w, r, "/users/{{name}}/", http.StatusSeeOther)
	//	fmt.Println(fmt.Sprintf("/users/%s/", name))
	http.Redirect(w, r, fmt.Sprintf("/users/%s/", name), http.StatusSeeOther)
}
func (a *Api) subscribe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name_u := vars["name"]
	name := r.FormValue("name")

	fmt.Println("name ", name)
	fmt.Println("name_u ", name_u)

	// u, _ := GetUserByName(name_u)
	u, _ := a.Service.GetPersonByName(name_u)
	_, err := a.Service.GetPersonByName(name)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/error/error_name_subs_notexist/", http.StatusSeeOther)
		return
	}
	err = a.Service.AddNameWork(name, u.ID)
	if err != nil {
		fmt.Println("Ошибка добавления на кого подписан")
		fmt.Println(err)
	}
	http.Redirect(w, r, fmt.Sprintf("/users/%s/subscribers/", name_u), http.StatusSeeOther)
}
func (a *Api) unsubscribe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name_u := vars["name"]
	name := r.FormValue("name")

	fmt.Println("name ", name)
	fmt.Println("name_u ", name_u)

	u, _ := a.Service.GetPersonByName(name_u)
	_, err := a.Service.GetPersonByName(name)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/error/error_name_subs_notexist/", http.StatusSeeOther)
		return
	}
	err = a.Service.DeleteNameWork(name, u.ID)
	//	delete(u.Subscribers, name)
	//	b, _ := json.Marshal(u.Subscribers)
	//	err = UpdateDatabase(name_u, b)
	if err != nil {
		fmt.Println("Ошибка удаления на кого подписан")
		fmt.Println(err)
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%s/subscribers/", name_u), http.StatusSeeOther)
}
