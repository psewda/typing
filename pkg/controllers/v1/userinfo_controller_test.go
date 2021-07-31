package v1_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/psewda/typing/mocks"
	ctrlv1 "github.com/psewda/typing/pkg/controllers/v1"
	"github.com/psewda/typing/pkg/signin/userinfo"
)

var _ = Describe("userinfo controller", func() {
	Context("get userinfo data", func() {
		It("should return valid userinfo when right setup", func() {
			mockContainer := mocks.NewMockContainer(mockCtrl)
			mockUserinfo := mocks.NewMockUserinfo(mockCtrl)
			rec := httptest.NewRecorder()
			mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockUserinfo, nil)
			user := &userinfo.User{
				ID:      "112295411320093",
				Name:    "username",
				Email:   "email@mail.com",
				Picture: "https://lh3.googleusercontent.com/AOh14GiShoGb1kvP=q01-b",
			}
			mockUserinfo.EXPECT().Get().Return(user, nil)
			req := httptest.NewRequest(http.MethodGet, uiRoute, nil)
			ctx := newCtx(req, rec, withAccessToken())
			ctrlv1.NewUserinfoController(mockContainer).GetUser(ctx)
			Expect(rec.Code).Should(Equal(http.StatusOK))

			var u userinfo.User
			json.NewDecoder(rec.Body).Decode(&u)
			Expect(u).ShouldNot(BeZero())
			Expect(u.Name).Should(Equal("username"))
			Expect(u.Email).Should(Equal("email@mail.com"))
		})

		It("should return error when inner error", func() {
			mockContainer := mocks.NewMockContainer(mockCtrl)
			mockUserinfo := mocks.NewMockUserinfo(mockCtrl)
			rec := httptest.NewRecorder()
			mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockUserinfo, nil)
			mockUserinfo.EXPECT().Get().Return(nil, errors.New("error"))

			req := httptest.NewRequest(http.MethodGet, uiRoute, nil)
			ctx := newCtx(req, rec, withAccessToken())
			err := ctrlv1.NewUserinfoController(mockContainer).GetUser(ctx)
			httpError := toHTTPError(err)
			Expect(httpError).Should(HaveOccurred())
			Expect(httpError.Code).Should(Equal(http.StatusInternalServerError))
		})
	})
})
