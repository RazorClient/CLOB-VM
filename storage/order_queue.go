// CLOB/storage/order_queue.go
package storage

// OrderQueue represents a doubly linked list of orders
type OrderQueue struct {
    head *Order // Unexported
    tail *Order // Unexported
    Size int
}

// NewOrderQueue creates a new empty order queue
func NewOrderQueue() *OrderQueue {
    return &OrderQueue{}
}

// Enqueue adds an order to the end of the queue
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
