package heap

// this file was copied from smart-cache-824project

type HeapItemFloat64 struct {
	label	string
	key	    float64
}

type MinHeapFloat64 struct {
	items		[]HeapItemFloat64
	labels  	map[string]int64
	Size		int64
}

func MakeMinHeapFloat64() *MinHeapFloat64 {
	h := &MinHeapFloat64{}
	h.items = make([]HeapItemFloat64, 0)
	h.labels = make(map[string]int64)
	return h
}

func (h *MinHeapFloat64) Init() {
	h.items = make([]HeapItemFloat64, 0)
	h.labels = make(map[string]int64)
}

func (h *MinHeapFloat64) MinHeapifyUp(c int64) {
	if c == 0 {
		return
	}
	p := (c - 1) / 2
	if h.items[p].key > h.items[c].key {
		// swap terms
		h.Swap(p, c)
		h.labels[h.items[p].label] = p
		h.labels[h.items[c].label] = c
		h.MinHeapifyUp(p)
	}
}

func (h *MinHeapFloat64) MinHeapifyDown(p int64) {
	if p >= h.Size {
		return
	}

	// check children
	l := 2 * p + 1
	r := 2 * p + 2
	if l >= h.Size {
		l = p
	}
	if r >= h.Size {
		r = p
	}

	// set child pointer
	var c int64
	if h.items[r].key > h.items[l].key {
		c = l
	} else {
		c = r
	}

	if h.items[p].key > h.items[c].key {
		// swap terms
		h.Swap(p, c)
		h.labels[h.items[p].label] = p
		h.labels[h.items[c].label] = c
		h.MinHeapifyDown(c)
	}
}

func (h *MinHeapFloat64) Insert(label string, key float64) {
	var i HeapItemFloat64
	i.label = label
	i.key = key
	h.items = append(h.items, i)
	h.labels[label] = h.Size
	h.Size++
	h.MinHeapifyUp(h.labels[label])
}

func (h *MinHeapFloat64) ExtractMin() string {
	// swap first and last terms
	h.Swap(0, h.Size - 1)
	h.labels[h.items[0].label] = 0
	delete(h.labels, h.items[h.Size - 1].label)
	label := h.items[h.Size - 1].label
	h.items = h.items[:(h.Size - 1)]
	h.Size--
	h.MinHeapifyDown(0)
	return label
}

func (h *MinHeapFloat64) ChangeKey(label string, key float64) {
	index, ok := h.labels[label]
	if ok {
		if key < h.items[index].key {
			h.items[index].key = key
			h.MinHeapifyUp(index)
		} else {
			h.items[index].key = key
			h.MinHeapifyDown(index)
		}
	}
}

func (h *MinHeapFloat64) Swap(i int64, j int64) {
	temp := h.items[i]
	h.items[i] = h.items[j]
	h.items[j] = temp
}

func (h *MinHeapFloat64) Contains(key string) bool {
	_, ok := h.labels[key]
	return ok
}

func (h *MinHeapFloat64) GetKeyList() []string {
	li := make([]string, len(h.items))
	for i, v := range h.items {
		li[i] = v.label
	}
	return li
}

func (h *MinHeapFloat64) GetKey(name string) float64 {
	index, ok := h.labels[name]
	if ok {
		return h.items[index].key
	} else {
		return 0.0
	}
}