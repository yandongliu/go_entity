package common

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var DEBUG_MODE = true

func PrintErr(err error) {
	if err != nil {
		DEBUG("error", err)
	}
}

func DEBUG(msg string, things ...interface{}) {
	if !DEBUG_MODE {
		return
	}
	strs := make([]string, 0)
	for _, thing := range things {
		str := fmt.Sprintf("%v", thing)
		strs = append(strs, str)
	}
	fmt.Println("(DEBUG)", time.Now(), " | ", msg, ":", strings.Join(strs, " "))
}

//Entity type entity
type Entity struct {
	ID    int
	Name  string
	Value string
}

func (p *Entity) Set(id int, name, value string) {
	p.ID = id
	p.Name = name
	p.Value = value
}

type IdPair struct {
	ID1 int
	ID2 int
}

type EntityPair struct {
	ParentEntity Entity
	ChildEntity  Entity
}

func GetRandomId() int {
	return rand.Int()
}

func GetURLParamFirstStr(params url.Values, key string, thenvalue string) string {
	if val, ok := params[key]; ok {
		return val[0]
	}
	return thenvalue
}

func GetURLParamFirstInt(params url.Values, key string, thenvalue int) int {
	if val, ok := params[key]; ok {
		if s, err := strconv.Atoi(val[0]); err == nil {
			return s
		}
	}
	return thenvalue
}

func String2Int(s string, thenvalue int) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return thenvalue
}
