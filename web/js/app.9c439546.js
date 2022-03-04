(function(t){function e(e){for(var r,n,i=e[0],s=e[1],c=e[2],m=0,p=[];m<i.length;m++)n=i[m],Object.prototype.hasOwnProperty.call(l,n)&&l[n]&&p.push(l[n][0]),l[n]=0;for(r in s)Object.prototype.hasOwnProperty.call(s,r)&&(t[r]=s[r]);u&&u(e);while(p.length)p.shift()();return o.push.apply(o,c||[]),a()}function a(){for(var t,e=0;e<o.length;e++){for(var a=o[e],r=!0,i=1;i<a.length;i++){var s=a[i];0!==l[s]&&(r=!1)}r&&(o.splice(e--,1),t=n(n.s=a[0]))}return t}var r={},l={app:0},o=[];function n(e){if(r[e])return r[e].exports;var a=r[e]={i:e,l:!1,exports:{}};return t[e].call(a.exports,a,a.exports,n),a.l=!0,a.exports}n.m=t,n.c=r,n.d=function(t,e,a){n.o(t,e)||Object.defineProperty(t,e,{enumerable:!0,get:a})},n.r=function(t){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(t,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(t,"__esModule",{value:!0})},n.t=function(t,e){if(1&e&&(t=n(t)),8&e)return t;if(4&e&&"object"===typeof t&&t&&t.__esModule)return t;var a=Object.create(null);if(n.r(a),Object.defineProperty(a,"default",{enumerable:!0,value:t}),2&e&&"string"!=typeof t)for(var r in t)n.d(a,r,function(e){return t[e]}.bind(null,r));return a},n.n=function(t){var e=t&&t.__esModule?function(){return t["default"]}:function(){return t};return n.d(e,"a",e),e},n.o=function(t,e){return Object.prototype.hasOwnProperty.call(t,e)},n.p="/";var i=window["webpackJsonp"]=window["webpackJsonp"]||[],s=i.push.bind(i);i.push=e,i=i.slice();for(var c=0;c<i.length;c++)e(i[c]);var u=s;o.push([0,"chunk-vendors"]),a()})({0:function(t,e,a){t.exports=a("56d7")},"034f":function(t,e,a){"use strict";a("85ec")},1034:function(t,e,a){"use strict";a("b8f6")},"56d7":function(t,e,a){"use strict";a.r(e);a("e260"),a("e6cf"),a("cca6"),a("a79d");var r=a("2b0e"),l=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{attrs:{id:"app"}},[a("el-container",{attrs:{direction:"vertical"}},[a("el-header",{staticClass:"title"},[t._v("Klog")]),a("el-container",{attrs:{direction:"horizontal"}},[a("el-aside",[a("el-row",[a("el-col",[a("el-menu",{attrs:{"default-active":"/services",router:""}},[a("el-menu-item",{attrs:{index:"/services"}},[a("span",[t._v("Services")])])],1)],1)],1)],1),a("el-main",[a("router-view")],1)],1)],1)],1)},o=[],n={name:"App",data:function(){return{activeIndex:"/"}}},i=n,s=(a("034f"),a("2877")),c=Object(s["a"])(i,l,o,!1,null,null,null),u=c.exports,m=a("5c96"),p=a.n(m),f=a("8c4f"),d=(a("0fae"),function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticStyle:{width:"100%","padding-top":"2%","pading-bottom":"2%"}},[a("el-table",{staticStyle:{width:"96%","margin-left":"2%"},attrs:{data:t.tableData,"span-method":t.objectSpanMethod}},[a("el-table-column",{attrs:{fixed:"",prop:"name",label:"节点名"}}),a("el-table-column",{attrs:{formatter:t.cellFormatter,label:"Cpu/Mem"}}),a("el-table-column",{attrs:{formatter:t.cellFormatter,label:"最大值"}}),a("el-table-column",{attrs:{formatter:t.cellFormatter,label:"最小值"}}),a("el-table-column",{attrs:{formatter:t.cellFormatter,label:"幅度"}}),a("el-table-column",{attrs:{formatter:t.cellFormatter,label:"平均值"}}),a("el-table-column",{attrs:{formatter:t.cellFormatter,label:"最快增幅"}}),a("el-table-column",{attrs:{formatter:t.cellFormatter,label:"磁盘使用情况"}})],1)],1)}),b=[],h=(a("ac1f"),a("1276"),a("bc3a")),v=a.n(h);function g(t){return""==t?"":t.split("T")[1].split(".")[0]}var y={name:"NodeTable",data:function(){return{tableData:[]}},created:function(){var t=this;v.a.get("https://klog.ciiplat.com/nodes").then((function(e){t.tableData=e.data})).catch((function(t){console.log(t)}))},methods:{objectSpanMethod:function(t){t.row,t.column;var e=t.rowIndex,a=t.columnIndex;if(0===a||7===a)return e%2===0?{rowspan:2,colspan:1}:{rowspan:0,colspan:0}},cellFormatter:function(t,e,a,r){if(r%2===0){if("Cpu/Mem"==e.label)return"Cpu";if("最大值"==e.label)return t.cpuSumMax+"% ("+g(t.cpuSumMaxTime)+")";if("最小值"==e.label)return t.cpuSumMin+"% ("+g(t.cpuSumMinTime)+")";if("幅度"==e.label)return t.cpuVolatility+"%";if("平均值"==e.label)return t.cpuSumAvg+"%";if("最快增幅"==e.label)return t.cpuMaxRatio+"%";if("磁盘使用情况"==e.label)return t.diskUsed+"Gi/"+t.diskTotal+"Gi "+t.diskUsedRatio+"%"}if(r%2===1){if("Cpu/Mem"==e.label)return"Mem";if("最大值"==e.label)return t.memMax+"Gi ("+g(t.memMaxTime)+")";if("最小值"==e.label)return t.memMin+"Gi ( "+g(t.memMinTime)+")";if("幅度"==e.label)return t.memVolatility+"Gi";if("平均值"==e.label)return t.memAvg+"Gi";if("最快增幅"==e.label)return t.memMaxRatio+"%";if("磁盘使用情况"==e.label)return t.diskUsed+"Gi/"+t.diskTotal+"Gi "+t.diskUsedRatio+"%"}}}},w=y,x=Object(s["a"])(w,d,b,!1,null,"6c78af68",null),S=(x.exports,function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticStyle:{width:"100%","padding-top":"2%","pading-bottom":"2%"}},[a("el-table",{staticStyle:{width:"96%","margin-left":"2%"},attrs:{data:t.tableData,height:"600"}},[a("el-table-column",{attrs:{fixed:"",prop:"namespace",label:"Namespace"}}),a("el-table-column",{attrs:{fixed:"",prop:"name",label:"Pod"}}),a("el-table-column",{attrs:{prop:"cpuSumMax",label:"Cpu最大值 (m)",sortable:""}}),a("el-table-column",{attrs:{formatter:t.cellFormatter,label:"Cpu最大时刻",sortable:""}}),a("el-table-column",{attrs:{prop:"cpuSumMin",label:"Cpu最小值 (m)",sortable:""}}),a("el-table-column",{attrs:{formatter:t.cellFormatter,label:"Cpu最小时刻",sortable:""}}),a("el-table-column",{attrs:{prop:"cpuVolatility",label:"Cpu幅度 (m)",sortable:""}}),a("el-table-column",{attrs:{prop:"cpuSumAvg",label:"Cpu平均值 (m)",sortable:""}}),a("el-table-column",{attrs:{prop:"cpuMaxRatio",label:"Cpu最快增幅 (m)",sortable:""}}),a("el-table-column",{attrs:{prop:"memMax",label:"Mem最大值 (Mi)",sortable:""}}),a("el-table-column",{attrs:{formatter:t.cellFormatter,label:"Mem最大时刻",sortable:""}}),a("el-table-column",{attrs:{prop:"memMin",label:"Mem最小值 (Mi)",sortable:""}}),a("el-table-column",{attrs:{formatter:t.cellFormatter,label:"Mem最小时刻",sortable:""}}),a("el-table-column",{attrs:{prop:"memVolatility",label:"Mem幅度 (Mi)",sortable:""}}),a("el-table-column",{attrs:{prop:"memAvg",label:"Mem平均值 (Mi)",sortable:""}}),a("el-table-column",{attrs:{prop:"memMaxRatio",label:"Mem最快增幅 (Mi)",sortable:""}})],1)],1)}),_=[];function k(t){return""==t?"":t.split("T")[1].split(".")[0]}var M={name:"PodTable",data:function(){return{tableData:[]}},created:function(){var t=this;v.a.get("http://172.16.179.13:8080/pods").then((function(e){t.tableData=e.data})).catch((function(t){console.log(t)}))},methods:{cellFormatter:function(t,e,a,r){return"Cpu最大时刻"==e.label?k(t.cpuSumMaxTime):"Cpu最小时刻"==e.label?k(t.cpuSumMinTime):"Mem最大时刻"==e.label?k(t.memMaxTime):"Mem最小时刻"==e.label?k(t.memMinTime):void 0}}},J=M,C=Object(s["a"])(J,S,_,!1,null,"65c28a62",null),$=(C.exports,function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"warning-table"},[a("el-table",{staticStyle:{width:"100%"},attrs:{data:t.tableData,"show-header":!1}},[a("el-table-column",{attrs:{type:"expand"},scopedSlots:t._u([{key:"default",fn:function(e){return[a("el-table",{staticClass:"demo-table-expand",staticStyle:{width:"100%"},attrs:{"label-position":"left",inline:"",data:e.row.children}},[a("el-table-column",{attrs:{prop:"type",label:"类型"}}),a("el-table-column",{attrs:{prop:"describe",label:"描述"}}),a("el-table-column",{attrs:{prop:"warningValue",label:"警戒值（%）"},scopedSlots:t._u([{key:"default",fn:function(e){return["Cpu"==e.row.type?[t._v(t._s(e.row.warningValue+"%"))]:[t._v(t._s(e.row.warningValue))]]}}],null,!0)}),a("el-table-column",{attrs:{prop:"actual",label:"实际值"},scopedSlots:t._u([{key:"default",fn:function(e){return["Cpu"==e.row.type?[t._v(" "+t._s(e.row.actual+"%"))]:[t._v(" "+t._s(e.row.actual))]]}}],null,!0)}),a("el-table-column",{attrs:{prop:"time",label:" 时间","min-width":100,formatter:t.timeFormatter}}),a("el-table-column",{attrs:{label:" 嫌疑人","min-width":100},scopedSlots:t._u([{key:"default",fn:function(e){return[a("div",t._l(e.row.suspect,(function(e,r){return a("div",{key:r},[a("p",[a("el-button",{attrs:{type:"warning",round:"",size:"mini"},on:{click:function(a){return t.open(e)}}},[t._v(t._s(e.namespace)+"/"+t._s(e.name))])],1)])})),0)]}}],null,!0)})],1)]}}])}),a("el-table-column",{attrs:{label:"节点名称",prop:"name"}})],1)],1)}),D=[],F=(a("d3b7"),a("4d63"),a("25f0"),a("4d90"),a("5319"),a("53ca"));function N(t,e){if(0===arguments.length||!t)return null;var a,r=e||"{y}-{m}-{d} {h}:{i}:{s}";"object"===Object(F["a"])(t)?a=t:("string"===typeof t&&(t=/^[0-9]+$/.test(t)?parseInt(t):t.replace(new RegExp(/-/gm),"/")),"number"===typeof t&&10===t.toString().length&&(t*=1e3),a=new Date(t));var l={y:a.getFullYear(),m:a.getMonth()+1,d:a.getDate(),h:a.getHours(),i:a.getMinutes(),s:a.getSeconds(),a:a.getDay()},o=r.replace(/{([ymdhisa])+}/g,(function(t,e){var a=l[e];return"a"===e?["日","一","二","三","四","五","六"][a]:a.toString().padStart(2,"0")}));return o}var E={data:function(){return{tableData:[]}},mounted:function(){this.refreshTableData()},methods:{open:function(t){this.$router.push({path:"/podinfo",query:t})},refreshTableData:function(){var t=this;v.a.get("https://klog.ciiplat.com/warnings").then((function(e){t.tableData=e.data})).catch((function(t){console.log(t)}))},timeFormatter:function(t,e){return N(new Date(t.time),"{h}:{i}:{s}")}}},z=E,T=(a("bfc8"),Object(s["a"])(z,$,D,!1,null,null,null)),U=(T.exports,function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticStyle:{width:"100%","padding-top":"2%","pading-bottom":"2%",height:"100%"}},[a("el-card",{staticClass:"box-card",staticStyle:{padding:"2% 1%"}},[a("div",{staticStyle:{padding:"5% 1%"}},[a("ul",{staticStyle:{"padding-inline-start":"0px","list-style":"none","text-align":"left"}},[a("li",{staticStyle:{"border-bottom":"1px solid #EBEEF5"}},[a("p",[a("lable",{staticStyle:{"margin-left":"2%"}},[t._v("Pod名称：")]),a("span",{staticStyle:{float:"right","margin-right":"2%"}},[t._v(t._s(this.suspect.name))])],1)]),a("li",{staticStyle:{"border-bottom":"1px solid #EBEEF5"}},[a("p",[a("lable",{staticStyle:{"margin-left":"2%"}},[t._v("命名空间：")]),a("span",{staticStyle:{float:"right","margin-right":"2%"}},[t._v(t._s(this.suspect.namespace))])],1)]),a("li",{staticStyle:{"border-bottom":"1px solid #EBEEF5"}},[a("p",[a("lable",{staticStyle:{"margin-left":"2%"}},[t._v("实际值：")]),a("span",{staticStyle:{float:"right","margin-right":"2%"}},["PodCpu"==this.suspect.type?[t._v(t._s(this.suspect.actualValue+"%"))]:t._e(),"PodMem"==this.suspect.type?[t._v(t._s(this.suspect.actualValue+"G"))]:t._e()],2)],1)]),a("li",{staticStyle:{"border-bottom":"1px solid #EBEEF5"}},[a("p",[a("lable",{staticStyle:{"margin-left":"2%"}},[t._v("浮动值：")]),a("span",{staticStyle:{float:"right","margin-right":"2%"}},["PodCpu"==this.suspect.type?[t._v(t._s(this.suspect.volatility+"%"))]:t._e(),"PodMem"==this.suspect.type?[t._v(t._s(this.suspect.volatility+"G"))]:t._e()],2)],1)]),a("li",[a("el-button",{staticStyle:{float:"right","margin-top":"8%","margin-right":"5%"},attrs:{type:"primary",icon:"el-icon-arrow-left",circle:""},on:{click:t.back}})],1)])])])],1)}),P=[],j={name:"PodInfo",data:function(){return{suspect:{}}},created:function(){this.suspect=this.$route.query},methods:{back:function(){this.$router.push({path:"/warning"})}}},R=j,O=Object(s["a"])(R,U,P,!1,null,"233fdd84",null),I=(O.exports,function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",[a("el-row",{staticClass:"row"},[a("el-col",{attrs:{span:24}},[a("el-button-group",[a("el-tooltip",{staticClass:"item",attrs:{effect:"light",content:"获取maven打包方式的Java后台服务CI文件",placement:"bottom"}},[a("el-button",{attrs:{type:"primary",icon:"el-icon-document-add",size:"mini",round:""},on:{click:function(e){t.drawerjava=!0}}},[t._v("Java")])],1),a("el-tooltip",{staticClass:"item",attrs:{effect:"light",content:"获取npm打包方式的前台服务CI文件",placement:"bottom"}},[a("el-button",{attrs:{type:"primary",icon:"el-icon-document-add",size:"mini",round:""},on:{click:function(e){t.drawernpm=!0}}},[t._v("Npm")])],1),a("el-tooltip",{staticClass:"item",attrs:{effect:"light",content:"获取dotnet服务CI文件",placement:"bottom"}},[a("el-button",{attrs:{type:"primary",icon:"el-icon-document-add",size:"mini",round:""},on:{click:function(e){t.drawernet=!0}}},[t._v(".net")])],1)],1),a("el-drawer",{attrs:{visible:t.drawerjava,"with-header":!1,size:"50%"},on:{"update:visible":function(e){t.drawerjava=e}}},[a("div",{staticStyle:{padding:"4% 10%"}},[a("el-tabs",[a("el-tab-pane",{staticStyle:{"margin-top":"4%"},attrs:{label:"单模块"}},[a("el-form",{ref:"formJava",attrs:{rules:t.rules,model:t.formJava,size:"small"}},[a("el-form-item",[a("span",[a("font",{staticStyle:{"font-size":"12px",color:"red"}},[t._v("单模块适用的场景，是一个代码仓库中，只存储一个微服务模块。")])],1)]),a("el-form-item",{attrs:{label:"分组",prop:"namespace"}},[a("el-input",{attrs:{placeholder:"请输入分组名，由管理员分配"},model:{value:t.formJava.namespace,callback:function(e){t.$set(t.formJava,"namespace",e)},expression:"formJava.namespace"}})],1),a("el-form-item",{attrs:{label:"名称",prop:"name"}},[a("el-input",{attrs:{placeholder:"请输入微服务英文名称，同一个分组内必须唯一"},model:{value:t.formJava.name,callback:function(e){t.$set(t.formJava,"name",e)},expression:"formJava.name"}})],1),a("el-form-item",{attrs:{label:"描述",prop:"describe"}},[a("el-input",{attrs:{placeholder:"请输入微服务的中文名称"},model:{value:t.formJava.describe,callback:function(e){t.$set(t.formJava,"describe",e)},expression:"formJava.describe"}})],1),a("el-form-item",{attrs:{label:"Git地址",prop:"gitUrl"}},[a("el-input",{attrs:{placeholder:"请输入gitlab仓库地址（拉取代码的url,通常已.git结尾）"},model:{value:t.formJava.gitUrl,callback:function(e){t.$set(t.formJava,"gitUrl",e)},expression:"formJava.gitUrl"}})],1),a("el-form-item",{attrs:{label:"分支名/Tag",prop:"onlyRefs"}},[a("el-input",{attrs:{placeholder:"测试环境请输入测试主分支，生产环境请输入'tags'"},model:{value:t.formJava.onlyRefs,callback:function(e){t.$set(t.formJava,"onlyRefs",e)},expression:"formJava.onlyRefs"}})],1),a("el-form-item",{attrs:{label:"产物",prop:"artifactsPaths"}},[a("el-input",{attrs:{placeholder:"请输入jar包路径，从根目录开始，例如target/demo.jar"},model:{value:t.formJava.artifactsPaths,callback:function(e){t.$set(t.formJava,"artifactsPaths",e)},expression:"formJava.artifactsPaths"}})],1),a("el-form-item",{attrs:{label:"端口",prop:"port"}},[a("el-input",{attrs:{placeholder:"请输入微服务端口，对应spring-boot的server.port配置"},model:{value:t.formJava.port,callback:function(e){t.$set(t.formJava,"port",e)},expression:"formJava.port"}})],1),a("el-form-item",{attrs:{label:"健康检查路径",prop:"health"}},[a("el-input",{attrs:{placeholder:"请输入包含server.servlet.context-path的路径，例如/${context-path}/actuator/health"},model:{value:t.formJava.health,callback:function(e){t.$set(t.formJava,"health",e)},expression:"formJava.health"}})],1),a("el-form-item",[a("span",[a("font",{staticStyle:{"font-size":"12px",color:"red"}},[t._v("获取配置文件后，解压到项目的根目录，一起提交。")])],1)]),a("el-form-item",[a("el-button",{attrs:{type:"primary"},on:{click:function(e){return t.submitFormJava("formJava")}}},[t._v("获取配置文件")]),a("el-button",{on:{click:function(e){return t.resetFormJava("formJava")}}},[t._v("重置")])],1)],1)],1),a("el-tab-pane",{staticStyle:{"margin-top":"4%",height:"750px","overflow-y":"scroll"},attrs:{label:"多模块"}},[a("el-form",{ref:"formJavaMul",attrs:{rules:t.rules,model:t.formJavaMul,size:"small"}},[a("el-form-item",[a("span",[a("font",{staticStyle:{"font-size":"12px",color:"red"}},[t._v("多模块适用的场景，是一个代码仓库中，存储了多个微服务模块。")])],1)]),a("el-form-item",{attrs:{label:"分组",prop:"namespace"}},[a("el-input",{attrs:{placeholder:"请输入分组名，由管理员分配"},model:{value:t.formJavaMul.namespace,callback:function(e){t.$set(t.formJavaMul,"namespace",e)},expression:"formJavaMul.namespace"}})],1),a("el-form-item",{attrs:{label:"Git地址",prop:"gitUrl"}},[a("el-input",{attrs:{placeholder:"请输入gitlab仓库地址（拉取代码的url,通常已.git结尾）"},model:{value:t.formJavaMul.gitUrl,callback:function(e){t.$set(t.formJavaMul,"gitUrl",e)},expression:"formJavaMul.gitUrl"}})],1),a("el-form-item",{attrs:{label:"分支名/Tag",prop:"onlyRefs"}},[a("el-input",{attrs:{placeholder:"测试环境请输入测试主分支，生产环境请输入'tags'"},model:{value:t.formJavaMul.onlyRefs,callback:function(e){t.$set(t.formJavaMul,"onlyRefs",e)},expression:"formJavaMul.onlyRefs"}})],1),a("el-form-item",[a("span",[a("font",{staticStyle:{"font-size":"12px",color:"red"}},[t._v("获取配置文件后，解压到项目的根目录，一起提交。")])],1)]),a("el-form-item",[a("el-button",{attrs:{type:"primary"},on:{click:function(e){return t.submitFormJavaMul("formJavaMul")}}},[t._v("获取配置文件")]),a("el-button",{on:{click:function(e){return t.addService()}}},[t._v("新增域名")]),a("el-button",{on:{click:function(e){return t.resetFormJavaMul("formJavaMul")}}},[t._v("重置")])],1),t._l(t.formJavaMul.modules,(function(e,r){return a("div",{key:e.key},[a("div",{staticStyle:{margin:"2% 2%",padding:"4% 4%","box-shadow":"2px 2px 5px #888888"}},[a("div",[a("h4",{staticStyle:{float:"left"}},[t._v("子模块"+t._s(r+1))]),a("el-button",{staticStyle:{float:"right"},attrs:{type:"warning",circle:"",icon:"el-icon-minus",size:"mini"},on:{click:function(a){return a.preventDefault(),t.removeService(e)}}})],1),a("table",{staticStyle:{width:"100%"}},[a("tr",[a("td",{staticClass:"td1"},[a("label",{attrs:{for:"name"}},[t._v("名称:")])]),a("td",{staticClass:"td2"},[a("el-form-item",{staticStyle:{margin:"0px 0px"},attrs:{rules:t.mulRules.name,prop:"modules."+r+".name"}},[a("el-input",{attrs:{label:"名称",prop:"name",placeholder:"请输入微服务英文名称，同一个分组内必须唯一"},model:{value:e.name,callback:function(a){t.$set(e,"name",a)},expression:"module.name"}})],1)],1)]),a("tr",[a("td",[a("label",{attrs:{for:"describe"}},[t._v("描述:")])]),a("td",[a("el-form-item",{staticStyle:{margin:"0px 0px"},attrs:{rules:t.mulRules.describe,prop:"modules."+r+".describe"}},[a("el-input",{attrs:{label:"描述",prop:"describe",placeholder:"请输入微服务的中文名称"},model:{value:e.describe,callback:function(a){t.$set(e,"describe",a)},expression:"module.describe"}})],1)],1)]),a("tr",[a("td",[a("label",{attrs:{for:"artifactsPaths"}},[t._v("产物:")])]),a("td",[a("el-form-item",{staticStyle:{margin:"0px 0px"},attrs:{rules:t.mulRules.artifactsPaths,prop:"modules."+r+".artifactsPaths"}},[a("el-input",{attrs:{label:"产物",prop:"artifactsPaths",placeholder:"请输入jar包路径，从根目录开始，例如target/demo.jar"},model:{value:e.artifactsPaths,callback:function(a){t.$set(e,"artifactsPaths",a)},expression:"module.artifactsPaths"}})],1)],1)]),a("tr",[a("td",[a("label",{attrs:{for:"port"}},[t._v("端口:")])]),a("td",[a("el-form-item",{staticStyle:{margin:"0px 0px"},attrs:{rules:t.mulRules.port,prop:"modules."+r+".port"}},[a("el-input",{attrs:{label:"端口",prop:"port",placeholder:"请输入微服务端口，对应spring-boot的server.port配置"},model:{value:e.port,callback:function(a){t.$set(e,"port",a)},expression:"module.port"}})],1)],1)]),a("tr",[a("td",[a("label",{attrs:{for:"health"}},[t._v("健康检查路径:")])]),a("td",[a("el-form-item",{key:e.key,staticStyle:{margin:"0px 0px"},attrs:{rules:t.mulRules.health,prop:"modules."+r+".health"}},[a("el-input",{attrs:{label:"健康检查路径",prop:"health",placeholder:"请输入包含server.servlet.context-path的路径，例如/${context-path}/actuator/health"},model:{value:e.health,callback:function(a){t.$set(e,"health",a)},expression:"module.health"}})],1)],1)])])])])}))],2)],1)],1)],1)]),a("el-drawer",{attrs:{visible:t.drawernpm,"with-header":!1,size:"50%"},on:{"update:visible":function(e){t.drawernpm=e}}},[a("div",{staticStyle:{padding:"4% 10%"}},[a("el-form",{ref:"formNpm",attrs:{rules:t.rules,model:t.formNpm,size:"small"}},[a("el-form-item",{attrs:{label:"分组",prop:"namespace"}},[a("el-input",{attrs:{placeholder:"请输入分组名，由管理员分配"},model:{value:t.formNpm.namespace,callback:function(e){t.$set(t.formNpm,"namespace",e)},expression:"formNpm.namespace"}})],1),a("el-form-item",{attrs:{label:"名称",prop:"name"}},[a("el-input",{attrs:{placeholder:"请输入微服务英文名称，同一个分组内必须唯一"},model:{value:t.formNpm.name,callback:function(e){t.$set(t.formNpm,"name",e)},expression:"formNpm.name"}})],1),a("el-form-item",{attrs:{label:"描述",prop:"describe"}},[a("el-input",{attrs:{placeholder:"请输入微服务的中文名称"},model:{value:t.formNpm.describe,callback:function(e){t.$set(t.formNpm,"describe",e)},expression:"formNpm.describe"}})],1),a("el-form-item",{attrs:{label:"Git地址",prop:"gitUrl"}},[a("el-input",{attrs:{placeholder:"请输入gitlab仓库地址（拉取代码的url,通常已.git结尾）"},model:{value:t.formNpm.gitUrl,callback:function(e){t.$set(t.formNpm,"gitUrl",e)},expression:"formNpm.gitUrl"}})],1),a("el-form-item",{attrs:{label:"分支名/Tag",prop:"onlyRefs"}},[a("el-input",{attrs:{placeholder:"测试环境请输入测试主分支，生产环境请输入'tags'"},model:{value:t.formNpm.onlyRefs,callback:function(e){t.$set(t.formNpm,"onlyRefs",e)},expression:"formNpm.onlyRefs"}})],1),a("el-form-item",{attrs:{label:"健康检查路径",prop:"health"}},[a("el-input",{attrs:{placeholder:"请输入一个能访问的静态资源，例如/index.html"},model:{value:t.formNpm.health,callback:function(e){t.$set(t.formNpm,"health",e)},expression:"formNpm.health"}})],1),a("el-form-item",[a("span",[a("font",{staticStyle:{"font-size":"12px",color:"red"}},[t._v("获取配置文件后，解压到项目的根目录，一起提交。")])],1)]),a("el-form-item",[a("el-button",{attrs:{type:"primary"},on:{click:function(e){return t.submitFormNpm("formNpm")}}},[t._v("获取配置文件")]),a("el-button",{on:{click:function(e){return t.resetFormNpm("formNpm")}}},[t._v("重置")])],1)],1)],1)])],1)],1),a("el-row",{staticClass:"row"},[a("el-col",{attrs:{span:24}},[a("el-table",{attrs:{data:t.tableData}},[a("el-table-column",{attrs:{fixed:"",prop:"namespace",label:"Namespace",width:"100"}}),a("el-table-column",{attrs:{fixed:"",prop:"name",label:"名称",width:"100"}}),a("el-table-column",{attrs:{prop:"description",label:"说明",width:"150"}}),a("el-table-column",{attrs:{prop:"kind",label:"负载类型",width:"100"}}),a("el-table-column",{attrs:{prop:"podStatus",formatter:t.podStatusFormatter,label:"状态",width:"100"}}),a("el-table-column",{attrs:{prop:"podNames",formatter:t.podNamesFormatter,label:"下辖Pod",width:"200"}}),a("el-table-column",{attrs:{prop:"startTime",formatter:t.startTimeFormatter,label:"启动时间",width:"200"}}),a("el-table-column",{attrs:{prop:"duration",formatter:t.durationFormatter,label:"存活时间（秒）",width:"200"}}),a("el-table-column",{attrs:{prop:"restartTimes",formatter:t.restartTimesFormatter,label:"重启次数",width:"100"}}),a("el-table-column",{attrs:{prop:"innerUrl",label:"内部调用地址",width:"300"}}),a("el-table-column",{attrs:{prop:"gatewayUrl",label:"外部调用地址",width:"400"}}),a("el-table-column",{attrs:{prop:"gitUrl",label:"仓库地址",width:"400"}}),a("el-table-column",{attrs:{prop:"replicas",label:"要求副本数",width:"100"}}),a("el-table-column",{attrs:{prop:"ready",label:"正常运行副本数",width:"150"}}),a("el-table-column",{attrs:{fixed:"right",label:"操作",width:"160"},scopedSlots:t._u([{key:"default",fn:function(e){return[a("el-tooltip",{staticClass:"item",attrs:{effect:"light",content:"重启",placement:"bottom"}},[a("el-button",{attrs:{type:"primary",icon:"el-icon-refresh",circle:""},on:{click:function(a){return t.restartClick(e.row)}}})],1),a("el-tooltip",{staticClass:"item",attrs:{effect:"light",content:"回滚",placement:"bottom"}},[a("el-button",{attrs:{type:"primary",icon:"el-icon-d-arrow-left",circle:""},on:{click:function(a){return t.rollbackClick(e.row)}}})],1),a("el-tooltip",{staticClass:"item",attrs:{effect:"light",content:"详情",placement:"bottom"}},[a("el-button",{attrs:{type:"primary",icon:"el-icon-info",circle:""},on:{click:function(a){return t.infoClick(e.row)}}})],1)]}}])})],1)],1)],1)],1)}),G=[],V=(a("4160"),a("c975"),a("a434"),a("b0c0"),a("159b"),{name:"PodTable",data:function(){var t=function(t,e,a){return e?e.length>15?a(new Error("长度不能超过15")):!/[a-z0-9]([-a-z0-9]*[a-z0-9])?(.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/g.test(e)||/[A-Z]/g.test(e)?a(new Error("必须由小写字母,数字,字符“-”或“.”组成")):void a():a(new Error("不能为空"))},e=function(t,e,a){return e?e.length>15?a(new Error("长度不能超过15")):!/[a-z0-9]([-a-z0-9]*[a-z0-9])?(.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/g.test(e)||/[A-Z]/g.test(e)?a(new Error("必须由小写字母,数字,字符“-”或“.”组成")):void a():a(new Error("不能为空"))},a=function(t,e,a){if(!e)return a(new Error("不能为空"));a()},r=function(t,e,a){if(!e)return a(new Error("不能为空"));a()},l=function(t,e,a){if(!e)return a(new Error("不能为空"));a()},o=function(t,e,a){if(!e)return a(new Error("不能为空"));a()},n=function(t,e,a){if(!e)return a(new Error("不能为空"));a()},i=function(t,e,a){if(!e)return a(new Error("不能为空"));a()};return{tableData:[],formJava:{namespace:"",name:"",describe:"",gitUrl:"",onlyRefs:"",artifactsPaths:"",port:"",health:""},formJavaMul:{namespace:"",onlyRefs:"",gitUrl:"",modules:[{name:"",describe:"",artifactsPaths:"",port:"",health:""}]},formNpm:{namespace:"",name:"",describe:"",gitUrl:"",onlyRefs:"",health:""},drawerjava:!1,drawernpm:!1,drawernet:!1,rules:{namespace:[{validator:t,trigger:"blur"}],name:[{validator:e,trigger:"blur"}],describe:[{validator:a,trigger:"blur"}],gitUrl:[{validator:r,trigger:"blur"}],onlyRefs:[{validator:l,trigger:"blur"}],artifactsPaths:[{validator:o,trigger:"blur"}],port:[{validator:n,trigger:"blur"}],health:[{validator:i,trigger:"blur"}]},mulRules:{name:[{validator:e,trigger:"blur"}],describe:[{validator:a,trigger:"blur"}],artifactsPaths:[{validator:o,trigger:"blur"}],port:[{validator:n,trigger:"blur"}],health:[{validator:i,trigger:"blur"}]}}},created:function(){this.getServiceList(),console.log(this.global.apiUrl)},methods:{getServiceList:function(){var t=this;v.a.get(this.global.apiUrl+"services",{params:{action:"list"}}).then((function(e){t.tableData=e.data})).catch((function(t){console.log(t)}))},podNamesFormatter:function(t,e,a,r){var l="";return a.forEach((function(t){l+=t+"\n"})),l},podStatusFormatter:function(t,e,a,r){var l="";return a.forEach((function(t){l+=t+"\n"})),l},startTimeFormatter:function(t,e,a,r){var l="";return a.forEach((function(t){var e=t.split(" ");l+=e[0]+" "+e[1]+"\n"})),l},durationFormatter:function(t,e,a,r){for(var l="",o=0;o<a.length;o++){var n=a[o],i=this.$options.methods.formatterSecondTime(n);l+=i+"\n"}return l},restartTimesFormatter:function(t,e,a,r){var l="";return a.forEach((function(t){l+=t+"\n"})),l},formatterSecondTime:function(t){var e=parseInt(t),a=0,r=0,l=0;e>60&&(a=parseInt(e/60),e=parseInt(e%60),a>60&&(r=parseInt(a/60),a=parseInt(a%60),r>24&&(l=parseInt(a/24),r=parseInt(a%24))));var o=parseInt(e)+"秒";return o=(a>0?parseInt(a):"0")+"分"+o,o=(r>0?parseInt(r):"0")+"小时"+o,o=(l>0?parseInt(l):"0")+"天"+o,o},restartClick:function(t){var e=this;v.a.post(this.global.apiUrl+"services",JSON.stringify({action:"restart",name:t.name,namespace:t.namespace,kind:t.kind})).then((function(t){"success"==t.data?e.$message({showClose:!0,message:"重启成功",type:"success"}):e.$message({showClose:!0,message:t.data,type:"error"})})).catch((function(t){console.log(t)}))},rollbackClick:function(t){var e=this;v.a.post(this.global.apiUrl+"services",JSON.stringify({action:"rollback",name:t.name,namespace:t.namespace,kind:t.kind})).then((function(t){"success"==t.data?e.$message({showClose:!0,message:"重启成功",type:"success"}):e.$message({showClose:!0,message:t.data,type:"error"})})).catch((function(t){console.log(t)}))},infoClick:function(t){this.$router.push({path:"/servicesinfo",query:t})},submitFormJava:function(t){var e=this;this.$refs[t].validate((function(t){if(!t)return console.log("error submit!!"),!1;v.a.put(e.global.apiUrl+"services?action=java",JSON.stringify(e.formJava)).then((function(t){"success"==t.data?(e.$message({showClose:!0,message:"创建成功",type:"success"}),window.open(e.global.apiUrl+"ci.zip")):e.$message({showClose:!0,message:t.data,type:"error"})})).catch((function(t){console.log(t)}))}))},resetFormJava:function(t){this.$refs[t].resetFields()},submitFormJavaMul:function(t){var e=this;console.log(this.formJavaMul),this.$refs[t].validate((function(t){if(!t)return console.log("error submit!!"),!1;v.a.put(e.global.apiUrl+"services?action=javaMul",JSON.stringify(e.formJavaMul)).then((function(t){"success"==t.data?(e.$message({showClose:!0,message:"创建成功",type:"success"}),window.open(e.global.apiUrl+"ci.zip")):e.$message({showClose:!0,message:t.data,type:"error"})})).catch((function(t){console.log(t)}))}))},resetFormJavaMul:function(t){this.$refs[t].resetFields()},removeService:function(t){var e=this.formJavaMul.modules.indexOf(t);-1!==e&&this.formJavaMul.modules.splice(e,1)},addService:function(){this.formJavaMul.modules.push({value:"",key:Date.now()})},submitFormNpm:function(t){var e=this;this.$refs[t].validate((function(t){if(!t)return console.log("error submit!!"),!1;v.a.put(e.global.apiUrl+"services?action=npm",JSON.stringify(e.formNpm)).then((function(t){"success"==t.data?(e.$message({showClose:!0,message:"创建成功",type:"success"}),window.open(e.global.apiUrl+"ci.zip")):e.$message({showClose:!0,message:t.data,type:"error"})})).catch((function(t){console.log(t)}))}))},resetFormNpm:function(t){this.$refs[t].resetFields()}}}),A=V,q=(a("63ee"),Object(s["a"])(A,I,G,!1,null,null,null)),B=q.exports,L=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"demo-image__lazy"},[a("el-input",{attrs:{placeholder:"请输入日期(例如:2000-01-01)"},on:{change:t.inputChange},model:{value:t.input,callback:function(e){t.input=e},expression:"input"}}),t._l(t.urls,(function(t){return a("div",{key:t,staticStyle:{overflow:"scroll"},attrs:{bac:""}},[a("el-image",{staticClass:"result-image",attrs:{src:t,lazy:""}})],1)}))],2)},H=[],K={name:"Echarts",created:function(){this.input=this.getCurrentTime(),this.getImagesName()},methods:{getImagesName:function(){var t=this;v.a.get("https://klog.ciiplat.com/chart?date="+this.input).then((function(e){t.urls=[];for(var a=0;a<e.data.length;a++)t.urls.push("https://klog.ciiplat.com/ai/result/"+t.input+"/"+e.data[a])})).catch((function(t){console.log(t)}))},getCurrentTime:function(){var t=(new Date).getFullYear(),e=(new Date).getMonth()+1<10?"0"+((new Date).getMonth()+1):(new Date).getMonth()+1,a=(new Date).getDate()<10?"0"+(new Date).getDate():(new Date).getDate();(new Date).getHours(),(new Date).getMinutes(),(new Date).getMinutes(),(new Date).getSeconds(),(new Date).getSeconds();return t+"-"+e+"-"+a},inputChange:function(t){this.input=t,this.getImagesName()}},data:function(){return{urls:[],input:""}}},Y=K,Z=(a("e3e9"),Object(s["a"])(Y,L,H,!1,null,null,null)),Q=(Z.exports,function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{},[a("el-row",{attrs:{gutter:20}},[a("el-col",{attrs:{span:24}},[a("el-button-group",{staticStyle:{float:"right","margin-top":"1%"}},[a("el-button",{attrs:{type:"primary",icon:"el-icon-back"},on:{click:t.back}}),a("el-button",{attrs:{type:"primary",icon:"el-icon-view"},on:{click:function(e){t.drawer=!0}}}),a("el-drawer",{attrs:{visible:t.drawer,"with-header":!1,size:"50%"},on:{"update:visible":function(e){t.drawer=e}}},[a("el-collapse",{staticStyle:{margin:"2% 2%"},model:{value:t.activeNames,callback:function(e){t.activeNames=e},expression:"activeNames"}},t._l(t.logData,(function(e){return a("el-collapse-item",{key:e,attrs:{title:e.name}},[a("div",{staticStyle:{overflow:"scroll",height:"900px","white-space":"pre-line"}},[t._v(t._s(e.content))])])})),1)],1)],1)],1)],1),a("el-row",{attrs:{gutter:20}},[a("el-col",{attrs:{span:24}},[a("el-row",{attrs:{gutter:20}},[a("el-col",{attrs:{span:24}},[a("div",{staticClass:"grid-content-desc"},[a("p",[t._v(" Deployment Describe ")]),a("el-table",{staticStyle:{width:"100%"},attrs:{data:t.deploymentData,size:"mini"}},[a("el-table-column",{attrs:{prop:"time",formatter:t.startTimeFormatter,label:"时间",width:"150"}}),a("el-table-column",{attrs:{prop:"kind",label:"资源类型",width:"200"}}),a("el-table-column",{attrs:{prop:"name",label:"资源名称",width:"200"}}),a("el-table-column",{attrs:{prop:"type",label:"类型",width:"200"}}),a("el-table-column",{attrs:{prop:"reason",label:"原因",width:"200"}}),a("el-table-column",{attrs:{prop:"message",label:"信息",width:"580"}})],1)],1)])],1),a("el-row",{attrs:{gutter:20}},[a("p",[t._v(" Replicaset Describe ")]),a("el-col",{attrs:{span:24}},[a("div",{staticClass:"grid-content-desc"},[a("el-table",{staticStyle:{width:"100%"},attrs:{data:t.replicasetData,size:"mini"}},[a("el-table-column",{attrs:{prop:"time",formatter:t.startTimeFormatter,label:"时间",width:"150"}}),a("el-table-column",{attrs:{prop:"kind",label:"资源类型",width:"200"}}),a("el-table-column",{attrs:{prop:"name",label:"资源名称",width:"200"}}),a("el-table-column",{attrs:{prop:"type",label:"类型",width:"200"}}),a("el-table-column",{attrs:{prop:"reason",label:"原因",width:"200"}}),a("el-table-column",{attrs:{prop:"message",label:"信息",width:"580"}})],1)],1)])],1),a("el-row",{attrs:{gutter:20}},[a("p",[t._v(" Pod Describe ")]),a("el-col",{attrs:{span:24}},[a("div",{staticClass:"grid-content-desc"},[a("el-table",{staticStyle:{width:"100%"},attrs:{data:t.podData,size:"mini"}},[a("el-table-column",{attrs:{prop:"time",formatter:t.startTimeFormatter,label:"时间",width:"150"}}),a("el-table-column",{attrs:{prop:"kind",label:"资源类型",width:"200"}}),a("el-table-column",{attrs:{prop:"name",label:"资源名称",width:"200"}}),a("el-table-column",{attrs:{prop:"type",label:"类型",width:"200"}}),a("el-table-column",{attrs:{prop:"reason",label:"原因",width:"200"}}),a("el-table-column",{attrs:{prop:"message",label:"信息",width:"580"}})],1)],1)])],1)],1)],1)],1)}),W=[],X={name:"ServicesInfo",data:function(){return{suspect:{},deploymentData:[],replicasetData:[],podData:[],drawer:!1,logData:[{name:"name1",content:"content1"},{name:"name2",content:"content2"}]}},created:function(){this.suspect=this.$route.query,this.getDesc(),this.getLog()},methods:{back:function(){this.$router.push({path:"/services"})},getDesc:function(){var t=this;v.a.get(this.global.apiUrl+"services",{params:{action:"desc",namespace:this.suspect.namespace,name:this.suspect.name}}).then((function(e){t.deploymentData=e.data.deployment,t.replicasetData=e.data.replicaset,t.podData=e.data.pod})).catch((function(t){console.log(t)}))},getLog:function(){var t=this;v.a.get(this.global.apiUrl+"services",{params:{action:"log",namespace:this.suspect.namespace,name:this.suspect.name}}).then((function(e){t.logData=e.data})).catch((function(t){console.log(t)}))},startTimeFormatter:function(t,e,a,r){var l="",o=a.split(" ");return l+=o[0]+" "+o[1]+"\n",l}}},tt=X,et=(a("1034"),Object(s["a"])(tt,Q,W,!1,null,"3e36d5d2",null)),at=et.exports,rt=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticStyle:{"text-align":"center"}},[a("P",[a("span",{staticStyle:{"font-family":"微软雅黑","font-size":"40px","font-weight":"bold",color:"#999"}},[t._v("微服务平台")])])],1)},lt=[],ot={name:"Index",data:function(){return{}},created:function(){},methods:{}},nt=ot,it=Object(s["a"])(nt,rt,lt,!1,null,"e7f2a896",null),st=it.exports,ct=[{path:"/",component:st,meta:{title:"Index"}},{path:"/services",component:B,meta:{title:"Services"}},{path:"/servicesinfo",component:at,meta:{title:"ServicesInfo"}}],ut=ct,mt=a("313e"),pt=a.n(mt),ft=window.location.href.split("#")[0],dt={apiUrl:ft};r["default"].config.productionTip=!1,r["default"].use(p.a),r["default"].use(f["a"]),r["default"].prototype.$echarts=pt.a,r["default"].prototype.global=dt;var bt=new f["a"]({routes:ut});document.title="Klog",new r["default"]({router:bt,render:function(t){return t(u)}}).$mount("#app")},"63ee":function(t,e,a){"use strict";a("8b78")},"85ec":function(t,e,a){},"8b78":function(t,e,a){},a0dc:function(t,e,a){},b8f6:function(t,e,a){},bfc8:function(t,e,a){"use strict";a("fcf7")},e3e9:function(t,e,a){"use strict";a("a0dc")},fcf7:function(t,e,a){}});
//# sourceMappingURL=app.9c439546.js.map