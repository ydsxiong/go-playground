package people

//declared the structure named Emp
type Emp struct {
	Name    string
	Address string
	Age     int
}

//Expose emp creation
func CreateEmployee(name, address string, age int) *Emp {
	return &Emp{name, address, age}
}

func ShowName(e *Emp) string {
	return e.Name
}
