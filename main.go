package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

func main() {
	filename := "stateful.yml"
	// read
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("no file")
		return
	}
	buf_slice := bytes.Split(buf, []byte("---"))
	decoder := scheme.Codecs.UniversalDeserializer()
	var decoded []runtime.Object
	for _, b := range buf_slice {
		obj, _, err := decoder.Decode(b, nil, nil)
		if err == nil && obj != nil {
			decoded = append(decoded, obj)
		}
	}
	// lets pick only 1 element
	decode := decoded[0]

	// write
	info, _ := runtime.SerializerInfoForMediaType(scheme.Codecs.SupportedMediaTypes(), "application/yaml")
	groupVersion := schema.GroupVersion{Group: decode.GetObjectKind().GroupVersionKind().Group, Version: decode.GetObjectKind().GroupVersionKind().Version}
	encoder := scheme.Codecs.EncoderForVersion(info.Serializer, groupVersion)
	yaml, err := runtime.Encode(encoder, decode)
	if err != nil {
		fmt.Println("encode error")
		return
	}
	err = ioutil.WriteFile("after-roundtrip.yml", yaml, 0644)
	if err != nil {
		fmt.Println("write error")
	}
	return
}
