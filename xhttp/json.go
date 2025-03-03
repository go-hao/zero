package xhttp

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func Json(ctx context.Context, w http.ResponseWriter, v interface{}) {
	httpx.OkJsonCtx(ctx, w, newResponse(v))
}
