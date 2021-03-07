package v1_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/psewda/typing/pkg/di"
	"github.com/psewda/typing/pkg/middlewares"
)

const (
	urlRoute     = "/api/v1/signin/auth/url"
	tokenRoute   = "/api/v1/signin/auth/token"
	refreshRoute = "/api/v1/signin/auth/refresh"
	revokeRoute  = "/api/v1/signin/auth/revoke"
)

var mockCtrl *gomock.Controller

func TestControllersV1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "controllers-v1-suite")
}

var _ = BeforeSuite(func() {
	mockCtrl = gomock.NewController(GinkgoT())
})

var _ = AfterSuite(func() {
	mockCtrl.Finish()
})

type keyValue struct {
	key   string
	value interface{}
}

func newCtx(r *http.Request, w http.ResponseWriter, objs ...keyValue) echo.Context {
	ctx := echo.New().NewContext(r, w)
	for _, obj := range objs {
		ctx.Set(obj.key, obj.value)
	}
	ctx.Logger().SetOutput(ioutil.Discard)
	return ctx
}

func withContainer(v di.Container) keyValue {
	return keyValue{
		key:   middlewares.KeyContainer,
		value: v,
	}
}

func toHTTPError(err error) *echo.HTTPError {
	if err != nil {
		return err.(*echo.HTTPError)
	}
	return nil
}
