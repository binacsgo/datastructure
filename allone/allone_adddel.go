package allone

type AllOneAddDel struct {
	hash  map[interface{}]int
	items []interface{}
	count int
}

func AllOneAddDelConstructor() AllOneAddDel {
	return AllOneAddDel{
		hash:  make(map[interface{}]int, 0),
		items: make([]interface{}, 0),
		count: 0,
	}
}

func (this *AllOneAddDel) Add(x any) bool {
	if _, ok := this.hash[x]; !ok {
		this.items = append(this.items, x)
		this.hash[x] = this.count
		this.count++
		return true
	}
	return false
}

func (this *AllOneAddDel) Del(x any) bool {
	if px, ok := this.hash[x]; ok {
		this.count--
		y := this.items[this.count]
		py := this.hash[y]
		this.items[px], this.hash[y] = this.items[py], px
		delete(this.hash, x)
		this.items = this.items[:this.count]
		return true
	}
	return false
}
