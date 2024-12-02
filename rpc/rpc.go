// CLOB/rpc/rpc.go

package rpc

import (
	"context"
	"strings"

	"CLOB/actions"
	"CLOB/consts"
	"CLOB/genesis"
	"CLOB/storage"
	"CLOB/utils"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/hypersdk/chain"
	"github.com/ava-labs/hypersdk/requester"
	"github.com/ava-labs/hypersdk/rpc"
)

type JSONRPCClient struct {
	requester *requester.EndpointRequester
	chainID   ids.ID
	genesis   *genesis.Genesis
}

// NewJSONRPCClient creates a new client object.
func NewJSONRPCClient(uri string, chainID ids.ID) *JSONRPCClient {
	uri = strings.TrimSuffix(uri, "/")
	uri += "/rpc"
	req := requester.New(uri, consts.Name)
	return &JSONRPCClient{req, chainID, nil}
}

// Genesis retrieves the genesis configuration.
func (cli *JSONRPCClient) Genesis(ctx context.Context) (*genesis.Genesis, error) {
	if cli.genesis != nil {
		return cli.genesis, nil
	}

	resp := new(GenesisReply)
	err := cli.requester.SendRequest(ctx, "genesis", nil, resp)
	if err != nil {
		return nil, err
	}
	cli.genesis = resp.Genesis
	return resp.Genesis, nil
}

// AddOrderArgs represents the arguments for adding an order.
type AddOrderArgs struct {
	OrderID   string  `json:"order_id"`
	Side      string  `json:"side"`       // "buy" or "sell"
	Price     float64 `json:"price"`      // 0 for market orders
	Quantity  float64 `json:"quantity"`
	OrderType string  `json:"order_type"` // "limit" or "market"
}

// AddOrderReply represents the response after adding an order.
type AddOrderReply struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// AddOrder sends an AddOrder request to the server.
func (cli *JSONRPCClient) AddOrder(ctx context.Context, args *AddOrderArgs) (*AddOrderReply, error) {
	resp := new(AddOrderReply)
	err := cli.requester.SendRequest(ctx, "addOrder", args, resp)
	return resp, err
}

// CancelOrderArgs represents the arguments for canceling an order.
type CancelOrderArgs struct {
	OrderID string `json:"order_id"`
}

// CancelOrderReply represents the response after canceling an order.
type CancelOrderReply struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// CancelOrder sends a CancelOrder request to the server.
func (cli *JSONRPCClient) CancelOrder(ctx context.Context, args *CancelOrderArgs) (*CancelOrderReply, error) {
	resp := new(CancelOrderReply)
	err := cli.requester.SendRequest(ctx, "cancelOrder", args, resp)
	return resp, err
}

// GetOrderArgs represents the arguments for retrieving an order.
type GetOrderArgs struct {
	OrderID string `json:"order_id"`
}

// GetOrderReply represents the response containing order details.
type GetOrderReply struct {
	OrderID   string  `json:"order_id"`
	Side      string  `json:"side"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	OrderType string  `json:"order_type"`
	Timestamp int64   `json:"timestamp"`
}

// GetOrder retrieves the details of an order.
func (cli *JSONRPCClient) GetOrder(ctx context.Context, args *GetOrderArgs) (*GetOrderReply, error) {
	resp := new(GetOrderReply)
	err := cli.requester.SendRequest(ctx, "getOrder", args, resp)
	if err != nil {
		if strings.Contains(err.Error(), ErrOrderNotFound.Error()) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return resp, nil
}

// GetOrderBookArgs represents the arguments for retrieving the order book.
type GetOrderBookArgs struct{}

// GetOrderBookReply represents the response containing the order book.
type GetOrderBookReply struct {
	Bids []storage.Order `json:"bids"`
	Asks []storage.Order `json:"asks"`
}

// GetOrderBook retrieves the current state of the order book.
func (cli *JSONRPCClient) GetOrderBook(ctx context.Context) (*GetOrderBookReply, error) {
	resp := new(GetOrderBookReply)
	err := cli.requester.SendRequest(ctx, "getOrderBook", nil, resp)
	return resp, err
}

// WaitForOrder waits until the order is available or a timeout occurs.
func (cli *JSONRPCClient) WaitForOrder(ctx context.Context, orderID string) (*GetOrderReply, error) {
	var order *GetOrderReply
	err := rpc.Wait(ctx, func(ctx context.Context) (bool, error) {
		resp, err := cli.GetOrder(ctx, &GetOrderArgs{OrderID: orderID})
		if err != nil {
			if err == ErrOrderNotFound {
				// Order not found yet
				return false, nil
			}
			// Some other error occurred
			return false, err
		}
		order = resp
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

// Parser implements chain.Parser for parsing actions and authentication.
type Parser struct {
	chainID ids.ID
	genesis *genesis.Genesis
}

func (p *Parser) ChainID() ids.ID {
	return p.chainID
}

func (p *Parser) Rules(t int64) chain.Rules {
	return p.genesis.Rules(t)
}

func (*Parser) Registry() (chain.ActionRegistry, chain.AuthRegistry) {
	return consts.ActionRegistry, consts.AuthRegistry
}

// Parser returns a chain.Parser for parsing actions and authentication.
func (cli *JSONRPCClient) Parser(ctx context.Context) (chain.Parser, error) {
	g, err := cli.Genesis(ctx)
	if err != nil {
		return nil, err
	}
	return &Parser{cli.chainID, g}, nil
}

// Error definitions.
var (
	ErrOrderNotFound = utils.NewError("order not found")
)
