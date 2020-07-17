package logging

import (
	"github.com/sirupsen/logrus"
	"testing"
)

type a struct {
	Animal string
	Size   int
}

var (
	one  = a{Animal: "walrus", Size: 10}
	list = []a{{Animal: "walrus10", Size: 10}, {Animal: "walrus11", Size: 11}}
)

func TestINFO(t *testing.T) {
	INFO("A group of walrus emerges from the ocean")
	INFO("A group of walrus %s from the %s", "emerges", "ocean")
	INFO(map[string]interface{}{"animal": "walrus", "size": 10}, "A group of walrus emerges from the ocean")
	INFO(map[string]interface{}{"animal": "walrus", "size": 10}, "A group of walrus %s from the %s", "emerges", "ocean")

	logrus.SetFormatter(&logrus.JSONFormatter{})

	INFO(one)
	INFO(list)
	INFO(one, "A group of walrus %s from the %s", "emerges", "ocean")
	INFO(list, "A group of walrus %s from the %s", "emerges", "ocean")
	INFO("A group of walrus emerges from the ocean")
	INFO("A group of walrus %s from the %s", "emerges", "ocean")
	INFO(map[string]interface{}{"animal": "walrus", "size": 10}, "A group of walrus emerges from the ocean")
	INFO(map[string]interface{}{"animal": "walrus", "size": 10}, "A group of walrus %s from the %s", "emerges", "ocean")
	INFO(map[string]interface{}{"msg": "walrus", "level": 10}, "A group of walrus %s from the %s", "emerges", "ocean")
	INFO([]string{"walrus", "size"}, "A group of walrus %s from the %s", "emerges", "ocean")
}
