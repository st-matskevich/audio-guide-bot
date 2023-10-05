package repository

type ConfigRepository interface {
	GetValue(key string) (*string, error)
}

func (repository *Repository) GetValue(key string) (*string, error) {
	reader, err := repository.DBProvider.Query("SELECT value FROM config WHERE key = $1", key)
	if err != nil {
		return nil, err
	}

	defer reader.Close()
	result := ""
	found, err := reader.NextRow(&result)
	if err != nil || !found {
		return nil, err
	}

	return &result, nil
}
