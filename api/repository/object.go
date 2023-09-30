package repository

type Object struct {
	Title     string `json:"title"`
	CoverPath string `json:"-"`
	AudioPath string `json:"-"`
}

type ObjectRepository interface {
	GetObject(code string) (Object, error)
}

func (repository *Repository) GetObject(code string) (Object, error) {
	reader, err := repository.DBProvider.Query("SELECT title, cover_path, audio_path FROM objects WHERE code = $1", code)
	if err != nil {
		return Object{}, err
	}

	defer reader.Close()
	result := Object{}
	err = reader.GetRow(&result.Title, &result.CoverPath, &result.AudioPath)
	if err != nil {
		return Object{}, err
	}

	return result, nil
}