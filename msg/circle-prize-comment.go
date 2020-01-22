package msg

//for snippet用于标准返回值的微服务接口

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-kit/kit/endpoint"
	tran "github.com/go-kit/kit/transport/http"
	"github.com/vhaoran/vchat/lib/ykit"
)

const (
	//todo
	CirclePrize_HANDLER_PATH = "/CirclePrize"
)

type (
	CirclePrizeService interface {
		//todo
		Exec(ctx context.Context, in *CirclePrizeIn) (*ykit.Result, error)
	}

	//input data
	//todo
	CirclePrizeIn struct {
		//0 prize 1 comment
		Action int                `json:"action omitempty"`
		ID     primitive.ObjectID `json:"id,omitempty"   bson:"_id,omitempty"`
		Text   string             `json:"text omitempty"`
		UID    int64              `json:"uid,omitempty"   bson:"uid,omitempty"`
	}

	//output data
	//Result struct {
	//	Code int         `json:"code"`
	//	Msg  string      `json:"msg"`
	//	Data interface{} `json:"data"`
	//}

	// handler implements
	CirclePrizeHandler struct {
		base ykit.RootTran
	}
)

func (r *CirclePrizeHandler) MakeLocalEndpoint(svc CirclePrizeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  CirclePrize ###########")
		spew.Dump(ctx)

		//todo
		in := request.(*CirclePrizeIn)
		return svc.Exec(ctx, in)
	}
}

//个人实现,参数不能修改
func (r *CirclePrizeHandler) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(CirclePrizeIn), ctx, req)
}

//个人实现,参数不能修改
func (r *CirclePrizeHandler) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response ykit.Result
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *CirclePrizeHandler) HandlerLocal(service CirclePrizeService,
	mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {

	ep := r.MakeLocalEndpoint(service)
	for _, f := range mid {
		ep = f(ep)
	}
	before := tran.ServerBefore(ykit.Jwt2ctx())

	opts := make([]tran.ServerOption, 0)
	opts = append(opts, before)
	opts = append(opts, options...)

	handler := tran.NewServer(
		ep,
		r.DecodeRequest,
		r.base.EncodeResponse,
		opts...)

	return handler
}

//sd,proxy实现,用于etcd自动服务发现时的handler
func (r *CirclePrizeHandler) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		MSTAG,
		//todo
		"POST",
		CirclePrize_HANDLER_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}
