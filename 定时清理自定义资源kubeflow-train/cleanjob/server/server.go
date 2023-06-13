package server

import (
	"cleanjob/client"
	trainjobs "cleanjob/model"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"time"
)

// 定义过期时间
const CleanTime = time.Second * 48

var dyClient = client.NewClient()

func CleanTfJob(ns string) {

	//要操作的资源对象
	gvr := schema.GroupVersionResource{
		Group:    "kubeflow.org",
		Version:  "v1",
		Resource: "tfjobs",
	}

	unstructObjList := getObjectList(gvr, ns)
	//引用原生的会报错
	//tf := &tfv1.TFJobList{}
	result := &trainjobs.TFJobList{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObjList.UnstructuredContent(), result)
	if err != nil {
		panic(err.Error())
	}
	for _, item := range result.Items {
		v := item.Status.CompletionTime
		if v != nil {
			t := time.Now()
			num := t.Sub(v.Time)
			if num > CleanTime {
				err = deleteObject(gvr, ns, item.Name)
				if err != nil {
					panic(err.Error())
				}
				fmt.Printf("成功删除命名空间%s下的TfJobs对象%s !", ns, item.Name)
			}
		}
	}

}
func PaddleJob(ns string) {

	//要操作的资源对象
	gvr := schema.GroupVersionResource{
		Group:    "kubeflow.org",
		Version:  "v1",
		Resource: "paddlejobs",
	}

	unstructObjList := getObjectList(gvr, ns)
	//引用原生的会报错
	//tf := &tfv1.TFJobList{}
	result := &trainjobs.PaddleJobList{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObjList.UnstructuredContent(), result)
	if err != nil {
		panic(err.Error())
	}
	for _, item := range result.Items {
		v := item.Status.CompletionTime
		if v != nil {
			t := time.Now()
			num := t.Sub(v.Time)
			if num > CleanTime {
				err = deleteObject(gvr, ns, item.Name)
				if err != nil {
					panic(err.Error())
				}
				fmt.Printf("成功删除命名空间%s下的paddlejob对象%s !", ns, item.Name)
			}
		}
	}

}
func PyTorchJob(ns string) {

	//要操作的资源对象
	gvr := schema.GroupVersionResource{
		Group:    "kubeflow.org",
		Version:  "v1",
		Resource: "pytorchjobs",
	}

	unstructObjList := getObjectList(gvr, ns)
	//引用原生的会报错
	//tf := &tfv1.TFJobList{}
	result := &trainjobs.TFJobList{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObjList.UnstructuredContent(), result)
	if err != nil {
		panic(err.Error())
	}
	for _, item := range result.Items {
		v := item.Status.CompletionTime
		if v != nil {
			t := time.Now()
			num := t.Sub(v.Time)
			if num > CleanTime {
				err = deleteObject(gvr, ns, item.Name)
				if err != nil {
					panic(err.Error())
				}
				fmt.Printf("成功删除命名空间%s下的pytorchjobs对象%s !", ns, item.Name)
			}
		}
	}

}
func MSJob(ns string) {

	//要操作的资源对象
	gvr := schema.GroupVersionResource{
		Group:    "mindspore.gitee.com",
		Version:  "v1",
		Resource: "msjobs",
	}

	unstructObjList := getObjectList(gvr, ns)
	//引用原生的会报错
	//tf := &tfv1.TFJobList{}
	result := &trainjobs.MSJobList{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObjList.UnstructuredContent(), result)
	if err != nil {
		panic(err.Error())
	}
	for _, item := range result.Items {
		v := item.Status.CompletionTime
		if v != nil {
			t := time.Now()
			num := t.Sub(v.Time)
			if num > CleanTime {
				err = deleteObject(gvr, ns, item.Name)
				if err != nil {
					panic(err.Error())
				}
				fmt.Printf("成功删除命名空间%s下的MSJob对象%s !", ns, item.Name)
			}
		}
	}

}

func GetNamespaceList() *v1.NamespaceList {
	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "namespaces",
	}
	unstructObjList, err := dyClient.
		//Resource是dynamicClient唯一的一个方法，参数为gvr
		Resource(gvr).
		//以list列表的方式查询
		List(context.TODO(), metav1.ListOptions{Limit: 100})

	result := &v1.NamespaceList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObjList.UnstructuredContent(), result)
	if err != nil {
		panic(err.Error())
	}
	return result
}

func deleteObject(gvr schema.GroupVersionResource, ns string, name string) error {
	return dyClient.Resource(gvr).Namespace(ns).Delete(context.TODO(), name, *metav1.NewDeleteOptions(0))

}
func getObjectList(gvr schema.GroupVersionResource, ns string) *unstructured.UnstructuredList {
	unstructObjList, err := dyClient.
		//Resource是dynamicClient唯一的一个方法，参数为gvr
		Resource(gvr).
		//指定查询的namespace
		Namespace(ns).
		//以list列表的方式查询
		List(context.TODO(), metav1.ListOptions{Limit: 100})

	if err != nil {
		fmt.Println("dynamicClient list pods failed ! err :", err)
		panic(err.Error())
	}
	return unstructObjList
}
