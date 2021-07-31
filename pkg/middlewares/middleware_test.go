package middlewares_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/psewda/typing/pkg/middlewares"
)

func TestMiddlewares(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "middlewares-suite")
}

var _ = Describe("middleware", func() {
	Context("authorization middleware", func() {
		It("should succeed when valid token", func() {
			ctx := newCtx()
			ctx.Request().Header.Add(echo.HeaderAuthorization, "Bearer qws3m55ydh5s80sb")
			middleware := middlewares.Authorization()
			handler := middleware(func(ctx echo.Context) error { return nil })
			Expect(handler(ctx)).Should(Succeed())
			Expect(ctx.Get(middlewares.ContextKeyAccessToken)).ShouldNot(BeZero())
		})

		It("should throw error when empty token", func() {
			ctx := newCtx()
			middleware := middlewares.Authorization()
			handler := middleware(func(ctx echo.Context) error { return nil })
			err := handler(ctx)
			httpError := toHTTPError(err)
			Expect(httpError).Should(HaveOccurred())
			Expect(httpError.Code).Should(Equal(http.StatusUnauthorized))
		})

		It("should throw error when malformed token", func() {
			ctx := newCtx()
			ctx.Request().Header.Add(echo.HeaderAuthorization, "malformed")
			middleware := middlewares.Authorization()
			handler := middleware(func(ctx echo.Context) error { return nil })
			err := handler(ctx)
			httpError := toHTTPError(err)

			Expect(httpError).Should(HaveOccurred())
			Expect(httpError.Code).Should(Equal(http.StatusUnauthorized))
		})
	})
})

func newCtx() echo.Context {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/any-url", nil)
	ctx := echo.New().NewContext(req, rec)
	ctx.Logger().SetOutput(ioutil.Discard)
	return ctx
}

func toHTTPError(err error) *echo.HTTPError {
	if err != nil {
		return err.(*echo.HTTPError)
	}
	return nil
}
