package api

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	gogotypes "github.com/gogo/protobuf/types"
	"golang.org/x/net/context"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog-controller/controller"
	"k8s.io/klog-controller/log"
	"k8s.io/klog/v2"
	"strings"
)

type EnvoyFilterController struct {
	beego.Controller
	IstioClient *versionedclient.Clientset
	Ctl         *controller.Controller
}

func (this *EnvoyFilterController) Get() {

	if efs, err := this.IstioClient.NetworkingV1alpha3().EnvoyFilters("istio-system").List(context.Background(), v1.ListOptions{}); err != nil {
		klog.Error(err.Error())
		this.Ctx.WriteString(err.Error())
	} else {
		whiteList := make([]log.WhiteList, 0)
		for _, name := range this.Ctl.Log.Spec.Context.EnvoyFilters {
			for _, ef := range efs.Items {
				if ef.Name == name {
					b, _ := ef.Spec.ConfigPatches[0].Patch.Value.Fields["config"].Marshal()
					ss := strings.Split(string(b), "-- function for user auth")[1]
					sr := strings.Split(ss, "-- whitelist")
					s := sr[1]
					se := strings.Split(s, "\n")
					for _, row := range se {
						if !strings.Contains(row, "path ==") && !strings.Contains(row, "string.sub(path") {
							continue
						} else {
							path := strings.Split(row, "\"")[1]
							if strings.Contains(row, "path ==") {
								whiteList = append(whiteList, log.WhiteList{Path: path, Type: "exact"})
							} else if strings.Contains(row, "string.sub(path") {
								whiteList = append(whiteList, log.WhiteList{Path: path, Type: "prefix"})
							}
						}
					}
				}
			}
		}
		if len(whiteList) != 0 {
			b, err := json.Marshal(whiteList)
			if err != nil {
				this.Ctx.WriteString(err.Error())
			} else {
				this.Ctx.WriteString(string(b))
			}
		} else {
			this.Ctx.WriteString("fail")
		}
	}
}

func (this *EnvoyFilterController) Post() {
	input := make([]log.WhiteList, 0)
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &input); err != nil {
		this.Data["json"] = err.Error()
	}
	if efs, err := this.IstioClient.NetworkingV1alpha3().EnvoyFilters("istio-system").List(context.Background(), v1.ListOptions{}); err != nil {
		klog.Error(err.Error())
		this.Ctx.WriteString(err.Error())
	} else {
		for _, name := range this.Ctl.Log.Spec.Context.EnvoyFilters {
			for _, ef := range efs.Items {
				if ef.Name == name {
					value := ef.Spec.ConfigPatches[0].Patch.Value.Fields["config"]
					b, _ := ef.Spec.ConfigPatches[0].Patch.Value.Fields["config"].Marshal()
					ss := strings.Split(string(b), "-- function for user auth")[1]
					sr := strings.Split(ss, "-- whitelist")
					var build string = ""
					for index, wl := range input {
						if wl.Type == "exact" {
							build = build + "  path == \"" + wl.Path + "\" "
							if index != len(input)-1 {
								build = build + "or\n"
							} else {
								build = build + "\n"
							}
						} else if wl.Type == "prefix" {
							build = build + "  string.sub(path,1," + fmt.Sprintf("%d", len(wl.Path)) + ") == \"" + wl.Path + "\" or\n"
						}
					}
					build = "-- function for user auth\n" + sr[0] + "-- whitelist\n  local inwhitelist =\n" + build + "  -- whitelist\n" + sr[2]
					v := value.Kind.(*gogotypes.Value_StructValue)
					v.StructValue.Fields["inlineCode"] = &gogotypes.Value{Kind: &gogotypes.Value_StringValue{
						StringValue: build,
					}}
					this.IstioClient.NetworkingV1alpha3().EnvoyFilters("istio-system").Update(context.Background(), &ef, v1.UpdateOptions{})
				}
			}
		}
		this.Ctx.WriteString("success")
	}
}
