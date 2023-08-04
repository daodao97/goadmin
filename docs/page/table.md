
列表组件

列表的 JsonSchema 定义为:

```json
{
  "headers": [], // 表头, 具体详细下方列类型定义
  "listApi": "/user/list", // 列表数据接口
  "list": [{"k":"v"}], // 列表数据, 在无 listApi 时, 可以指定数据
  "infoApi": "/user/schema", // 列表渲染配置拉去接口, 响应数据同 列表组件的 JsonSchema
  "filter": [], // 同 表单组件的 JsonSchema 筛选条件
  "batchButton": [], // 批量操作按钮, 表头左上方
  "normalButton": [], // 页面操作按钮, 表头右上方
  "rowButton": [], // 行按钮, 表头操作列
  "showPagination": true, // 是否显示分页组件
  "showFilter": true, // 是否显示 搜索条件表单
  "selectedNotice": {
    "text": "当前勾选了 {selectedCount} 条数据, 平均值为: {score|sum|avg}",
    "position": "afterBatchButton" // 显示位置 afterBatchButton 批量按钮后(默认), beforePagination 表单控件前
  },  // 可以缺省为 "selectedNotice": "{selectedCount} 条数据"
  "tableProps": {}, // 原生组件的属性 https://element-plus.gitee.io/zh-CN/component/table.html#table-attributes
  "tabs": [
    {
      "value": 1,
      "field": "mytype",
      "label": "选项卡名",
      "icon": "" // 图标
    } // 当 当前选项卡选中时, 数据请求参数为 /user/list?mytype=1
  ],

  "orderBy": {   // 设置默认的排序规则, 非必须
      "field": "id", // 默认排序的字段
      "mod": "desc" //  desc or asc
   },

  "orderByMulit": [{   // 设置默认的排序规则, 非必须
      "field": "id", // 默认排序的字段
      "mod": "desc" //  desc or asc
   }], 多条

"exportCurrentPageAble": true, // 显示导出按钮

"dragSort": true // 列表拖动排序  true 或者 api 或者 json

// 当 dragSort = true, 列表拖动后不会有其他动作, 可以配置表单的 normalButton 最终完成数据提交

// 当 dragSort = /xx/api 列表拖动后会立即请求该接口, 进行数据提交

// 当 dragSort = {

           sortApi: /xxx/xxapi,

           confirm: true, // 默认false, 设置为true时 拖动后需要二次确认才会生效

          // 二次确认文案: 默认 '此操作将会改变数据顺序, {oldIndex} => {newIndex}, 是否确认此操作?'  支持当前行变量替换

          "confirmMsg": "[{title}] 分区 的顺序将从 {oldIndex} 调整到 {newIndex}, 是否确认此操作?", 

     } 是, 列表拖动后会立即请求该接口, 进行数据提交

}
```

说明:

headers 为必须字段
当列表为 后端拉取数据时 listApi 为必须, 页面定义为 通用CRUD 时, 会自动补充
列类型定义

```json
{
  "field": "data_field", // 数据的字段名
  "label": "表头",
  "info": "表头的提示文字",
  "type": "列类型" // 可选值为 空字符串, enmu, image 等, 具体见下方
  "sortable": false, // 是否支持排序 默认 false
  "props": {}, // 根据具体的列类型而定
    "fake": false // 是否为伪字段 默认 false, 后端脚手架用于生成 select ** from, 为 true 是将不查询该字段
}
```

排序

```json
{
    "type": "sort-index",
    "label": "排序",
    "fake": true,
    "cellProps": {
        "ctrl": false,
        "width": 100
    },
    "info": "点击单元格修改排序"
}
```

索引排序和拖动排序,  等非实时提交的形式, 一般借助按钮 提交 排序后的数据

```json
{
    "props": {
        "type": "primary"
    },
    "text": "发布",
    "type": "api",
    "extra": {
        "method": "POST",
        "url": "/modules/publish",
        "sourceData": {
            "ids": "tableList", // 将替换为 整个表单的数据 进行提交
            "page_id": "data[1].id" // 格局这个key path从行数据中获取数据并合并到 post body 中
        }
    }
}

```


