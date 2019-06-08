package dblib

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/lib/pq"

	"github.com/yandongliu/go_entity/common"
)

//ReadAllEntities read all children id into a map
func ReadAllEntities(id int, m map[int][]int, v map[int]common.Entity, db *sql.DB) {
	children_ids := ReadChildIds(id, db)
	//common.DEBUG("children_ids", children_ids)
	if len(children_ids) > 0 {
		m[id] = children_ids
	}
	for _, id2 := range children_ids {
		//common.DEBUG("id2", id2)
		e, ok := v[id2]
		if !ok {
			e = ReadEntity(id2, db)
			v[id2] = e
		}
		ReadAllEntities(id2, m, v, db)
	}
}

//ReadEntityPairs aaa
func ReadChildIds(id int, db *sql.DB) []int {
	var id2 int
	r := []int{}
	rows, err := db.Query("SELECT entity_id2 FROM entity_entity where entity_id1 = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No entity_pair with that ID.")
	case err != nil:
		log.Fatal(err)
	}

	for rows.Next() {
		if err := rows.Scan(&id2); err != nil {
			log.Fatal(err)
		}
		r = append(r, id2)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return r
}

func ReadEntity(id int, db *sql.DB) common.Entity {
	var (
		name  string
		value string
	)
	e := common.Entity{}
	rows, err := db.Query("SELECT name, value FROM entity where id = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No entity with that ID.")
	case err != nil:
		log.Fatal(err)
	}

	if rows.Next() {
		if err := rows.Scan(&name, &value); err != nil {
			log.Fatal(err)
		}
		e.Set(id, name, value)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return e
}

//ReadEntityByNameAndValue Find entity with given name/value. This is case insensitive
func ReadEntityByNameAndValue(name, value string, db *sql.DB) common.Entity {
	e := common.Entity{ID: -1}
	rows, err := db.Query(
		"SELECT id, name, value FROM entity where LOWER(name) = $1 and LOWER(value) = $2",
		strings.ToLower(name),
		strings.ToLower(value))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	switch {
	case err == sql.ErrNoRows:
		return e
	case err != nil: //other errors, log it
		log.Fatal(err)
	}

	if rows.Next() {
		if err := rows.Scan(&e.ID, &e.Name, &e.Value); err != nil {
			log.Fatal(err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return e
}

func CreateEntity(name, value string, db *sql.DB) common.Entity {
	e := ReadEntityByNameAndValue(name, value, db)
	if e.ID == -1 {
		common.DEBUG("not found", e)
		_, err := db.Exec(
			"INSERT into ENTITY (name, value) VALUES ($1, $2)",
			strings.ToUpper(name),
			value)
		if err != nil { //other errors, log it
			panic(err)
		}
	} else {
		common.DEBUG("found", e)
	}
	return e
}

func GetDB() *sql.DB {
	db, err := sql.Open("postgres", "user=entity_user1 dbname=entity_db password=qwerty123 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
