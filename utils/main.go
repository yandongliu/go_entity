package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/yandongliu/go_entity/common"
	"github.com/yandongliu/go_entity/dblib"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "dbinit" {
			common.DEBUG("db")
		}
	} else {
		db := dblib.GetDB()
		defer db.Close()
		MapEntities := make(map[int]common.Entity)
		MapEntityChildren := make(map[int][]int)
		MapEntities[1] = dblib.ReadEntity(1, db)
		dblib.ReadAllEntities(1, MapEntityChildren, MapEntities, db)
		common.DEBUG("children result", MapEntityChildren)
		common.DEBUG("entities result", MapEntities)
	}
}

func doInsert(id int, name, value string, db *sql.DB) {
	sqlStatement := `insert into entity (id, name, value) values ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, id, name, value)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func insertFromFile() {
	db, err := sql.Open("postgres", "user=entity_user1 dbname=entity_db password=qwerty123 sslmode=disable")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	args := os.Args[1:]
	common.DEBUG("args", args, len(args))
	if len(args) > 1 {
		common.DEBUG("results", args)
	} else {
		doInsert(1, "hello", "aa", db)
	}

}
