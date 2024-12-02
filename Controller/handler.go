// controller/handler.go

package controller

import (
	"errors"
	"net/http"

	"CLOB/genesis"
	"CLOB/storage"
	"CLOB/utils"
	"CLOB/vm"

	"github.com/ava-labs/hypersdk/chain"
	"github.com/ava-labs/hypersdk/codec"
	"github.com/ava-labs/hypersdk/ids"
)

// Define controller-specific errors
var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrInvalidOrder     = errors.New("invalid order")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

// Handler manages HTTP requests related to the Order Book Matching Engine
type Handler struct {
	*vm.Handler // Embed standard VM handler functionality

	c *Controller
}

// GenesisReply represents the response for the Genesis endpoint
type GenesisReply struct {
	Genesis *genesis.Genesis `json:"genesis"`
}

// Genesis handles the Genesis endpoint, returning the genesis configuration
func (h *Handler) Genesis(_ *http.Request, _ *struct{}, reply *GenesisReply) error {
	reply.Genesis = h.c.genesis
	return nil
}

// AddOrderArgs represents the request payload for adding a new order
type AddOrderArgs struct {
	OrderID   string  `json:"order_id"`
	Side      string  `json:"side"`       // "buy" or "sell"
	Price     float64 `json:"price"`      // 0 for market orders
	Quantity  float64 `json:"quantity"`
	OrderType string  `json:"order_type"` // "limit" or "market"
}

// AddOrderReply represents the response after adding a new order
type AddOrderReply struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// AddOrder handles the addition of a new order to the order book
func (h *Handler) AddOrder(req *http.Request, args *AddOrderArgs, reply *AddOrderReply) error {
	ctx, span := h.c.inner.Tracer().Start(req.Context(), "Handler.AddOrder")
	defer span.End()

	// Validate order type
	if args.OrderType != "limit" && args.OrderType != "market" {
		reply.Success = false
		reply.Message = "invalid order type"
		return ErrInvalidOrder
	}

	// Validate order side
	if args.Side != "buy" && args.Side != "sell" {
		reply.Success = false
		reply.Message = "invalid order side"
		return ErrInvalidOrder
	}

	// Check for duplicate order ID
	exists, err := storage.OrderExists(ctx, h.c.metaDB, args.OrderID)
	if err != nil {
		return err
	}
	if exists {
		reply.Success = false
		reply.Message = "order ID already exists"
		return ErrOrderAlreadyExists
	}

	// Create Order struct
	order := &storage.Order{
		ID:        args.OrderID,
		Side:      storage.Side(args.Side),
		Price:     args.Price,
		Quantity:  args.Quantity,
		Timestamp: h.c.inner.Clock().Now(),
		OrderType: storage.OrderType(args.OrderType),
	}

	// Execute AddOrder action
	addOrderAction := &actions.AddOrderAction{Order: order}
	if err := h.c.vm.ExecuteAction(addOrderAction); err != nil {
		reply.Success = false
		reply.Message = err.Error()
		return err
	}

	reply.Success = true
	reply.Message = "order added successfully"
	return nil
}

// CancelOrderArgs represents the request payload for canceling an order
type CancelOrderArgs struct {
	OrderID string `json:"order_id"`
}

// CancelOrderReply represents the response after canceling an order
type CancelOrderReply struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// CancelOrder handles the cancellation of an existing order
func (h *Handler) CancelOrder(req *http.Request, args *CancelOrderArgs, reply *CancelOrderReply) error {
	ctx, span := h.c.inner.Tracer().Start(req.Context(), "Handler.CancelOrder")
	defer span.End()

	// Check if the order exists
	exists, err := storage.OrderExists(ctx, h.c.metaDB, args.OrderID)
	if err != nil {
		return err
	}
	if !exists {
		reply.Success = false
		reply.Message = "order not found"
		return ErrOrderNotFound
	}

	// Execute CancelOrder action
	cancelOrderAction := &actions.CancelOrderAction{OrderID: args.OrderID}
	if err := h.c.vm.ExecuteAction(cancelOrderAction); err != nil {
		reply.Success = false
		reply.Message = err.Error()
		return err
	}

	reply.Success = true
	reply.Message = "order canceled successfully"
	return nil
}

// GetOrderArgs represents the request payload for retrieving an order
type GetOrderArgs struct {
	OrderID string `json:"order_id"`
}

// GetOrderReply represents the response containing order details
type GetOrderReply struct {
	OrderID   string  `json:"order_id"`
	Side      string  `json:"side"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	OrderType string  `json:"order_type"`
	Timestamp int64   `json:"timestamp"`
}

// GetOrder handles retrieving details of a specific order
func (h *Handler) GetOrder(req *http.Request, args *GetOrderArgs, reply *GetOrderReply) error {
	ctx, span := h.c.inner.Tracer().Start(req.Context(), "Handler.GetOrder")
	defer span.End()

	order, err := storage.GetOrder(ctx, h.c.metaDB, args.OrderID)
	if err != nil {
		return err
	}
	if order == nil {
		return ErrOrderNotFound
	}

	reply.OrderID = order.ID
	reply.Side = string(order.Side)
	reply.Price = order.Price
	reply.Quantity = order.Quantity
	reply.OrderType = string(order.OrderType)
	reply.Timestamp = order.Timestamp.Unix()

	return nil
}

// GetOrderBookArgs represents the request payload for retrieving the order book
type GetOrderBookArgs struct{}

// GetOrderBookReply represents the response containing the current state of the order book
type GetOrderBookReply struct {
	Bids []storage.Order `json:"bids"`
	Asks []storage.Order `json:"asks"`
}

// GetOrderBook handles retrieving the current state of the order book
func (h *Handler) GetOrderBook(req *http.Request, args *GetOrderBookArgs, reply *GetOrderBookReply) error {
	ctx, span := h.c.inner.Tracer().Start(req.Context(), "Handler.GetOrderBook")
	defer span.End()

	bids, asks, err := storage.GetOrderBookState(ctx, h.c.metaDB)
	if err != nil {
		return err
	}

	reply.Bids = bids
	reply.Asks = asks
	return nil
}

// ListOrdersArgs represents the request payload for listing all orders of a user
type ListOrdersArgs struct {
	Address string `json:"address"`
}

// ListOrdersReply represents the response containing a list of orders
type ListOrdersReply struct {
	Orders []storage.Order `json:"orders"`
}

// ListOrders handles listing all orders associated with a specific address
func (h *Handler) ListOrders(req *http.Request, args *ListOrdersArgs, reply *ListOrdersReply) error {
	ctx, span := h.c.inner.Tracer().Start(req.Context(), "Handler.ListOrders")
	defer span.End()

	address, err := utils.ParseAddress(args.Address)
	if err != nil {
		return err
	}

	orders, err := storage.GetOrdersByAddress(ctx, h.c.metaDB, address)
	if err != nil {
		return err
	}

	reply.Orders = orders
	return nil
}
