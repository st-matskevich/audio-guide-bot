package repository

type Ticket struct {
	ID   int64
	Code string
	Used bool
}

type TicketRepository interface {
	CreateTicket(code string) error
	GetTicket(code string) (*Ticket, error)
	ActivateTicket(code string) (bool, error)
}

func (repository *Repository) CreateTicket(code string) error {
	_, err := repository.DBProvider.Exec("INSERT INTO tickets(code) VALUES ($1)", code)
	if err != nil {
		return err
	}

	return nil
}

func (repository *Repository) GetTicket(code string) (*Ticket, error) {
	reader, err := repository.DBProvider.Query("SELECT ticket_id, used FROM tickets WHERE code = $1", code)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	result := Ticket{}
	found, err := reader.NextRow(&result.ID, &result.Used)
	if err != nil || !found {
		return nil, err
	}

	result.Code = code
	return &result, nil
}

func (repository *Repository) ActivateTicket(code string) (bool, error) {
	updated, err := repository.DBProvider.Exec("UPDATE tickets SET used = true WHERE code = $1 AND used = false", code)
	if err != nil {
		return false, err
	}

	return updated == 1, nil
}
