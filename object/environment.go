package object

// 为变量名和值创建 map

func NewEnvirnment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

type Environment struct {
	store map[string]Object
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) Get(name string) (obj Object, ok bool) {
	obj, ok = e.store[name]
	return
}
