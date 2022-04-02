package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v3"
	"log"
	"time"
)

type RoomTemplate struct {
	Pusher_common_client struct {
		EtcdUrls []string `yaml:"etcd_urls"`
	}
	Generatecenter_common_client struct {
		EtcdUrls []string `yaml:"etcd_urls"`
	}
	Common_client struct {
		EtcdUrls []string `yaml:"etcd_urls"`
	}
	Room struct {
		ListenPort  string   `yaml:"listen_port"`
		RpcBasePath string   `yaml:"rpc_base_path"`
		GameType    string   `yaml:"game_type"`
		IsTestEnv   string   `yaml:"is_test_env"`
		EtcdUrls    []string `yaml:"etcd_urls"`
	}
	Log struct {
		Level string `yaml:"level"`
	}
	Pprof struct {
		Enabled bool   `yaml:"enabled"`
		Address string `yaml:"address"`
	}
	Report struct {
		Address      string `yaml:"address"`
		Level        string `yaml:"level"`
		KafkaTopic   string `yaml:"kafka_topic"`
		KafkaBrokers string `yaml:"kafka_brokers"`
	}
}

func main() {

	var RoomConfitTemplate RoomTemplate

	sc := []constant.ServerConfig{{
		IpAddr: "192.168.101.171",
		Port:   8848,
	}}

	cc := constant.ClientConfig{
		NamespaceId:         "7bde5507-ffea-46bf-a81a-a8c9e55153cb", //命名空间Id
		TimeoutMs:           5000,                                   // 请求Nacos服务端的超时时间，默认是10000ms
		NotLoadCacheAtStart: true,                                   // 在启动的时候不读取缓存在CacheDir的service信息
		LogDir:              "log",                                  // 日志存储路径
		CacheDir:            "cache",                                // 缓存service信息的目录，默认是当前运行目录
		LogLevel:            "debug",                                // 日志默认级别，值必须是：debug,info,warn,error，默认值是info
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	content, err2 := configClient.GetConfig(vo.ConfigParam{
		DataId: "config-dev.yaml",
		Group:  "DEFAULT_GROUP",
	})
	if err2 != nil {
		panic(err2)
	}

	fmt.Println(content) //字符串 - yaml

	err = yaml.Unmarshal([]byte(content), &RoomConfitTemplate)

	if err != nil {
		log.Fatalf("解析失败,%v", err)
	}

	fmt.Println("roomConfitTemplate:", RoomConfitTemplate)

	//url1 := roomConfitTemplate.PusherCommonClient.EtcdUrls
	//url2 := roomConfitTemplate.GeneratecenterCommonClient.EtcdUrls
	//listen_port := roomConfitTemplate.Room.ListenPort
	//rpc_base_path := roomConfitTemplate.Room.RpcBasePath
	//game_type := roomConfitTemplate.Room.GameType
	//is_test_env := roomConfitTemplate.Room.IsTestEnv
	//etcd_urls := roomConfitTemplate.Room.EtcdUrls
	//
	//fmt.Printf("url1:%s,url2:%s,listen_port:%s,rpc_base_path:%s,game_type:%s,is_test_env:%s,etcd_urls:%s", url1, url2, listen_port, rpc_base_path, game_type, is_test_env, etcd_urls)

	configClient.ListenConfig(vo.ConfigParam{
		DataId: "config-dev.yaml",
		Group:  "DEFAULT_GROUP",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置文件发生了变化...")
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
			yaml.Unmarshal([]byte(data), &RoomConfitTemplate)
		},
	})

	t1 := time.NewTimer(time.Second * 5)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 5)
			fmt.Println("roomConfitTemplate:", RoomConfitTemplate)
		}
	}
}
