package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"
)

func main() {
	appsv1.AddToScheme(scheme.Scheme)
	decode := scheme.Codecs.UniversalDeserializer().Decode
	stream, err := ioutil.ReadFile("/Users/i2tsuki/Repo/nulab/backlog-k8s-clusters/clusters/dev/com-dev/backlog-web/kustomize/deployment.yaml")
	if err != nil {
		log.Printf("failed to read the file.\n")
	}
	obj, gKV, err := decode(stream, nil, nil)
	if err != nil {
		log.Printf(`failed to decode the manifest file.\n`)
	}
	fmt.Printf("obj: %v, gKV: %v\n", obj, gKV)
	if gKV.Kind == "Deployment" {
		d := obj.(*appsv1.Deployment)
		for _, container := range d.Spec.Template.Spec.Containers {
			if container.Name == "backlog-web" {
				cpu := container.Resources.Requests.Cpu()
				cpu.SetMilli(2)
				container.Resources.Requests["cpu"] = *cpu
			}
		}
		fmt.Print(os.Getwd())
		w, _ := os.OpenFile("deployment.yaml", os.O_RDWR|os.O_CREATE, 0644)
		b, _ := yaml.Marshal(&d)
		w.Write(b)
	} else {
		log.Printf(`Group Kind in the manifest file must be "Deployment".\n`)
	}
}
