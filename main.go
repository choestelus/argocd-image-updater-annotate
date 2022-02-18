package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

const (
	kimage = "argocd-image-updater.argoproj.io/%s.kustomize.image-name"
	kstrat = "argocd-image-updater.argoproj.io/%s.update-strategy"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("args needed")
	}
	rawK, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	resources, err := kio.FromBytes(rawK)
	if err != nil {
		panic(err)
	}
	resource := resources[0]
	imgs, err := resource.Pipe(yaml.Get("images"))
	if err != nil {
		panic(err)
	}

	services, err := parseImageNodesName(imgs)
	if err != nil {
		panic(err)
	}
	services = extractServiceName(services...)
	fmt.Println("services: ", services)

	for _, svc := range services {
		if err := resource.PipeE(yaml.SetAnnotation(annotateImage(svc), fmt.Sprintf("%s-image", svc))); err != nil {
			panic(err)
		}

		if err := resource.PipeE(yaml.SetAnnotation(annotateStrategy(svc), "latest")); err != nil {
			panic(err)
		}
	}
}

func annotateStrategy(service string) string {
	return fmt.Sprintf(kstrat, service)
}

func annotateImage(service string) string {
	return fmt.Sprintf(kimage, service)
}

func parseImageNodes(r *yaml.RNode) ([][]string, error) {
	result, err := r.ElementValuesList([]string{"name", "newName", "newTag"})
	if err != nil {
		return nil, fmt.Errorf("failed to get key [name, newName, newTag] from resource node: %w", err)
	}
	return result, err
}

func parseImageNodesName(r *yaml.RNode) ([]string, error) {
	result, err := r.ElementValues("name")
	if err != nil {
		return nil, fmt.Errorf("failed to get key [name]: %w", err)
	}
	return result, nil
}

func extractServiceName(images ...string) []string {

	trimmedName := []string{}
	for _, img := range images {
		trimmedName = append(trimmedName, strings.TrimSuffix(img, "-image"))
	}
	return trimmedName
}
