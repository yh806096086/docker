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
	"strings"
	"io"
)


type containerCreateConfig struct {
	Config container.HostConfig
	HostConfig container.HostConfig
	NetworkConfig network.NetworkingConfig
}

type Login struct {
	User string `json:"user"`
	Password string `json:"password"`
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

	// /api/v1/image/{id}
	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-1]

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

	// /api/v1/image/{id}/inspect

	path := strings.Split(r.URL.Path, "/")

	id := path[3]
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

	defer r.Body.Close()

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

	// /api/v1/container/{id}/stop

	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-2]
	timeout := 1 *time.Second
	if ok:= cli.ContainerStop(context.Background(), id, &timeout); ok != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}
	fmt.Fprint(w, []byte("stop container"))
	return
}

func statsContainer(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, []byte("state container"))

	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-2]
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

	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-2]

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

	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-1]
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

//docker dick stats

func dockerDiskStats(w http.ResponseWriter, r *http.Request) {

	usage, err := cli.DiskUsage(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	d, err := json.Marshal(usage)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, d)
}

func dockerLogin(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, []byte("Server Error"))
		return
	}

	defer r.Body.Close()

	l := Login{}
	if ok := json.Unmarshal(body, l); ok != nil {
		fmt.Printf("%s\n", err.Error())
	}

	authBody, err := cli.RegistryLogin(context.Background(),
										types.AuthConfig{Username:l.User, Password:l.Password})
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	fmt.Println(authBody)
	io.WriteString(w, "docker login ok")
}

func dockerInfo(w http.ResponseWriter, r *http.Request) {

	info, err := cli.Info(context.Background())
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	d, err := json.Marshal(info)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	fmt.Fprint(w, d)
}

func Run(cmd *cobra.Command, args []string) {

	fset := cmd.LocalFlags()
	lis, err := fset.GetString("listen")
	if err != nil {
		panic(err)
	}


	r := mux.NewRouter()

	r.HandleFunc("/api/v1/test", testHandler).Methods("GET")

	//镜像
	r.HandleFunc("/api/v1/image", listImages).Methods("GET")
	r.HandleFunc("/api/v1/image/{id}", deleteImage).Methods("DELETE")
	r.HandleFunc("/api/v1/image", pullImage).Methods("POST")
	r.HandleFunc("/api/v1/image/{id}/inspect", inspectImage).Methods("GEt")

	//容器
	r.HandleFunc("/api/v1/container", listContainer).Methods("GET")
	r.HandleFunc("/api/v1/container/{id}/stats", statsContainer).Methods("GET")
	r.HandleFunc("/api/v1/container/{id}/stop", stopContainer).Methods("PUT")
	r.HandleFunc("/api/v1/container/{id}/", deleteContainer).Methods("DELETE")
	r.HandleFunc("/api/v1/container/{id}/inspect", inspectContainer).Methods("GET")
	r.HandleFunc("/api/v1/container", runContainer).Methods("POST")

	//docker domain
	r.HandleFunc("/api/v1/docker/version", dockerVersion).Methods("GET")
	r.HandleFunc("/api/v1/dcoker/disk", dockerDiskStats).Methods("GET")
	r.HandleFunc("/api/v1/docker/login", dockerLogin).Methods("POST")
	r.HandleFunc("/api.v1/docker/info", dockerInfo).Methods("GET")

	s := http.Server{
		Addr: lis,
		Handler: r,
		ReadTimeout: 1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
