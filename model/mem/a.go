package mem

import (
	"sparkai/model"
)

var WSConnContainers map[string]*model.WSConnContainer

func init() {
	WSConnContainers = make(map[string]*model.WSConnContainer)
}
