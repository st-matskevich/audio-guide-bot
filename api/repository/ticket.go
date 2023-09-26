package repository

type TicketRepository interface {
	CreateTicket(code string) error
	ActivateTicket(code string) (bool, error)
}

func (repository *Repository) CreateTicket(code string) error {
	_, err := repository.DBProvider.Exec("INSERT INTO tickets(code) VALUES ($1)", code)
	if err != nil {
		return err
	}

	return nil
}

func (repository *Repository) ActivateTicket(code string) (bool, error) {
	updated, err := repository.DBProvider.Exec("UPDATE tickets SET used = true WHERE code = $1 AND used = false", code)
	if err != nil {
		return false, err
	}

	return updated == 1, nil
}
