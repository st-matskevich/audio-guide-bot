package repository

import "github.com/st-matskevich/audio-guide-bot/api/db"

type Repository struct {
	DBProvider db.DBProvider
}
