package tables

type Insert interface {
	insertQuery() error
}

func InsertQuery(i Insert) error {
	err := i.insertQuery()
	return err
}
