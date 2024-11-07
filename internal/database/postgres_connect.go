package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"log"
	"metrika/internal/utils"
)

const createTableGoals = `CREATE TABLE IF NOT EXISTS public.goals (
	id bigserial NOT NULL,
	goal_id int8 NOT NULL,
	name text NULL,
	status bool NOT NULL,
	fk_counters_metrika_counter_id int8 NOT NULL,
	CONSTRAINT goals_goal_id_key UNIQUE (goal_id),
	CONSTRAINT goals_pkey PRIMARY KEY (id),
	CONSTRAINT goals_fk_counters_metrika_counter_id_fkey FOREIGN KEY (fk_counters_metrika_counter_id) REFERENCES public.counters_metrika(counter_id));`

const createTableCounters = `CREATE TABLE IF NOT EXISTS public.counters_metrika (
        id bigserial NOT NULL,
        counter_id int8 NOT NULL,
        status varchar NULL,
        owner_login varchar NULL,
        name varchar NULL,
        site varchar NULL,
        site_two varchar NULL,
        domain varchar NULL,
        CONSTRAINT counters_metrika_counter_id_key UNIQUE(counter_id),
        CONSTRAINT counters_metrika_pkey PRIMARY KEY (id));`

func PostgresConnect() (*sqlx.DB, error) {
	token, err := utils.GetToken()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}
	pool, err := sqlx.Connect("postgres", token.ClientTable)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = pool.Exec(createTableCounters)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = pool.Exec(createTableGoals)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Postgres connected")

	return pool, nil
}
