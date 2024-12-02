package storage

// BuyHeap implements a max-heap for managing buy orders.
type BuyHeap []*PriceLevel

// SellHeap implements a min-heap for managing sell orders.
type SellHeap []*PriceLevel

func (bh BuyHeap) Len() int            { return len(bh) }
func (bh BuyHeap) Swap(i, j int)       { bh[i], bh[j] = bh[j], bh[i] }
func (bh BuyHeap) Less(i, j int) bool  { return bh[i].Price > bh[j].Price }
func (bh *BuyHeap) Push(x interface{}) { *bh = append(*bh, x.(*PriceLevel)) }
func (bh *BuyHeap) Pop() interface{} {
    old := *bh
    n := len(old)
    x := old[n-1]
    *bh = old[0 : n-1]
    return x
}

func (sh SellHeap) Len() int            { return len(sh) }
func (sh SellHeap) Swap(i, j int)       { sh[i], sh[j] = sh[j], sh[i] }
func (sh SellHeap) Less(i, j int) bool  { return sh[i].Price < sh[j].Price }
func (sh *SellHeap) Push(x interface{}) { *sh = append(*sh, x.(*PriceLevel)) }
func (sh *SellHeap) Pop() interface{} {
    old := *sh
    n := len(old)
    x := old[n-1]
    *sh = old[0 : n-1]
    return x
}
