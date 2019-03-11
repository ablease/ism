package commands_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pivotal-cf/ism/commands"
	"github.com/pivotal-cf/ism/commands/commandsfakes"
	"github.com/pivotal-cf/ism/osbapi"
)

var _ = Describe("Platform Register Command", func() {

	var (
		fakeUI                *commandsfakes.FakeUI
		fakePlatformRegistrar *commandsfakes.FakePlatformRegistrar

		registerCommand PlatformRegisterCommand

		executeErr error
	)

	BeforeEach(func() {
		fakeUI = &commandsfakes.FakeUI{}
		fakePlatformRegistrar = &commandsfakes.FakePlatformRegistrar{}

		registerCommand = PlatformRegisterCommand{
			UI:                fakeUI,
			PlatformRegistrar: fakePlatformRegistrar,
		}
	})

	JustBeforeEach(func() {
		executeErr = registerCommand.Execute(nil)
	})

	When("given all required args", func() {
		BeforeEach(func() {
			registerCommand.Name = "platform-1"
			registerCommand.URL = "test-url"
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("displays that the platform was registered", func() {
			text, data := fakeUI.DisplayTextArgsForCall(0)
			Expect(text).To(Equal("Platform '{{.PlatformName}}' registered."))
			Expect(data[0]).To(HaveKeyWithValue("PlatformName", "platform-1"))
		})

		It("registers the platform", func() {
			broker := fakePlatformRegistrar.RegisterArgsForCall(0)

			Expect(broker).To(Equal(&osbapi.Platform{
				Name: "platform-1",
				URL:  "test-url",
			}))
		})

		When("registering the broker errors", func() {
			BeforeEach(func() {
				fakePlatformRegistrar.RegisterReturns(errors.New("error-registering-platform"))
			})

			It("propagates the error", func() {
				Expect(executeErr).To(MatchError("error-registering-platform"))
			})
		})
	})
})
