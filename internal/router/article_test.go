package router

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"xmlt/internal/api/v1"
	"xmlt/internal/domain"
	"xmlt/internal/expand/public"
	"xmlt/internal/service"
)

func TestArticleController_GetByID(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.ArticleService

		req      *http.Request
		wantResp public.Response
	}{
		// TODO: Add test cases.
		{
			name: "查找文章",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/api/v1/article/read/1", nil)
			}(),
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				ms := service.NewMockArticleService(ctrl)
				ms.EXPECT().Get(gomock.Any(), uint64(1)).Return(domain.Article{
					ID:      1,
					Title:   "这是标题",
					Content: "这是内容",
					Utime:   now,
					Ctime:   now,
				}, nil)
				return ms
			},
			wantResp: public.Response{
				Status: 200,
				Data: v1.ArticleVO{
					ID:      1,
					Title:   "这是标题",
					Content: "这是内容",
					Utime:   now.String(),
					Ctime:   now.String(),
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			v1 := server.Group("/api/v1")
			RouterGroupCtrl.ArticleRouter.InitApiRouter(v1)

			// 准备记录结果
			w := httptest.NewRecorder()

			// 模拟执行请求
			server.ServeHTTP(w, tc.req)

			// 解析响应
			resp := public.Response{}
			decoder := json.NewDecoder(w.Body)
			err := decoder.Decode(&resp)
			fmt.Println(err.Error())
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}
