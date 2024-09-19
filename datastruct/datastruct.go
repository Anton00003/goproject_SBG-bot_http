package datastruct

//import "time"

type Person struct {
	Name        string
	ID          int64
	Date        string
	Previous    string
	Subscribers map[string]int
}
