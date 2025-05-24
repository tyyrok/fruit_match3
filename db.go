package main

import (
	"log"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

const createResutsTable = `
CREATE TABLE IF NOT EXISTS m3tops (
	id SERIAL PRIMARY KEY,
	date DATE DEFAULT now() NOT NULL,
	scores INTEGER NOT NULL DEFAULT '0',
	created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);
CREATE INDEX IF NOT EXISTS ix_m3tops_scores ON m3tops (scores);
`

func InitDB(dbpool *pgxpool.Pool) {
	_, err := dbpool.Exec(
			context.Background(), createResutsTable,
		)
	if err != nil {
		log.Fatalf("Error with accessing/creating table %s\n", err)
	}
}
