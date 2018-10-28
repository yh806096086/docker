package agent

import (

	"net/http"
	"time"
	"log"
	"context"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"fmt"
	"github.com/docker/docker/api/types"
	"encoding/json"
	"io/ioutil"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)


type containerCreateConfig struct {
	Config container.HostConfig
	HostConfig container.HostConfig
	NetworkConfig network.NetworkingConfig
}

var cli *client.Client

func init() {
	c, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	cli = c
}

func testHandler(w http.ResponseWriter, r *http.Request) {


	fmt.Fprint(w, []byte("test ok"))
}

//镜像操作
func listImages(w http.ResponseWriter, r *http.Request) {

	//w.Header().Set("Code", http.StatusOK)

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		fmt.Fprint(w, []byte("list image error"))
	}
	body, err := json.Marshal(images)
	if err != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}
	fmt.Fprint(w, body)
	return
}


func deleteImage(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, []byte("delete container"))
	if ok := r.ParseForm(); ok != nil {
		fmt.Fprint(w, []byte("Server Error"))
	}
	id := r.Form.Get("id")
	if _, ok := cli.ImageRemove(context.Background(), id, types.ImageRemoveOptions{}); ok != nil {
		fmt.Fprint(w, []byte("Server Error"))
	}
	fmt.Fprint(w, []byte("delete image ok"))
	return
}

func pullImage(w http.ResponseWriter, r *http.Request) {
	//.Fprint(w, []byte("pull image"))
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}
	opt := types.ImagePullOptions{}
	if ok := json.Unmarshal(body, opt); ok != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}
	resp, ok := cli.ImagePull(context.Background(), "", opt)
	if ok != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}

	defer resp.Close()
	fmt.Fprint(w, []byte("pull image ok"))
	return
}

func inspectImage(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, []byte("inspect image"))
	if ok := r.ParseForm(); ok != nil {
		fmt.Fprint(w, []byte("Server Error"))
	}
	id := r.Form.Get("id")


	ins, _, ok := cli.ImageInspectWithRaw(context.Background(), id);
	if ok != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}
	d, err := json.Marshal(ins)
	if err != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}
	fmt.Fprint(w, d)
	return
}

//容器操作
func listContainer(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, []byte("containers"))
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}
	d, err := json.Marshal(containers)
	if err != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}
	fmt.Fprint(w, d)
	return
}

func runContainer(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, []byte("run container"))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}

	opt := containerCreateConfig{}
	if ok := json.Unmarshal(body, opt); ok != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}

	fmt.Println(opt)

	//创建新容器
	containerCreateBody, err := cli.ContainerCreate(context.Background(), &container.Config{},
						&container.HostConfig{}, &network.NetworkingConfig{}, "test")

	if err != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}

	// container run
	if ok := cli.ContainerStart(context.Background(), containerCreateBody.ID, types.ContainerStartOptions{}); ok != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}
	fmt.Fprint(w, []byte("create container ok"))
	return


}

func stopContainer(w http.ResponseWriter, r *http.Request) {
	if ok := r.ParseForm(); ok != nil {
		fmt.Fprint(w, []byte("Server Error"))
	}
	id := r.Form.Get("id")

	timeout := 1 *time.Second
	if ok:= cli.ContainerStop(context.Background(), id, &timeout); ok != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}
	fmt.Fprint(w, []byte("stop container"))
	return
}

func stateContainer(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, []byte("state container"))
	if ok := r.ParseForm(); ok != nil {
		fmt.Fprint(w, []byte("Server Error"))
	}
	id := r.Form.Get("id")

	stats, err := cli.ContainerStats(context.Background(), id, false)
	if err != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}
	d, err := json.Marshal(stats)
	if err != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}
	fmt.Fprint(w, d)
	return
}

func inspectContainer(w http.ResponseWriter, r *http.Request) {
	if ok := r.ParseForm(); ok != nil {
		fmt.Fprint(w, []byte("Server Error"))
	}
	id := r.Form.Get("id")

	inspect, err := cli.ContainerInspect(context.Background(), id);
	if err != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}
	d, err := json.Marshal(inspect)
	if err != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}
	fmt.Fprint(w, d)
	return
}

func deleteContainer(w http.ResponseWriter, r *http.Request) {
	if ok := r.ParseForm(); ok != nil {
		fmt.Fprint(w, []byte("Server Error"))
	}
	id := r.Form.Get("id")

	if ok := cli.ContainerRemove(context.Background(),id, types.ContainerRemoveOptions{}); ok != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}
	fmt.Fprint(w, []byte("delete container"))
	return
}

//docker daemon
func dockerVersion(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, []byte("docker version"))
	v, err := cli.ServerVersion(context.Background())
	if err != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}
	d, err := json.Marshal(v)
	if err != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return

	}
	fmt.Fprint(w, d)
	return
}



func Run(cmd *cobra.Command, args []string) {

	fset := cmd.LocalFlags()
	lis, err := fset.GetString("listen")
	if err != nil {
		panic(err)
	}


	r := mux.NewRouter()

	//镜像
	r.HandleFunc("/api/v1/test", testHandler).Methods("GET")
	r.HandleFunc("/api/v1/image", listImages).Methods("GET")
	r.HandleFunc("/api/v1/image/{id}", deleteImage).Methods("DELETE")
	r.HandleFunc("/api/v1/image", pullImage).Methods("POST")
	r.HandleFunc("/api/v1/image/inspect", inspectImage).Methods("GEt")

	//容器
	r.HandleFunc("/api/v1/container", listContainer).Methods("GET")
	r.HandleFunc("/api/v1/container/{id}/state", stateContainer).Methods("GET")
	r.HandleFunc("/api/v1/container/{id}/stop", stopContainer).Methods("DELETE")
	r.HandleFunc("/api/v1/container/{id}/", deleteContainer).Methods("DELETE")
	r.HandleFunc("/api/v1/container/{id}/inspect", inspectContainer).Methods("GET")
	r.HandleFunc("/api/v1/container", runContainer).Methods("POST")

	r.HandleFunc("/api/v1/docker/version", dockerVersion).Methods("GET")

	s := http.Server{
		Addr: lis,
		Handler: r,
		ReadTimeout: 1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
