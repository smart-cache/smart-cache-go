package heap

type HeapItemInt64 struct {
	label	string
	key	    int64
}

type MinHeapInt64 struct {
	items		[]HeapItemInt64
	labels  	map[string]int64
	Size		int64
}

func MakeMinHeapInt64() *MinHeapInt64 {
	h := &MinHeapInt64{}
	h.items = make([]HeapItemInt64, 0)
	h.labels = make(map[string]int64)
	return h
}

func (h *MinHeapInt64) MinHeapifyUp(c int64) {
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

func (h *MinHeapInt64) MinHeapifyDown(p int64) {
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

func (h *MinHeapInt64) Insert(label string, key int64) {
	if h.Contains(label) {
		h.ChangeKey(label, key)
	} else {
		var i HeapItemInt64
		i.label = label
		i.key = key
		h.items = append(h.items, i)
		h.labels[label] = h.Size
		h.Size++
		h.MinHeapifyUp(h.labels[label])
	}
}

func (h *MinHeapInt64) ExtractMin() string {
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

func (h *MinHeapInt64) ChangeKey(label string, key int64) {
	index, ok := h.labels[label]
	if ok {
		if key < h.items[index].key {
			h.items[index].key = key
			h.MinHeapifyUp(index)
		} else {
			h.items[index].key = key
			h.MinHeapifyDown(index)
		}
	} else {
		h.Insert(label, key)
	}
}

func (h *MinHeapInt64) Swap(i int64, j int64) {
	temp := h.items[i]
	h.items[i] = h.items[j]
	h.items[j] = temp
}

func (h *MinHeapInt64) Contains(key string) bool {
	_, ok := h.labels[key]
	return ok
}

func (h *MinHeapInt64) GetKeyList() []string {
	li := make([]string, len(h.items))
	for i, v := range h.items {
		li[i] = v.label
	}
	return li
}

func (h *MinHeapInt64) GetKey(name string) int64 {
	index, ok := h.labels[name]
	if ok {
		return h.items[index].key
	} else {
		return 0
	}
}