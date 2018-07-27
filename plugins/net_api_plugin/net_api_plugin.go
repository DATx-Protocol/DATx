package net_api_plugin

import (
    "datx_chain/utils/helper"
    "datx_chain/chainlib/application"
	"glog"
	"net/http"
    "time"
    str "strings"
)

type net_api_plugin struct {
}    

func (self *net_api_plugin) plugin_startup () {
    glog.Info ("starting net_api_plugin")
    // 在 plugins 的池子里找到 net_plugin 的指针
    // net_plugin 通过 boost::core::demangle 获得名称，再根据名称在map里找到指针
    net_mgr := app().Find ("net_plugin")
    // add_api 添加的是 api_description （map键值对）
    app().Find ("http_plugin").add_api ({
        // result    是eos中的一个结构体，可以与string相互转换，但是转换方法特别复杂，为方便起见，这里用string
        // error_results 是eos中的一个结构体，可以与string相互转换，不是特别复杂，为方便起见，这里也用string
        // elog      是eos的日志库，记录了__FILE__, __LINE__, __func__，以及具体的错误信息，为方便起见，这里用glog
        "/v1/net/connect" :
        func (string, var body string, func cb (int,string)) fc::exception {
            try {
                if (body == "")
                    body = "{}"
                result := net_mgr.connect (body)  // connect 方法需要在各插件中实现
                cb (201, result)
            } catch (fc::eof_exception& e) {
                results := "400" + "Bad Request" + e  // error_results结构体，这里直接用string
                cb (400, results)
                glog.Errorf ("Unable to parse arguments  %s", body)
            } catch (fc::exception& e) {
                results := "500" + "Internal Service Error" + e  // error_results结构体，这里直接用string
                cb (500, results))
                glog.Errorf ("Exception encountered while processing  %s %s", "call net.connect", e)
            }
        },
        "/v1/net/disconnect" :
        func (string, var body string, func cb (int,string)) fc::exception {
            try {
                if (body == "")
                    body = "{}"
                result := net_mgr.disconnect (body)  // disconnect 方法需要在各插件中实现
                cb (201, result)
            } catch (fc::eof_exception& e) {
                results := "400" + "Bad Request" + e  // error_results结构体，这里直接用string
                cb (400, results)
                glog.Errorf ("Unable to parse arguments  %s", body)
            } catch (fc::exception& e) {
                results := "500" + "Internal Service Error" + e  // error_results结构体，这里直接用string
                cb (500, results))
                glog.Errorf ("Exception encountered while processing  %s %s", "call net.disconnect", e)
            }
        },
        "/v1/net/status" : 
        func (string, var body string, func cb (int,string)) fc::exception {
            try {
                if (body == "")
                    body = "{}"
                result := net_mgr.status (body)  // status 方法需要在各插件中实现
                cb (201, result)
            } catch (fc::eof_exception& e) {
                results := "400" + "Bad Request" + e  // error_results结构体，这里直接用string
                cb (400, results)
                glog.Errorf ("Unable to parse arguments  %s", body)
            } catch (fc::exception& e) {
                results := "500" + "Internal Service Error" + e  // error_results结构体，这里直接用string
                cb (500, results))
                glog.Errorf ("Exception encountered while processing  %s %s", "call net.status", e)
            }
        },
        "/v1/net/connections" : 
        func (string, var body string, func cb (int,string)) fc::exception {
            try {
                if (body == "")
                    body = "{}"
                net_mgr.connections ()  // connections 方法需要在各插件中实现
                eosio::detail::net_api_plugin_empty result  // 这里的result和前面的不一样
                cb (201, result)
            } catch (fc::eof_exception& e) {
                results := "400" + "Bad Request" + e  // error_results结构体，这里直接用string
                cb (400, results)
                glog.Errorf ("Unable to parse arguments  %s", body)
            } catch (fc::exception& e) {
                results := "500" + "Internal Service Error" + e  // error_results结构体，这里直接用string
                cb (500, results))
                glog.Errorf ("Exception encountered while processing  %s %s", "call net.connections", e)
            }
        },
    })
}

func (self *net_api_plugin) plugin_initialize() {
    log.Println("net_api_plugin initialize")

	err, data := helper.GetFileHelper("net_api_plugin.yaml")
	if err != nil {
		log.Printf("net_api_plugin init chain config error={%v}", err)
		return
	}

	var config controller.CtlConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Printf("net_api_plugin init chain config unmarshal config  error={%v}", err)
		return
	}

    log.Printf("net_api_plugin init config={%v}", config)

    lipstr := config.http-server-address
    lippos := str.Index (lipstr, ':')
    host := lipstr[0:lippos]

    if host != "localhost" && host != "127.0.0.1" {
        glog.Warningf ("\n"
           "*************************************\n"
           "*                                   *\n"
           "*  --  Net API NOT on localhost  -- *\n"
           "*                                   *\n"
           "*   this may be abused if exposed   *\n"
           "*                                   *\n"
           "*************************************\n")
    }
}