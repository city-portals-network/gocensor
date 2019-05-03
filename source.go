package main

//Source retur
type Source interface {
	parse() (map[string]bool, error)
	append(string) error
	delete(string) error
}

type MysqlSource struct {
	MySQL *MySQL
}

func (s *MysqlSource) parse() (map[string]bool, error) {
	return nil, nil
}
