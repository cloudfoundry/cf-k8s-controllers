package payloads_test

import (
	. "code.cloudfoundry.org/korifi/api/payloads"
	"code.cloudfoundry.org/korifi/api/repositories"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("ManifestApplicationProcess", func() {
	const spaceGUID = "the-space-guid"

	Describe("ToProcessCreateMessage", func() {
		const appGUID = "the-app-guid"
		var processInfo ManifestApplicationProcess

		When("all fields are specified", func() {
			BeforeEach(func() {
				processInfo = ManifestApplicationProcess{
					Type:                         "web",
					Command:                      stringPointer("start-web.sh"),
					DiskQuota:                    stringPointer("512M"),
					HealthCheckHTTPEndpoint:      stringPointer("/stuff"),
					HealthCheckInvocationTimeout: int64Pointer(90),
					HealthCheckType:              stringPointer("http"),
					Instances:                    intPointer(3),
					Memory:                       stringPointer("1G"),
					Timeout:                      int64Pointer(60),
				}
			})

			It("returns a CreateProcessMessage with those values", func() {
				message := processInfo.ToProcessCreateMessage(appGUID, spaceGUID)

				Expect(message).To(Equal(repositories.CreateProcessMessage{
					AppGUID:     appGUID,
					SpaceGUID:   spaceGUID,
					Type:        "web",
					Command:     "start-web.sh",
					DiskQuotaMB: 512,
					HealthCheck: repositories.HealthCheck{
						Type: "http",
						Data: repositories.HealthCheckData{
							HTTPEndpoint:             "/stuff",
							TimeoutSeconds:           60,
							InvocationTimeoutSeconds: 90,
						},
					},
					DesiredInstances: 3,
					MemoryMB:         1024,
				}))
			})

			Describe("HealthCheckType", func() {
				When("HealthCheckType is 'none' (legacy alias for 'process')", func() {
					const noneHealthCheckType = "none"

					It("converts the type to 'process'", func() {
						processInfo.HealthCheckType = stringPointer(noneHealthCheckType)

						message := processInfo.ToProcessCreateMessage(appGUID, spaceGUID)

						Expect(message.HealthCheck.Type).To(Equal("process"))
					})
				})

				When("HealthCheckType is specified as some other valid type", func() {
					It("passes the type through to the message", func() {
						processInfo.HealthCheckType = stringPointer("port")

						message := processInfo.ToProcessCreateMessage(appGUID, spaceGUID)

						Expect(message.HealthCheck.Type).To(Equal("port"))
					})
				})
			})
		})

		When("only type is specified", func() {
			BeforeEach(func() {
				processInfo = ManifestApplicationProcess{}
			})

			When(`type is "web"`, func() {
				BeforeEach(func() {
					processInfo.Type = "web"
				})

				It("returns a CreateProcessMessage with defaulted values", func() {
					message := processInfo.ToProcessCreateMessage(appGUID, spaceGUID)

					Expect(message).To(Equal(repositories.CreateProcessMessage{
						Type:             "web",
						AppGUID:          appGUID,
						SpaceGUID:        spaceGUID,
						DesiredInstances: 1,
						Command:          "",
						DiskQuotaMB:      1024,
						HealthCheck: repositories.HealthCheck{
							Type: "process",
							Data: repositories.HealthCheckData{
								HTTPEndpoint:             "",
								TimeoutSeconds:           0, // this isn't nullable
								InvocationTimeoutSeconds: 0, // this isn't nullable
							},
						},
						MemoryMB: 1024,
					}))
				})
			})

			When(`type is not "web"`, func() {
				BeforeEach(func() {
					processInfo.Type = "worker"
				})

				It("returns a CreateProcessMessage with defaulted values", func() {
					message := processInfo.ToProcessCreateMessage(appGUID, spaceGUID)

					Expect(message).To(Equal(repositories.CreateProcessMessage{
						Type:             "worker",
						AppGUID:          appGUID,
						SpaceGUID:        spaceGUID,
						DesiredInstances: 0,
						Command:          "",
						DiskQuotaMB:      1024,
						HealthCheck: repositories.HealthCheck{
							Type: "process",
							Data: repositories.HealthCheckData{
								HTTPEndpoint:             "",
								TimeoutSeconds:           0, // this isn't nullable
								InvocationTimeoutSeconds: 0, // this isn't nullable
							},
						},
						MemoryMB: 1024,
					}))
				})
			})
		})
	})

	Describe("ToProcessPatchMessage", func() {
		const processGUID = "the-process-guid"
		var processInfo ManifestApplicationProcess

		BeforeEach(func() {
			processInfo = ManifestApplicationProcess{Type: "web"}
		})

		Describe("HealthCheckType", func() {
			When("HealthCheckType is specified as 'none' (legacy alias for 'process')", func() {
				const noneHealthCheckType = "none"

				It("converts the type to 'process'", func() {
					processInfo.HealthCheckType = stringPointer(noneHealthCheckType)

					message := processInfo.ToProcessPatchMessage(processGUID, spaceGUID)

					Expect(message.HealthCheckType).To(Equal(stringPointer("process")))
				})
			})

			When("HealthCheckType is specified as some other valid type", func() {
				It("passes the type through to the message", func() {
					processInfo.HealthCheckType = stringPointer("port")

					message := processInfo.ToProcessPatchMessage(processGUID, spaceGUID)

					Expect(message.HealthCheckType).To(Equal(stringPointer("port")))
				})
			})

			When("HealthCheckType is unspecified", func() {
				It("returns a message with HealthCheckType unset", func() {
					Expect(
						processInfo.ToProcessPatchMessage(processGUID, spaceGUID).HealthCheckType,
					).To(BeNil())
				})
			})
		})

		When("DiskQuota is specified", func() {
			BeforeEach(func() {
				processInfo.DiskQuota = stringPointer("1G")
			})

			It("returns a message with DiskQuotaMB set to the parsed value", func() {
				Expect(
					processInfo.ToProcessPatchMessage(processGUID, spaceGUID).DiskQuotaMB,
				).To(PointTo(BeEquivalentTo(1024)))
			})
		})

		When("DiskQuota is unspecified", func() {
			It("returns a message with DiskQuotaMB unset", func() {
				Expect(
					processInfo.ToProcessPatchMessage(processGUID, spaceGUID).DiskQuotaMB,
				).To(BeNil())
			})
		})

		When("Memory is specified", func() {
			BeforeEach(func() {
				processInfo.Memory = stringPointer("1G")
			})

			It("returns a message with MemoryMB set to the parsed value", func() {
				Expect(
					processInfo.ToProcessPatchMessage(processGUID, spaceGUID).MemoryMB,
				).To(PointTo(BeEquivalentTo(1024)))
			})
		})

		When("Memory is unspecified", func() {
			It("returns a message with MemoryMB unset", func() {
				Expect(
					processInfo.ToProcessPatchMessage(processGUID, spaceGUID).MemoryMB,
				).To(BeNil())
			})
		})

		When("Instances is specified", func() {
			BeforeEach(func() {
				processInfo.Instances = intPointer(3)
			})

			It("returns a message with DesiredInstances set to the parsed value", func() {
				Expect(
					processInfo.ToProcessPatchMessage(processGUID, spaceGUID).DesiredInstances,
				).To(PointTo(BeEquivalentTo(3)))
			})
		})

		When("Instances is unspecified", func() {
			It("returns a message with DesiredInstances unset", func() {
				Expect(
					processInfo.ToProcessPatchMessage(processGUID, spaceGUID).DesiredInstances,
				).To(BeNil())
			})
		})
	})
})

func stringPointer(s string) *string {
	return &s
}

func intPointer(i int) *int {
	return &i
}

func int64Pointer(i int64) *int64 {
	return &i
}
