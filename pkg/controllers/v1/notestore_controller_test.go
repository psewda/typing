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
	"github.com/psewda/typing/pkg/storage/notestore"
)

var _ = Describe("notestore controller", func() {
	var (
		mockContainer *mocks.MockContainer
		mockNotestore *mocks.MockNotestore
		rec           *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		mockContainer = mocks.NewMockContainer(mockCtrl)
		mockNotestore = mocks.NewMockNotestore(mockCtrl)
		rec = httptest.NewRecorder()
	})

	Context("create new note", func() {
		It("should create the note successfully", func() {
			newReq := func(j string) *http.Request {
				reader := strings.NewReader(j)
				req := httptest.NewRequest(http.MethodPost, notesRoute, reader)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				return req
			}

			By("correct input")
			{
				note := &notestore.Note{
					ID:   "n0hd6hd12tes4",
					Name: "note",
				}
				mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockNotestore, nil)
				mockNotestore.EXPECT().Create(gomock.Any()).Return(note, nil)
				req := newReq(`{"name": "note"}`)
				ctx := newCtx(req, rec, withContainer(mockContainer), withAccessToken())
				ctx.SetPath(notesRoute)

				ctrlv1.CreateNote(ctx)
				Expect(rec.Code).Should(Equal(http.StatusCreated))

				loc := rec.Header().Get(echo.HeaderLocation)
				Expect(loc).Should(HavePrefix(notesRoute))

				var n notestore.Note
				json.NewDecoder(rec.Body).Decode(&n)
				Expect(n).ShouldNot(BeZero())
				Expect(n.Name).Should(Equal("note"))
			}

			By("wrong input")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockNotestore, nil)
				req := newReq(`{"name": ""}`)
				ctx := newCtx(req, rec, withContainer(mockContainer), withAccessToken())

				err := ctrlv1.CreateNote(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusBadRequest))
			}

			By("inner error")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockNotestore, nil)
				mockNotestore.EXPECT().Create(gomock.Any()).Return(nil, errors.New("error"))
				req := newReq(`{"name": "note"}`)
				ctx := newCtx(req, rec, withContainer(mockContainer), withAccessToken())

				err := ctrlv1.CreateNote(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusInternalServerError))
			}
		})
	})

	Context("get all notes", func() {
		It("should return a list off all notes", func() {
			By("correct setup")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockNotestore, nil)
				notes := []*notestore.Note{
					{
						ID:   "hftg5wgs5dfs7",
						Name: "note1",
					},
					{
						ID:   "n0hd6hd12tes4",
						Name: "note2",
					},
				}
				mockNotestore.EXPECT().GetAll().Return(notes, nil)
				req := httptest.NewRequest(http.MethodGet, notesRoute, nil)
				ctx := newCtx(req, rec, withContainer(mockContainer), withAccessToken())

				ctrlv1.GetNotes(ctx)
				Expect(rec.Code).Should(Equal(http.StatusOK))

				var n []*notestore.Note
				json.NewDecoder(rec.Body).Decode(&n)
				Expect(n).ShouldNot(BeNil())
				Expect(len(n)).Should(Equal(2))
				Expect(n[0].Name).Should(Equal("note1"))
				Expect(n[1].Name).Should(Equal("note2"))
			}

			By("inner error")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockNotestore, nil)
				mockNotestore.EXPECT().GetAll().Return(nil, errors.New("error"))
				req := httptest.NewRequest(http.MethodGet, notesRoute, nil)
				ctx := newCtx(req, rec, withContainer(mockContainer), withAccessToken())

				err := ctrlv1.GetNotes(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusInternalServerError))
			}
		})
	})

	Context("get note by id", func() {
		It("should return the note associated with the id", func() {
			By("correct id")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockNotestore, nil)
				note := &notestore.Note{
					ID:          "hftg5wgs5dfs7",
					Name:        "note",
					Description: "desc",
					Labels:      []string{"label1", "label2"},
				}
				mockNotestore.EXPECT().Get(gomock.Any()).Return(note, nil)
				req := httptest.NewRequest(http.MethodGet, noteRouteWithID, nil)
				ctx := newCtx(req, rec, withContainer(mockContainer), withAccessToken())

				ctrlv1.GetNote(ctx)
				Expect(rec.Code).Should(Equal(http.StatusOK))

				var n notestore.Note
				json.NewDecoder(rec.Body).Decode(&n)
				Expect(n).ShouldNot(BeNil())
				Expect(n.Name).Should(Equal("note"))
				Expect(n.Description).Should(Equal("desc"))
			}

			By("wrong id")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockNotestore, nil)
				mockNotestore.EXPECT().Get(gomock.Any()).Return(nil, nil)
				req := httptest.NewRequest(http.MethodGet, noteRouteWithID, nil)
				ctx := newCtx(req, rec, withContainer(mockContainer), withAccessToken())

				err := ctrlv1.GetNote(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusNotFound))
			}

			By("inner error")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockNotestore, nil)
				mockNotestore.EXPECT().Get(gomock.Any()).Return(nil, errors.New("error"))
				req := httptest.NewRequest(http.MethodGet, noteRouteWithID, nil)
				ctx := newCtx(req, rec, withContainer(mockContainer), withAccessToken())

				err := ctrlv1.GetNote(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusInternalServerError))
			}
		})
	})

	Context("update note", func() {
		It("should update the note successfully", func() {
			newReq := func(j string) *http.Request {
				reader := strings.NewReader(j)
				req := httptest.NewRequest(http.MethodPut, noteRouteWithID, reader)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				return req
			}

			By("correct input")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockNotestore, nil)
				note := &notestore.Note{
					ID:          "hftg5wgs5dfs7",
					Name:        "note",
					Description: "desc",
					Labels:      []string{"label1", "label2"},
				}
				mockNotestore.EXPECT().Update(gomock.Any(), gomock.Any()).Return(note, nil)
				req := newReq(`{"name": "note"}`)
				ctx := newCtx(req, rec, withContainer(mockContainer), withAccessToken())

				ctrlv1.UpdateNote(ctx)
				Expect(rec.Code).Should(Equal(http.StatusOK))

				var n notestore.Note
				json.NewDecoder(rec.Body).Decode(&n)
				Expect(n).ShouldNot(BeNil())
				Expect(n.Name).Should(Equal("note"))
				Expect(n.Description).Should(Equal("desc"))
			}

			By("wrong note id")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockNotestore, nil)
				mockNotestore.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, nil)
				req := newReq(`{"name": "note"}`)
				ctx := newCtx(req, rec, withContainer(mockContainer), withAccessToken())

				err := ctrlv1.UpdateNote(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusNotFound))
			}

			By("inner error")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockNotestore, nil)
				mockNotestore.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				req := newReq(`{"name": "note"}`)
				ctx := newCtx(req, rec, withContainer(mockContainer), withAccessToken())

				err := ctrlv1.UpdateNote(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusInternalServerError))
			}
		})
	})

	Context("delete note", func() {
		It("should delete the note successfully", func() {
			By("correct note id")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockNotestore, nil)
				mockNotestore.EXPECT().Delete(gomock.Any()).Return(true, nil)
				req := httptest.NewRequest(http.MethodDelete, noteRouteWithID, nil)
				ctx := newCtx(req, rec, withContainer(mockContainer), withAccessToken())
				ctrlv1.DeleteNote(ctx)
				Expect(rec.Code).Should(Equal(http.StatusNoContent))
			}

			By("wrong note id")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockNotestore, nil)
				mockNotestore.EXPECT().Delete(gomock.Any()).Return(false, nil)
				req := httptest.NewRequest(http.MethodDelete, noteRouteWithID, nil)
				ctx := newCtx(req, rec, withContainer(mockContainer), withAccessToken())

				err := ctrlv1.DeleteNote(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusNotFound))
			}

			By("inner error")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any(), gomock.Any()).Return(mockNotestore, nil)
				mockNotestore.EXPECT().Delete(gomock.Any()).Return(false, errors.New("error"))
				req := httptest.NewRequest(http.MethodDelete, noteRouteWithID, nil)
				ctx := newCtx(req, rec, withContainer(mockContainer), withAccessToken())

				err := ctrlv1.DeleteNote(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusInternalServerError))
			}
		})
	})
})
