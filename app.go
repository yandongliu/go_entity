package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/yandongliu/go_entity/common"
	"github.com/yandongliu/go_entity/dblib"
)

type Controller struct {
	db *sql.DB
}

func (ct *Controller) indexHandler(c *gin.Context) {
	m := make(map[int][]int)
	v := make(map[int]common.Entity)
	ents := []common.Entity{}
	v[1] = dblib.ReadEntity(1, ct.db)
	dblib.ReadAllEntities(1, m, v, ct.db)
	common.DEBUG("m", m)
	common.DEBUG("v", v)
	params := c.Request.URL.Query()
	id := common.GetURLParamFirstInt(params, "id", 1)
	for _, id2 := range m[id] {
		ents = append(ents, v[id2])
	}
	common.DEBUG("ents", ents)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"entity":   v[id],
		"entities": ents,
	})
}

func (ct *Controller) createHandler(c *gin.Context) {
	name := c.PostForm("name")
	value := c.PostForm("value")
	common.DEBUG("handler create", name, value)
	_, err := dblib.CreateEntity(name, value, ct.db)
	if err != nil {
		c.String(200, "Error!")
	} else {
		c.Redirect(http.StatusSeeOther, "/")
	}
}

func (ct *Controller) entityHandler(c *gin.Context) {
	params := c.Request.URL.Query()
	entity := dblib.ReadEntity(common.GetURLParamFirstInt(params, "id", 1), ct.db)
	common.DEBUG("entity, entities", entity)
	c.HTML(http.StatusOK, "entity.html", gin.H{
		"entity": entity,
	})
}

func setRoutersAndRun(ct *Controller) {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.GET("/", ct.indexHandler)
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.POST("/create", ct.createHandler)
	router.GET("/entity/:id", ct.entityHandler)
	//router.Run(":8080")
	router.Run() //if using gin, port is 3000
}

func main() {
	db := dblib.GetDB()
	defer db.Close()
	ct := &Controller{db: db}
	setRoutersAndRun(ct)
}
