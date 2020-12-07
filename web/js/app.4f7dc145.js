(function(e){function t(t){for(var l,o,i=t[0],u=t[1],c=t[2],s=0,p=[];s<i.length;s++)o=i[s],Object.prototype.hasOwnProperty.call(a,o)&&a[o]&&p.push(a[o][0]),a[o]=0;for(l in u)Object.prototype.hasOwnProperty.call(u,l)&&(e[l]=u[l]);f&&f(t);while(p.length)p.shift()();return n.push.apply(n,c||[]),r()}function r(){for(var e,t=0;t<n.length;t++){for(var r=n[t],l=!0,i=1;i<r.length;i++){var u=r[i];0!==a[u]&&(l=!1)}l&&(n.splice(t--,1),e=o(o.s=r[0]))}return e}var l={},a={app:0},n=[];function o(t){if(l[t])return l[t].exports;var r=l[t]={i:t,l:!1,exports:{}};return e[t].call(r.exports,r,r.exports,o),r.l=!0,r.exports}o.m=e,o.c=l,o.d=function(e,t,r){o.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:r})},o.r=function(e){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},o.t=function(e,t){if(1&t&&(e=o(e)),8&t)return e;if(4&t&&"object"===typeof e&&e&&e.__esModule)return e;var r=Object.create(null);if(o.r(r),Object.defineProperty(r,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var l in e)o.d(r,l,function(t){return e[t]}.bind(null,l));return r},o.n=function(e){var t=e&&e.__esModule?function(){return e["default"]}:function(){return e};return o.d(t,"a",t),t},o.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},o.p="/";var i=window["webpackJsonp"]=window["webpackJsonp"]||[],u=i.push.bind(i);i.push=t,i=i.slice();for(var c=0;c<i.length;c++)t(i[c]);var f=u;n.push([0,"chunk-vendors"]),r()})({0:function(e,t,r){e.exports=r("56d7")},"034f":function(e,t,r){"use strict";r("85ec")},"56d7":function(e,t,r){"use strict";r.r(t);r("e260"),r("e6cf"),r("cca6"),r("a79d");var l=r("2b0e"),a=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{attrs:{id:"app"}},[r("LogTable",{attrs:{msg:"log for k8s"}})],1)},n=[],o=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticStyle:{width:"100%"}},[r("el-table",{staticStyle:{width:"96%","margin-left":"2%"},attrs:{data:e.tableData,"span-method":e.objectSpanMethod}},[r("el-table-column",{attrs:{fixed:"",prop:"Name",label:"节点名"}}),r("el-table-column",{attrs:{formatter:e.cellFormatter,label:"Cpu/Mem"}}),r("el-table-column",{attrs:{formatter:e.cellFormatter,label:"最大值"}}),r("el-table-column",{attrs:{formatter:e.cellFormatter,label:"最小值"}}),r("el-table-column",{attrs:{formatter:e.cellFormatter,label:"幅度"}}),r("el-table-column",{attrs:{formatter:e.cellFormatter,label:"平均值"}}),r("el-table-column",{attrs:{formatter:e.cellFormatter,label:"最快增幅"}}),r("el-table-column",{attrs:{formatter:e.cellFormatter,label:"硬盘使用情况"}})],1)],1)},i=[],u=(r("ac1f"),r("1276"),r("bc3a")),c=r.n(u);function f(e){return""==e?"":e.split("T")[1].split(".")[0]}var s={name:"LogTable",data:function(){return{tableData:[]}},created:function(){var e=this;c.a.get("https://klog.ciiplat.com/nodes").then((function(t){console.log(t.data),e.tableData=t.data})).catch((function(e){console.log(e)}))},methods:{objectSpanMethod:function(e){e.row,e.column;var t=e.rowIndex,r=e.columnIndex;if(0===r||7===r)return t%2===0?{rowspan:2,colspan:1}:{rowspan:0,colspan:0}},cellFormatter:function(e,t,r,l){if(l%2===0){if("Cpu/Mem"==t.label)return"Cpu";if("最大值"==t.label)return e.CpuSumMax+"% ("+f(e.CpuSumMaxTime)+")";if("最小值"==t.label)return e.CpuSumMin+"% ("+f(e.CpuSumMinTime)+")";if("幅度"==t.label)return e.CpuVolatility+"%";if("平均值"==t.label)return e.CpuSumAvg+"%";if("最快增幅"==t.label)return e.CpuMaxRatio+"%";if("磁盘使用情况"==t.label)return e.DiskUsed+"Gi/"+e.DiskTotal+"Gi "+e.DiskUsedRatio+"%"}if(l%2===1){if("Cpu/Mem"==t.label)return"Mem";if("最大值"==t.label)return e.MemMax+"Gi ("+f(e.MemMaxTime)+")";if("最小值"==t.label)return e.MemMin+"Gi ( "+f(e.MemMinTime)+")";if("幅度"==t.label)return e.MemVolatility+"Gi";if("平均值"==t.label)return e.MemAvg+"Gi";if("最快增幅"==t.label)return e.MemMaxRatio+"%";if("磁盘使用情况"==t.label)return e.DiskUsed+"Gi/"+e.DiskTotal+"Gi "+e.DiskUsedRatio+"%"}}}},p=s,m=r("2877"),b=Object(m["a"])(p,o,i,!1,null,"28e28dee",null),d=b.exports,h={name:"App",components:{LogTable:d}},M=h,v=(r("034f"),Object(m["a"])(M,a,n,!1,null,null,null)),g=v.exports,y=r("5c96"),w=r.n(y);r("0fae");l["default"].config.productionTip=!1,l["default"].use(w.a),new l["default"]({render:function(e){return e(g)}}).$mount("#app")},"85ec":function(e,t,r){}});
//# sourceMappingURL=app.4f7dc145.js.map