文本

```json
{
  "field": "name",
  "label": "姓名"
}
// 此处没有定义 type, 默认直接渲染为文本
```

标签

```json
{
  "field": "status", // 通常对应表单中 type 为 radio/checkbox/select 类型的字段
  "label": "状态",
  "type": "enum",
  "options": [
    {"value": 1, "label": "启用"},
    {"value": 2, "label": "禁用"}
  ],
  "state": { // 标签的颜色值
    "1": "success",
    "2": "info"
    }
}
```

图片

```json
{
  "field": "cover"
  "label": "封面图",
  "type": "image",
}
// 当字段值为 http://xxxx.jpg 图片地址时 使用
```

链接

```json
{
  "field": "doc_url"
  "label": "文档地址",
  "type": "link",
}
// 当字段值为 http://xxxx 可访问链接时使用
```

html
```json
{
  "field": "html_code"
  "label": "一段html",
  "type": "html",
}
// 当字段为 一段html 代码时 使用
```

字面量模板
```json
{
  "field": "test_field"
  "label": "信息",

    "type": "tpl",
  "tpl": "姓名: {name}, 年龄: {age}", // name, age 为 list 子项的 key
}
```

json
```json
{
  "field": "comm_conf"
  "label": "配置",
  "type": "json",
}
// 当内容为 复杂的kv数据时 直接渲染使用
```

图表
```json
{
  "field": "matched_data",
  "label": "匹配数据",
  "fake": true,
  "type": "chart",
  "cellProps": {  // 图表控件的属性, 详见 图表 章节
    "options": {
      "series": [ // 数据序列
        { "name": "总量", "type": "line"},
        { "name": "命中数", "type": "line"}
      ]
    }
  }
}
```

icon
```json
{
  "field": "icon",
  "label": "图标",
  "type": "icon"
}
```

富文本 rich
```json
{
  "field": "data",
  "label": "简易富文本",
  "type": "rich"
}
// 字段值格式为 "今天是 {date|yellow|back}, 天气: {weather|前景色|背景色}"
```

列的编辑模式

```json
{
  "edit": true, // 当 定义 edit=true 时, 列表单元格将以 表格的形式渲染
  // 以下字段同表单字段
  "type": "cascader", // type 可选值同 表单控件
  "field": "role_ids",
  "label": "角色",
  "props": {
    "disabled": true,
    "optionsApi": "/role/tree",
    "props": {
      "checkStrictly": true,
      "multiple": true
    },
    "style": {
      "width": "100%"
    }
  }
}
```

列表关联数据

假设当前列表 的底层表 为 A, A.bid 为 B 表主键 B.id, 现在 A 的列表中需要显示 B.name 可以使用关联数据配置, 快速实现

```json
{
	"hasOne": [
		"bilibili_tv.B:bid->id,name as bname"
	],
	"headers": [
		{
			"field": "id",
			"label": "ID",
			"hidden": true // 不需要显示时可以隐藏
		},
		{
			"field": "bid",
			"label": "BID",
			"hidden": true // 不需要显示时可以隐藏
		},
		{
			"field": "bname",
			"label": "BNAME"
			"fake": true // 标明为伪字段
		},
		{
			"field": "bname",
			"label": "BNAME",
			"hasMany": "bilibili_tv.B:bid->id,name as bname" // 一对多的关联数据, 语法同 hasOne
		}
	]
}
```

 以上配置即可快速展示关联数据, hasOne的具体语法为: `[pool.]db.table:[local_key→]foreign_key,other_key`

pool: db的连接名,  见 db.toml database.name 节点, 默认为 default, 可缺省

db : database 库名

table: 表名

local_key: 本表用于关联 的字段, 默认为 id, 可缺省

foreign_key: 关联表中用于 join 的字段名

other_key: 关联表中需要 查出的字段, 可以多个, 可以使用别名

