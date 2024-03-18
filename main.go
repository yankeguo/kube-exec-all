package main

import (
	"bytes"
	"context"
	"flag"
	"log"
	"os"

	"github.com/yankeguo/rg"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

var (
	optScript string
	optShell  string
)

func main() {
	var err error
	defer func() {
		if err == nil {
			return
		}
		log.Println("exited with error:", err.Error())
		os.Exit(1)
	}()
	defer rg.Guard(&err)

	flag.StringVar(&optScript, "script", "script.sh", "the script file to run")
	flag.StringVar(&optShell, "shell", "/bin/sh", "the shell to use")
	flag.Parse()

	buf := rg.Must(os.ReadFile(optScript))

	cfg := rg.Must(rest.InClusterConfig())

	client := rg.Must(kubernetes.NewForConfig(cfg))

	ctx := context.Background()

	list := rg.Must(client.CoreV1().Pods("").List(ctx, metav1.ListOptions{}))

	for _, item := range list.Items {
		for _, container := range item.Spec.Containers {
			execute(ctx, executeOptions{
				script:    buf,
				config:    cfg,
				client:    client,
				shell:     optShell,
				namespace: item.Namespace,
				name:      item.Name,
				container: container.Name,
			})
		}
	}
}

type executeOptions struct {
	script    []byte
	config    *rest.Config
	client    *kubernetes.Clientset
	shell     string
	namespace string
	name      string
	container string
}

func execute(ctx context.Context, opts executeOptions) {
	var err error
	defer func() {
		if err == nil {
			return
		}
		log.Printf("failed execution (%s/%s:%s): %s", opts.namespace, opts.name, opts.container, err.Error())
	}()
	defer rg.Guard(&err)

	req := opts.client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(opts.name).
		Namespace(opts.namespace).
		SubResource("exec")

	req.VersionedParams(&v1.PodExecOptions{
		Container: opts.container,
		Command:   []string{opts.shell},
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
	}, scheme.ParameterCodec)

	exe := rg.Must(remotecommand.NewSPDYExecutor(opts.config, "POST", req.URL()))

	rg.Must0(exe.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  bytes.NewReader(opts.script),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}))
}
