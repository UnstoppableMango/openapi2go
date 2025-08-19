package petstore

type User struct {
	Id         int
	Username   string
	Firstname  string
	Lastname   string
	Email      string
	Password   string
	Phone      string
	Userstatus int
}

type Tag struct {
	Id   int
	Name string
}

type Pet struct {
	Id        int
	Name      string
	Category  object
	Photourls array
	Tags      array
	Status    string
}

type ApiResponse struct {
	Code    int
	Type    string
	Message string
}

type Order struct {
	Id       int
	Petid    int
	Quantity int
	Shipdate string
	Status   string
	Complete bool
}

type Category struct {
	Id   int
	Name string
}
