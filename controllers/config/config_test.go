package config_test

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"

	"code.cloudfoundry.org/korifi/controllers/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LoadFromPath", func() {
	var (
		configPath string
		retConfig  *config.ControllerConfig
		retErr     error
	)

	BeforeEach(func() {
		// Setup filesystem
		var err error
		configPath, err = os.MkdirTemp("", "config")
		Expect(err).NotTo(HaveOccurred())

		config := config.ControllerConfig{
			CFProcessDefaults: config.CFProcessDefaults{
				MemoryMB:    1024,
				DiskQuotaMB: 512,
			},
			CFRootNamespace:             "rootNamespace",
			PackageRegistrySecretName:   "packageRegistrySecretName",
			TaskTTL:                     "taskTTL",
			WorkloadsTLSSecretName:      "workloadsTLSSecretName",
			WorkloadsTLSSecretNamespace: "workloadsTLSSecretNamespace",
			BuilderName:                 "buildReconciler",
			RunnerName:                  "statefulset-runner",
		}
		configYAML, err := yaml.Marshal(config)
		Expect(err).NotTo(HaveOccurred())

		Expect(os.WriteFile(filepath.Join(configPath, "file1"), configYAML, 0o644)).To(Succeed())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(configPath)).To(Succeed())
	})

	JustBeforeEach(func() {
		retConfig, retErr = config.LoadFromPath(configPath)
	})

	It("loads the configuration from all the files in the given directory", func() {
		Expect(retErr).NotTo(HaveOccurred())
		Expect(*retConfig).To(Equal(config.ControllerConfig{
			CFProcessDefaults: config.CFProcessDefaults{
				MemoryMB:    1024,
				DiskQuotaMB: 512,
			},
			CFRootNamespace:             "rootNamespace",
			PackageRegistrySecretName:   "packageRegistrySecretName",
			TaskTTL:                     "taskTTL",
			WorkloadsTLSSecretName:      "workloadsTLSSecretName",
			WorkloadsTLSSecretNamespace: "workloadsTLSSecretNamespace",
			BuilderName:                 "buildReconciler",
			RunnerName:                  "statefulset-runner",
		}))
	})
})

var _ = Describe("ParseTaskTTL", func() {
	var (
		taskTTLString string
		taskTTL       time.Duration
		parseErr      error
	)

	BeforeEach(func() {
		taskTTLString = ""
	})

	JustBeforeEach(func() {
		cfg := config.ControllerConfig{
			TaskTTL: taskTTLString,
		}

		taskTTL, parseErr = cfg.ParseTaskTTL()
	})

	It("return 30 days by default", func() {
		Expect(parseErr).NotTo(HaveOccurred())
		Expect(taskTTL).To(Equal(30 * 24 * time.Hour))
	})

	When("entering something parseable by tools.ParseDuration", func() {
		BeforeEach(func() {
			taskTTLString = "1d12h30m5s20ns"
		})

		It("parses ok", func() {
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(taskTTL).To(Equal(24*time.Hour + 12*time.Hour + 30*time.Minute + 5*time.Second + 20*time.Nanosecond))
		})
	})

	When("entering something that cannot be parsed", func() {
		BeforeEach(func() {
			taskTTLString = "foreva"
		})

		It("returns an error", func() {
			Expect(parseErr).To(HaveOccurred())
		})
	})
})
