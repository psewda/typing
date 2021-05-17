package googleuserinfo_test

import (
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/errs"
	"github.com/psewda/typing/pkg/signin/userinfo/googleuserinfo"
)

func TestGoogleUserinfo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "googleuserinfo-suite")
}

var _ = Describe("google userinfo", func() {
	Context("get userinfo", func() {
		It("should return userinfo data", func() {
			By("right setup")
			{
				j := `{
					"id": "112295411320093",
					"name": "username",
					"email": "email@mail.com",
					"picture": "https://lh3.googleusercontent.com/AOh14GiShoGb1kvP=q01-b"
				}`
				client := utils.ClientWithJSON(j, http.StatusOK)
				gui, _ := googleuserinfo.New(client)
				u, err := gui.Get()

				Expect(u).ShouldNot(BeNil())
				Expect(u.Name).Should(Equal("username"))
				Expect(u.Email).Should(Equal("email@mail.com"))
				Expect(err).ShouldNot(HaveOccurred())
			}

			By("authorization failure")
			{
				client := utils.ClientWithJSON("{}", http.StatusUnauthorized)
				gui, _ := googleuserinfo.New(client)
				_, err := gui.Get()
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(BeAssignableToTypeOf(errs.NewUnauthorizedError()))
			}

			By("inner error")
			{
				client := utils.ClientWithJSON("error", http.StatusInternalServerError)
				gui, _ := googleuserinfo.New(client)
				_, err := gui.Get()
				Expect(err).Should(HaveOccurred())
			}
		})
	})
})
