package msg

//for snippet用于标准返回值的微服务接口

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-kit/kit/endpoint"
	tran "github.com/go-kit/kit/transport/http"
	"github.com/vhaoran/vchat/common/ypage"
	"github.com/vhaoran/vchat/lib/ykit"
)

const (
	//todo
	MsgHis_private_HANDLER_PATH = "/PageMsgHis"
)

type (
	MsgHisPrivateService interface {
		//todo
		Exec(in *MsgHisPrivateIn) (*ypage.PageResult, error)
	}

	//input data
	//todo
	MsgHisPrivateIn struct {
		From int64 `json:"from omitempty"`
		To   int64 `json:"to omitempty"`
		ypage.PageBean
	}

	//output data
	//Result struct {
	//	Code int         `json:"code"`
	//	Msg  string      `json:"msg"`
	//	Data interface{} `json:"data"`
	//}

	// handler implements
	MsgHisPrivateHandler struct {
		base ykit.RootTran
	}
)

func (r *MsgHisPrivateHandler) MakeLocalEndpoint(svc MsgHisPrivateService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  PageMsgHis ###########")
		spew.Dump(ctx)

		//todo
		in := request.(*MsgHisPrivateIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *MsgHisPrivateHandler) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(MsgHisPrivateIn), ctx, req)
}

//个人实现,参数不能修改
func (r *MsgHisPrivateHandler) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *ypage.PageResult
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *MsgHisPrivateHandler) HandlerLocal(service MsgHisPrivateService,
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

	//handler = loggingMiddleware()
	return handler
}

//sd,proxy实现,用于etcd自动服务发现时的handler
func (r *MsgHisPrivateHandler) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		MSTAG,
		//todo
		"POST",
		MsgHis_private_HANDLER_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}
