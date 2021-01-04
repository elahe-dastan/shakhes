package heap

type Frequency struct {
	DocId string
	Freq  int
}

// A FrequencyHeap is a max-heap of Frequency
type FrequencyHeap []Frequency

func (f FrequencyHeap) Len() int { return len(f) }
// changed the less function so it became a max heap
func (f FrequencyHeap) Less(i, j int) bool { return f[i].Freq > f[j].Freq }
func (f FrequencyHeap) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

func (f *FrequencyHeap) Pop() interface{} {
	old := *f
	n := len(old)
	x := old[n-1]
	*f = old[0 : n-1]
	return x
}

func (f *FrequencyHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*f = append(*f, x.(Frequency))
}
