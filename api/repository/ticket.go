package repository

type TicketRepository interface {
	CreateTicket(code string) error
	ActivateTicket(ticketID int64) (bool, error)
}

func (repository *Repository) CreateTicket(code string) error {
	_, err := repository.DBProvider.Exec("INSERT INTO tickets(code) VALUES ($1)", code)
	if err != nil {
		return err
	}

	return nil
}

func (repository *Repository) ActivateTicket(ticketID int64) (bool, error) {
	return false, nil
}
