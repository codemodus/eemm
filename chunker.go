package main

type uniqIDs []uint32

func (ids uniqIDs) chunks(size int) [][]uint32 {
	var cs [][]uint32
	var c []uint32

	for i, id := range ids {
		c = append(c, id)

		if len(c)%size == 0 || i == len(ids)-1 {
			cs = append(cs, c)
			c = nil
		}
	}

	return cs
}
