package main

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/yandongliu/go_entity/common"
	"github.com/yandongliu/go_entity/dblib"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "dbinit" {
			common.DEBUG("Init DB Data")
			initDBData()
		} else if args[0] == "dbdelete" {
			common.DEBUG("Delete DB Data")
			deleteDBData()
		}
	} else {
		showAllData()
	}
}
func showAllData() {
	db := dblib.GetDB()
	defer db.Close()
	MapEntities := make(map[int]common.Entity)
	MapEntityChildren := make(map[int][]int)
	MapEntities[1] = dblib.ReadEntity(1, db)
	dblib.ReadAllEntities(1, MapEntityChildren, MapEntities, db)
	common.DEBUG("children result", MapEntityChildren)
	common.DEBUG("entities result", MapEntities)
}

func deleteDBData() {
	db := dblib.GetDB()
	defer db.Close()
	_, err := db.Exec("delete from ENTITY_ENTITY")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("delete from ENTITY")
	if err != nil {
		panic(err)
	}
}

func initDBData() {
	db := dblib.GetDB()
	defer db.Close()
	eRoot, err := dblib.CreateEntity("ROOT", "root", db)
	common.PrintErr(err)
	//common.DEBUG("Created entity", eRoot, err)
	eBusiness, err := dblib.CreateEntity("BUSINESS", "business", db)
	common.PrintErr(err)
	//common.DEBUG("Created entity", eBusiness, err)
	eRetailer, err := dblib.CreateEntity("RETAILER", "retailer", db)
	common.PrintErr(err)
	err = dblib.CreateEntityPair(eRoot.ID, eBusiness.ID, db)
	common.PrintErr(err)
	err = dblib.CreateEntityPair(eBusiness.ID, eRetailer.ID, db)
	common.PrintErr(err)
}

/* SQL

create table entity (id SERIAL PRIMARY KEY, name VARCHAR(255), value VARCHAR(255));

create table entity_entity (
	id1 INTEGER REFERENCES entity(id),
	id2 INTEGER REFERENCES entity(id),
	UNIQUE (id1, id2)
);
Create index idx_entity_entity on entity_entity(id1);
*/
