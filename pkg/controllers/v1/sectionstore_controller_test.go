package v1_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/psewda/typing/mocks"
	ctrlv1 "github.com/psewda/typing/pkg/controllers/v1"
	"github.com/psewda/typing/pkg/storage/sectionstore"
)

var _ = Describe("sectionstore controller", func() {
	var (
		mockContainer    *mocks.MockContainer
		mockSectionstore *mocks.MockSectionstore
		rec              *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		mockContainer = mocks.NewMockContainer(mockCtrl)
		mockSectionstore = mocks.NewMockSectionstore(mockCtrl)
		rec = httptest.NewRecorder()
	})

	Context("create new section", func() {
		newReq := func(j string) *http.Request {
			reader := strings.NewReader(j)
			req := httptest.NewRequest(http.MethodPost, sectionsRoute, reader)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			return req
		}

		It("should create the section when correct input", func() {
			section := &sectionstore.Section{
				ID:   "n0hd6hd12tes4",
				Name: "section",
			}
			mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockSectionstore, nil)
			mockSectionstore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(section, nil)
			req := newReq(`{"name": "section"}`)
			ctx := newCtx(req, rec, withAccessToken())

			ctrlv1.NewSectionstoreController(mockContainer).CreateSection(ctx)
			Expect(rec.Code).Should(Equal(http.StatusCreated))

			loc := rec.Header().Get(echo.HeaderLocation)
			Expect(loc).Should(HavePrefix(sectionsRoute))

			var s sectionstore.Section
			json.NewDecoder(rec.Body).Decode(&s)
			Expect(s).ShouldNot(BeZero())
			Expect(s.Name).Should(Equal("section"))
		})

		It("should return error when wrong input", func() {
			mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockSectionstore, nil)
			req := newReq(`{"name": ""}`)
			ctx := newCtx(req, rec, withAccessToken())

			err := ctrlv1.NewSectionstoreController(mockContainer).CreateSection(ctx)
			httpError := toHTTPError(err)
			Expect(httpError).Should(HaveOccurred())
			Expect(httpError.Code).Should(Equal(http.StatusBadRequest))
		})

		It("should return error when inner error", func() {
			mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockSectionstore, nil)
			mockSectionstore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			req := newReq(`{"name": "section"}`)
			ctx := newCtx(req, rec, withAccessToken())

			err := ctrlv1.NewSectionstoreController(mockContainer).CreateSection(ctx)
			httpError := toHTTPError(err)
			Expect(httpError).Should(HaveOccurred())
			Expect(httpError.Code).Should(Equal(http.StatusInternalServerError))
		})
	})

	Context("get all sections", func() {
		It("should return all sections when correct setup", func() {
			mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockSectionstore, nil)
			sections := []*sectionstore.Section{
				{
					ID:   "hftg5wgs5dfs7",
					Name: "section1",
				},
				{
					ID:   "n0hd6hd12tes4",
					Name: "section2",
				},
			}
			mockSectionstore.EXPECT().GetAll(gomock.Any()).Return(sections, nil)
			req := httptest.NewRequest(http.MethodGet, sectionsRoute, nil)
			ctx := newCtx(req, rec, withAccessToken())

			ctrlv1.NewSectionstoreController(mockContainer).GetSections(ctx)
			Expect(rec.Code).Should(Equal(http.StatusOK))

			var s []*sectionstore.Section
			json.NewDecoder(rec.Body).Decode(&s)
			Expect(s).ShouldNot(BeNil())
			Expect(len(s)).Should(Equal(2))
			Expect(s[0].Name).Should(Equal("section1"))
			Expect(s[1].Name).Should(Equal("section2"))
		})

		It("should return error when inner error", func() {
			mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockSectionstore, nil)
			mockSectionstore.EXPECT().GetAll(gomock.Any()).Return(nil, errors.New("error"))
			req := httptest.NewRequest(http.MethodGet, sectionsRoute, nil)
			ctx := newCtx(req, rec, withAccessToken())

			err := ctrlv1.NewSectionstoreController(mockContainer).GetSections(ctx)
			httpError := toHTTPError(err)
			Expect(httpError).Should(HaveOccurred())
			Expect(httpError.Code).Should(Equal(http.StatusInternalServerError))
		})
	})

	Context("get section by id", func() {
		It("should return the section when correct id", func() {
			mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockSectionstore, nil)
			section := &sectionstore.Section{
				ID:   "hftg5wgs5dfs7",
				Name: "section",
			}
			mockSectionstore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(section, nil)
			req := httptest.NewRequest(http.MethodGet, sectionRouteWithID, nil)
			ctx := newCtx(req, rec, withAccessToken())

			ctrlv1.NewSectionstoreController(mockContainer).GetSection(ctx)
			Expect(rec.Code).Should(Equal(http.StatusOK))

			var s sectionstore.Section
			json.NewDecoder(rec.Body).Decode(&s)
			Expect(s).ShouldNot(BeNil())
			Expect(s.Name).Should(Equal("section"))
		})

		It("should return error when inner error", func() {
			mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockSectionstore, nil)
			mockSectionstore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			req := httptest.NewRequest(http.MethodGet, sectionRouteWithID, nil)
			ctx := newCtx(req, rec, withAccessToken())

			err := ctrlv1.NewSectionstoreController(mockContainer).GetSection(ctx)
			httpError := toHTTPError(err)
			Expect(httpError).Should(HaveOccurred())
			Expect(httpError.Code).Should(Equal(http.StatusInternalServerError))
		})
	})

	Context("update section", func() {
		newReq := func(j string) *http.Request {
			reader := strings.NewReader(j)
			req := httptest.NewRequest(http.MethodPut, sectionRouteWithID, reader)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			return req
		}

		It("should succeed when correct input", func() {
			mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockSectionstore, nil)
			section := &sectionstore.Section{
				ID:   "hftg5wgs5dfs7",
				Name: "section",
			}
			mockSectionstore.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(section, nil)
			req := newReq(`{"name": "section"}`)
			ctx := newCtx(req, rec, withAccessToken())

			ctrlv1.NewSectionstoreController(mockContainer).UpdateSection(ctx)
			Expect(rec.Code).Should(Equal(http.StatusOK))

			var s sectionstore.Section
			json.NewDecoder(rec.Body).Decode(&s)
			Expect(s).ShouldNot(BeNil())
			Expect(s.Name).Should(Equal("section"))
		})

		It("should return error when wrong input", func() {
			mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockSectionstore, nil)
			req := newReq(`{"name": ""}`)
			ctx := newCtx(req, rec, withAccessToken())

			err := ctrlv1.NewSectionstoreController(mockContainer).UpdateSection(ctx)
			httpError := toHTTPError(err)
			Expect(httpError).Should(HaveOccurred())
			Expect(httpError.Code).Should(Equal(http.StatusBadRequest))
		})

		It("should return error when inner error", func() {
			mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockSectionstore, nil)
			mockSectionstore.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			req := newReq(`{"name": "section"}`)
			ctx := newCtx(req, rec, withAccessToken())

			err := ctrlv1.NewSectionstoreController(mockContainer).UpdateSection(ctx)
			httpError := toHTTPError(err)
			Expect(httpError).Should(HaveOccurred())
			Expect(httpError.Code).Should(Equal(http.StatusInternalServerError))
		})
	})

	Context("delete note", func() {
		It("should succeed when correct section id", func() {
			mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockSectionstore, nil)
			mockSectionstore.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			req := httptest.NewRequest(http.MethodDelete, sectionRouteWithID, nil)
			ctx := newCtx(req, rec, withAccessToken())
			ctrlv1.NewSectionstoreController(mockContainer).DeleteSection(ctx)
			Expect(rec.Code).Should(Equal(http.StatusNoContent))
		})

		It("should return error when inner error", func() {
			mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockSectionstore, nil)
			mockSectionstore.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("error"))
			req := httptest.NewRequest(http.MethodDelete, sectionRouteWithID, nil)
			ctx := newCtx(req, rec, withAccessToken())

			err := ctrlv1.NewSectionstoreController(mockContainer).DeleteSection(ctx)
			httpError := toHTTPError(err)
			Expect(httpError).Should(HaveOccurred())
			Expect(httpError.Code).Should(Equal(http.StatusInternalServerError))
		})
	})
})
