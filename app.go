package main

import (
	"database/sql"
	"fmt"
	//"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func DEBUG(msg string, things ...interface{}) {
	fmt.Println("(DEBUG)", time.Now(), ":", msg)
	for _, thing := range things {
		fmt.Println(thing)
	}
}

type TypePair struct {
	Id   int
	Name string
}

func (p *TypePair) Set(id int, name string) {
	p.Id = id
	p.Name = name
}

type TypeEntityCategory struct {
	Entity   TypePair
	Category TypePair
}

func (e *TypeEntityCategory) Set(id int, name string, cid int, cname string) {
	e.Entity.Id = id
	e.Entity.Name = name
	e.Category.Id = cid
	e.Category.Name = cname
}

type Controller struct {
	db *sql.DB
}

func (ct *Controller) categoryHandler(c *gin.Context) {
	cate, cates := readCategories(c.Param("id"), ct.db)
	ents := readEntitiesByCategoryId(c.Param("id"), ct.db)
	c.HTML(http.StatusOK, "category.html", gin.H{
		"category":   cate,
		"categories": cates,
		"entities":   ents,
	})
}

func (ct *Controller) entityHandler(c *gin.Context) {
	entity, entities := readEntitiesByEntityId(c.Param("id"), ct.db)
	DEBUG("entity, entities", entity, entities)
	c.HTML(http.StatusOK, "entity.html", gin.H{
		"entity":   entity,
		"entities": entities,
	})
}

func readCategoryByEntity(eid int, db *sql.DB) TypePair {
	var (
		id   int
		name string
	)
	cate := TypePair{Id: -1, Name: ""}
	err := db.QueryRow("SELECT c.id, c.name FROM category c, entity e where e.category_id = c.id and e.id = $1", eid).Scan(&id, &name)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No category with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		cate = TypePair{Id: id, Name: name}
	}
	return cate
}

// read entities by category id
func readEntitiesByCategoryId(cid string, db *sql.DB) []TypePair {
	var (
		id   int
		name string
	)
	entities := []TypePair{}
	rows, err := db.Query("SELECT id, name FROM entity where category_id = $1", cid)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No category with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		log.Printf("OK")
	}

	for rows.Next() {
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		entities = append(entities, TypePair{Id: id, Name: name})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return entities
}

func readEntitiesByEntityId(_eid string, db *sql.DB) (TypeEntityCategory, []TypeEntityCategory) {
	var eid, cid int
	var ename, cname string
	var entity TypeEntityCategory
	var entities []TypeEntityCategory

	/* select entity id/name, category id/name */
	DEBUG("Read entity")
	err := db.QueryRow("SELECT e.id, e.name, c.id, c.name FROM entity e, category c WHERE e.id = $1 and e.category_id = c.id", _eid).Scan(&eid, &ename, &cid, &cname)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No category with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		entity.Set(eid, ename, cid, cname)
	}

	/* select entity assocaite with given entity */
	DEBUG("Read entities")
	rows, err := db.Query("SELECT e.id, e.name, c.id, c.name FROM entity e, category c, entity_entity ee WHERE e.id = ee.entity_id2 and ee.entity_id1 = $1 and e.category_id = c.id", _eid)
	for rows.Next() {
		if err := rows.Scan(&eid, &ename, &cid, &cname); err != nil {
			log.Fatal(err)
		}
		entities = append(entities, TypeEntityCategory{TypePair{eid, ename}, TypePair{cid, cname}})
		//DEBUG("aasdsdf", Entity{Pair{eid, ename}, Pair{cid, cname}})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return entity, entities
}

func readCategories(id string, db *sql.DB) (TypePair, []TypePair) {
	var (
		cid  int
		name string
	)
	var cate TypePair
	err := db.QueryRow("SELECT id, name FROM category WHERE id = $1", id).Scan(&cid, &name)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No category with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		cate = TypePair{Id: cid, Name: name}
	}
	rows, err := db.Query("SELECT c.id, c.name FROM category c, category_category cc WHERE cc.pid = $1 and cc.id = c.id", id)
	if err != nil {
		log.Fatal(err)
	}
	cates := []TypePair{}
	for rows.Next() {
		if err := rows.Scan(&cid, &name); err != nil {
			log.Fatal(err)
		}
		cates = append(cates, TypePair{Id: cid, Name: name})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return cate, cates
}

func main() {
	db, err := sql.Open("postgres", "user=entity_user1 dbname=entity_db password=qwerty123")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	ct := &Controller{db: db}
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.LoadHTMLGlob("templates/*.html")
	router.GET("/category/:id", ct.categoryHandler)
	router.GET("/entity/:id", ct.entityHandler)
	router.Run(":8080")
}
