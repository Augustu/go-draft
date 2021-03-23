package http

type Action int

const (
	List Action = iota
	Show
	Add
	Update
	Delete
)

type Restfuler interface {
	Gets(resp Response, req Request)
	Get(resp Response, req Request)
	Put(resp Response, req Request)
	Delete(resp Response, req Request)
	Post(resp Response, req Request)
}

type Request struct{}

type Response struct{}

type method func(resp Response, req Request)

type Restful struct {
	resources map[string]*Handler
}

type Handler struct {
	methods map[Action]method
}

func (r *Restful) Regist(pattern string, h Restfuler) error {

	// TODO beautify this using for loop
	handler := Handler{
		methods: make(map[Action]method),
	}
	handler.methods[List] = h.Gets
	handler.methods[Show] = h.Get
	handler.methods[Update] = h.Put
	handler.methods[Delete] = h.Delete
	handler.methods[Add] = h.Post

	r.resources[pattern] = &handler
	return nil
}

func New() *Restful {
	return &Restful{
		resources: make(map[string]*Handler),
	}
}
