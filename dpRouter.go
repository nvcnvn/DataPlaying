package main

func (h *handler) initSubRoutes() {
	h._defaultHandle = HandleIndex
	h._subRoutes = []route{
		route{pattern: "index", fn: HandleIndex},
	}
}
