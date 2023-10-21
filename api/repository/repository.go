package repository

import "github.com/st-matskevich/audio-guide-bot/api/provider/db"

type Repository struct {
	DBProvider db.DBProvider
}
