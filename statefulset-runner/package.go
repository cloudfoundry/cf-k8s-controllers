package statefulsetrunner

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o fake -fake-name Client sigs.k8s.io/controller-runtime/pkg/client.Client
//counterfeiter:generate -o fake -fake-name EventRecorder k8s.io/client-go/tools/record.EventRecorder
//counterfeiter:generate -o fake -fake-name StatusWriter sigs.k8s.io/controller-runtime/pkg/client.StatusWriter
