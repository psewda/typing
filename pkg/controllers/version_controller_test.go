package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/psewda/typing"
	"github.com/psewda/typing/pkg/controllers"
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "controllers-suite")
}

var _ = Describe("version controller", func() {
	Context("get version", func() {
		It("should return correct version string", func() {
			req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
			rec := httptest.NewRecorder()
			ctx := echo.New().NewContext(req, rec)

			controllers.GetVersion(ctx)
			Expect(rec.Code).Should(Equal(http.StatusOK))
			Expect(rec.Body.String()).Should(ContainSubstring(typing.GetVersionString()))
		})
	})
})
