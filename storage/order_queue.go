package storage

// OrderQueue represents a doubly linked list of orders
type OrderQueue struct {
    head *Order 
    tail *Order 
    Size int    // Number of orders in the queue
}

// NewOrderQueue creates a new empty order queue
// Returns a pointer to the newly created OrderQueue
func NewOrderQueue() *OrderQueue {
    return &OrderQueue{}
}

// Enqueue adds an order to the end of the queue
// Takes a pointer to an Order as input
// If the queue is empty, it sets both head and tail to the new order
// Otherwise, it updates the tail and links the new order
func (oq *OrderQueue) Enqueue(order *Order) {
    if oq.tail == nil {
        oq.head = order
        oq.tail = order
    } else {
        oq.tail.next = order
        order.prev = oq.tail
        oq.tail = order
    }
    oq.Size++
}

// Dequeue removes and returns the order from the front of the queue
// Returns the order at the front of the queue or nil if the queue is empty
func (oq *OrderQueue) Dequeue() *Order {
    if oq.head == nil {
        return nil
    }
    order := oq.head
    oq.head = oq.head.next
    if oq.head != nil {
        oq.head.prev = nil
    } else {
        oq.tail = nil
    }
    order.next = nil
    oq.Size--
    return order
}

// Remove removes a specific order from the queue
// Takes a pointer to an Order as input
// Adjusts the next and prev pointers of surrounding orders
// Updates head or tail if necessary
func (oq *OrderQueue) Remove(order *Order) {
    if order.prev != nil {
        order.prev.next = order.next
    } else {
        oq.head = order.next
    }
    if order.next != nil {
        order.next.prev = order.prev
    } else {
        oq.tail = order.prev
    }
    order.next = nil
    order.prev = nil
    oq.Size--
}
