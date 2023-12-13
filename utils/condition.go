package utils

// Condition структура условий
type Condition struct {
	Equal       map[string]interface{}
	NotEqual    map[string]interface{}
	Order       []*Order
	LimitOffset *LimitOffset
	ForUpdate   bool
	Upsert      bool
}

// Order структура сортировки
type Order struct {
	Field string
	Asc   bool
}

// LimitOffset структура ограничения
type LimitOffset struct {
	Offset int64
	Limit  int64
}
