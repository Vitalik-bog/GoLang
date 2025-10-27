package models

type Passport struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}

type Department struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type Employee struct {
	Id         int        `json:"id"`
	Name       string     `json:"name"`
	Surname    string     `json:"surname"`
	Phone      string     `json:"phone"`
	CompanyId  int        `json:"companyId"`
	Passport   Passport   `json:"passport"`
	Department Department `json:"department"`
}
