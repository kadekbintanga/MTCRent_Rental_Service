package option

type CustomerOption struct {
	ID   int
	UUID string
}

type CustomerSaveOption struct {
	ID        int
	UUID      string
	Name      string
	IDNumber  string
	SIMNumber string
	Phone     string
	StatusId  int
}
