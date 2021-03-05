package middlewares_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/psewda/typing/pkg/di"
	"github.com/psewda/typing/pkg/middlewares"
)

func TestMiddlewares(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "middlewares-suite")
}

var _ = Describe("middleware", func() {
	Context("authorization middleware", func() {
		It("should validate authorization token", func() {
			By("valid token")
			{
				ctx := newCtx()
				ctx.Request().Header.Add(echo.HeaderAuthorization, "Bearer qws3m55ydh5s80sb")
				middleware := middlewares.Authorization()
				handler := middleware(func(ctx echo.Context) error { return nil })
				Expect(handler(ctx)).Should(Succeed())
				Expect(ctx.Get(middlewares.KeyAccessToken)).ShouldNot(BeZero())
			}

			By("empty token")
			{
				ctx := newCtx()
				middleware := middlewares.Authorization()
				handler := middleware(func(ctx echo.Context) error { return nil })
				err := handler(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusUnauthorized))
			}

			By("malformed token")
			{
				ctx := newCtx()
				ctx.Request().Header.Add(echo.HeaderAuthorization, "malformed")
				middleware := middlewares.Authorization()
				handler := middleware(func(ctx echo.Context) error { return nil })
				err := handler(ctx)
				httpError := toHTTPError(err)

				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusUnauthorized))
			}
		})
	})

	Context("dependencies middleware", func() {
		It("should inject container instance in the context", func() {
			ctx := newCtx()
			container := di.New()
			middleware := middlewares.Dependencies(container)
			handler := middleware(func(ctx echo.Context) error { return nil })
			Expect(handler(ctx)).Should(Succeed())
			Expect(ctx.Get(middlewares.KeyContainer)).ShouldNot(BeNil())
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