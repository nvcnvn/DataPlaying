package main

func (h *handler) initSubRoutes() {
	h._defaultHandle = HandleIndex
	h._subRoutes = []route{
		route{pattern: "index", fn: HandleIndex},
		route{pattern: "ajax/summary", fn: HandleAjaxSummary},
		route{pattern: "ajax/data", fn: HandleAjaxData},
		route{pattern: "ajax/sorteddata", fn: HandleAjaxSortedData},
		route{pattern: "ajax/clustering", fn: HandleAjaxClustering},
	}
}
