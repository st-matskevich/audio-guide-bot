package repository

type Cover struct {
	Index int    `json:"index"`
	Path  string `json:"-"`
}

type Object struct {
	ID        int64   `json:"-"`
	Title     string  `json:"title"`
	Covers    []Cover `json:"covers"`
	AudioPath string  `json:"-"`
}

type ObjectRepository interface {
	GetObject(code string) (*Object, error)
}

func (repository *Repository) GetObject(code string) (*Object, error) {
	reader, err := repository.DBProvider.Query("SELECT object_id, title, audio_path FROM objects WHERE code = $1", code)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	result := Object{}
	found, err := reader.NextRow(&result.ID, &result.Title, &result.AudioPath)
	if err != nil || !found {
		return nil, err
	}

	reader, err = repository.DBProvider.Query("SELECT index, path FROM covers WHERE object_id = $1", result.ID)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	result.Covers = []Cover{}
	row := Cover{}
	for {
		ok, err := reader.NextRow(&row.Index, &row.Path)
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}

		result.Covers = append(result.Covers, row)
	}

	return &result, nil
}
