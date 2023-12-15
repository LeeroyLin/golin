package gcfg

type CfgBase[T interface{}] struct {
	dataMap map[int32]T
	dataArr []T
}

func (c *CfgBase[T]) GetCount() int {
	c.GetById(0)
	return len(c.dataArr)
}

func (c *CfgBase[T]) GetById(id int32) T {
	return c.dataMap[id]
}

func (c *CfgBase[T]) Foreach(handler func(cfg T)) {
	for _, v := range c.dataArr {
		handler(v)
	}
}

func (c *CfgBase[T]) InitData(data []T) map[int32]T {
	c.dataArr = data
	c.dataMap = make(map[int32]T)

	return c.dataMap
}

func (c *CfgBase[T]) AddMapData(id int32, data T) {
	c.dataMap[id] = data
}
