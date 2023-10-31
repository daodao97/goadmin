表单控件

表单组件的通用字段,  通用字段  即所有表单组件支持的属性, 标明  必须的必须字段, 未标记的则为非必须.

```json
{
    "type":  "input",  //  表单控件的类型, 可选值为下方具体控件列表, 必须
    "field":  "name",  //  字段名一般为数据库中的 column_name, 必须
    "label":  "名称",  //  列表渲染是的可读标记名, 供运营识别, 必须
    "info":  "请注意, 此处是表单填写时的提示语",  //  提示文字, 辅助说明字段的语义
    "rules":  "required|max:10",  //  字段规则, 约束字段为 必须/最多x个字符 等, 具体见方法表单规则
    "props": {    //  表单组件的扩展属性, 一般可以参加 element-plus 的原生组件
      }
    "validate":  "required,min=0,max=100",  //  此处规则约束为后端脚手架使用, 值参见https://github.com/asaskevich/govalidator
    //  tips: rules 为前端组件调用, validate 为后端脚手架数据验证, 理想模式下 前后端应该使用同一份约束, 后期计划统一
    "disabled":  false,  //  禁用状态
    "form":  true,  //  当 form=true  时, 后端脚手架会从db中查询该字段, 当 db 中无该字段, 又需要在表单中显示时, 设置为  false
    "placeholder":  "输入框默认提示",

   "value":  111,

   "submitConfirm": true// 值类型为 bool or object as {title: '', message : ''}

}
```



说明

type,  field,  label  为必须字段, 其他为非必须
输入框 input

```json
{
    "type":  "input",
    "field":  "name",
    "label":  "名称",
    "props": {
        "showWordLimit":  true,  //  显示剩余可输入字数
        "showCopy":  false,  //  显示 复制到粘贴板 控件
        "mask":  "99-999999"//  文本规则, 值规则参见 https://github.com/RobinHerbots/Inputmask
    }
}
```

数字 number

```json
{
    "type":  "number",
    "field":  "age",
    "label":  "年龄",
        "rules":  "max:100",
    "props": {
        "min": 1,  //  最小值
        "max": 100,  //  最大值
        "step": 5,  //  步长
        "precision"  : 0  //  数据精度, 默认为整数
    }
}
```

数字区间 number-range

```json
{
    "type":  "number-range",
    "field":  "age",
    "label":  "年龄区间",
        "rules":  "max:100",
    "props": {}
}
```

子表单 sub-form

```json
{
    "type":  "sub-form",
    "field":  "website",
    "lable":  "系统模块",
    "props": {
        "formItems": [],  //  同表单的 控件配置
        "options": {},  //  同表单的 样式通知, 无 提交/取消 按钮
        "repeat":  false,  //  为  true  表示为 可重复的子表单
    }
}
```

下拉框 select

```json
{
        "type":  "select",
    "field":  "type",
    "lable":  "类型",
    "options": [  //  备选项列表
        {"value": 1,  "label":  "备选项1",  "xxx":  "2", "when": ["field_a", "=", 1]}, // 备选项支持 when 条件过滤
        {"value": 2,  "label":  "备选项2",  "xxx":  "1", "previewHtml": "<span>当鼠标悬浮在这个备选项时, 可以显示一些额外的辅助信息</span>" }
    ],
    "props": {
        "selectApi":  "/xxx/types",  //  备选项为远程响应, 接口响应数据的格式同 options

              "immediate": false,    // 是否立即调用  selectApi 进行下拉数据请求, 默认false(当有用户输入是才请求), 若为true 则会立即请求 selectApi 进行下拉数据渲染(此处要注意默认返回的数据不能过多, 否则会造成前端渲染卡顿)
        "multiple":  false,  //  是否为多选
        "multipleLimit": 10,  //  多选时可选项目的个数上限
        "allowCreate":  false,  //  允许创建

            "value":"value" // 备选项
        "valueKey":  "value",  //  备选项数据的 value key, 当备选项数据不是标准的 value/label  时使用
        "labelKey":  "label",  //  备选项数据的 label key, 当备选项数据不是标准的 value/label  时使用
        "tpl":  "",  //  label 的字面量显示模板, 默认无 例如  'ID:{id}|名称:{label}'  {xxx} 为 options 的元素key
        "fullWidth":  false,  //  控件宽度是否撑满屏幕
        "effectData": {}  //  当前字段变动时的副作用, 格式为 {k:v} 可以覆盖 与字段 `type` 同级的同名字段的值,
        //  具体值为 merge(effectData, options.{value})
        //  若定义了 effectData 为 {"a":1,  "xxx":  ""}, 当前选中 value=1时 则覆盖的实际值为 {"a":1,  "xxx":  "1"

            "watch": "field_name" // ["field_name"] 当前组件在form某些字段变化时, 需要重新渲染的时候可以 使用观察者模式
    }
}
```

单选 radio

```json
{
    "type":  "radio",
    "field":  "status",
    "label":  "状态",
    "options": [
        {"value": 1,  "label":  "启用"},
        {"value": 2,  "label":  "禁用"}
    ]
}
```

多选 checkbox

```json
{
    "type":  "checkbox",
    "field":  "type",
    "label":  "支持类型",
    "options": [
        {"value": 1,  "label":  "类型1"},
        {"value": 2,  "label":  "类型2"},
        {"value": 3,  "label":  "类型3"}
    ]
}
```

// 备选项远程模式

```json
{
         field: 'a',
         label: 'checkbox',
         type: 'checkbox',
         props: {
            valueKey: 'page_id',
            labelKey: 'title',
            optionsApi: '/tv_channel_region/list?state=2'
         }
      }
````

日期 date

```json
{
    "type":  "date",
    "field":  "mydate",
    "lable":  "日期",
    "props": {  //  具体参见 https://element-plus.gitee.io/zh-CN/component/date-picker.html#attributes
        "type":  "date"  //  可选  date  / daterange, 默认  date
    }
    //  日期组件的值 默认格式为 2021-01-01
}
```

时间 time

```json
{
    "type":  "time",
    "field":  "mytime",
    "lable":  "时间",
    "props": {  //  具体参见 https://element-plus.gitee.io/en-US/component/time-picker.html#attributes
    }
}
```

日期时间 datetime

```json
{
    "type":  "datetime",
    "field":  "mydate",
    "lable":  "日期时间",
    //  值得默认格式为 2021-01-01 01:01:01
}
```

图片上传 image

```json
{
    "type":  "image",
    "field":  "cover",
    "label":  "封面",
    "props": {
        "action":  "/api/upload",  //  上传接口 默认  /api/upload
        "limit": 1,  //  上传文件数, 默认1
        "format": [],  //  支持的 扩展类型, 默认无, 例如 ["jpg",  "png"]
        "maxSize": 0,  //  文件最大 * kb, 默认无限制
    }
}
```

文件上传 file

```json
{
    "type":  "file",
    "field":  "mid_file",
    "label":  "文件上传",
    "props": {
        "action":  "/api/upload",  //  上传接口 默认  /api/upload
        "limit": 1,  //  上传文件数, 默认1
        "format": [],  //  支持的 扩展类型, 默认无, 例如 ["jpg",  "png"]
        "maxSize": 0,  //  文件最大 * kb, 默认无限制
    }
}
```

JSON json

```json
{
    "type":  "json",
    "field":  "comm_conf",
    "label":  "自定义配置"
}
```

代码编辑器 code

```json
{
    "type":  "code",
    "field":  "sql",
    "label":  "自定义sql",
    "props": {
        "mod":  "sql",  //  可选 sql/html/php/go/javascript/...
        "options": {}  //  https://microsoft.github.io/monaco-editor/api/interfaces/monaco.editor.ieditoroptions.html
    }
}
```

规则引擎 rule

```json
{
    "type":  "rule",
    "field":  "bu_rule",
    "label":  "业务下发规则",
    "props": {
        "formItems": [  //  同表单控件
            {
                "type":  "select",
                "label":  "年级",
                "options": [
                    {  "value": 1,  "label":  "一年级"},
                    {  "value": 2,  "label":  "二年级"}
                ],
                "operatorOptions": [  //  此处增加 操作符 备选项, 非必须, 默认 =,<,>,>=,<=,!=,in,not  in,has,not has
                    {  "value":  "=","label":  "是"}
                ]
            },
            {
                "type":  "number",
                "field":  "age",
                "label":  "年龄"
            }
        ]
    }
}
```

关于  规则引擎  的详细描述  doc

KV数据 filter

```json
{
    "type":  "filter",
    "field":  "bu_rule",
    "label":  "业务下发规则",
    "props": {
        "formItems": []  //  同表单控件
    }
}
```
  
//  数据结构为 {k:v}

数据展示 show

当需要根据当前表单数据展示某些数据时使用

```json
{
    "type":  "show",
    "field":  "__show__",
    "label":  "人群量级",
    "form":  false,
    "props": {
        "tpl":  "当前命中 {count} 人",  //  count 的值可以来源于当前表单, 也可以来源于 dataApi 的响应
        "dataApi":  "/xxx/xxx?type={type}"  //  如果定义了 dataApi 将从后端拉去数据供渲染模板使用,
        "watch":  "watch_field_name",  //  被监听的字段名, 当其变化是会刷新 show 组件的显示
    }
}
```

开关 switch

```json
{
    "type":  "switch",
    "field":  "status",
    "label":  "开关",
    "props": {}  //  https://element-plus.gitee.io/en-CN/component/switch.html#attributes
}
```

级联选择器 cascader

```json
{
    "type":  "cascader",
    "field":  "role_ids",
    "label":  "用户角色",
    "props": {    //  https://element-plus.gitee.io/zh-CN/component/cascader.html#cascader-attributes
        "optionsApi":  "/role/resource",  //  备选数据接口

              "ctrl": true, // default is false, 显示便捷操作控件, 全选, 全不选, 反选
        "props": {    //  https://element-plus.gitee.io/zh-CN/component/cascader.html#props
            "checkStrictly":  false,  //  是否严格的遵守父子节点不互相关联
            "multiple":  true  //  是否多选
        }
    }
}
```

级联选择面板

```json
{
    "type":  "cascader-panel",
    "field":  "resource_ids",
    "label":  "角色资源",
    "props": {    //  https://element-plus.gitee.io/zh-CN/component/cascader.html#cascader-attributes
        "optionsApi":  "/role/resource",  //  备选数据接口
        "props": {    //  https://element-plus.gitee.io/zh-CN/component/cascader.html#props
            "checkStrictly":  false,  //  是否严格的遵守父子节点不互相关联
            "multiple":  true,  //  是否多选
            "emitPath":  false,  //  当子项都被选中时, 值中只保留父节点
        }
    }
}
```

滑块 slider

```json
{
    "type":  "slider",
    "field":  "percent",
    "label":  "百分比",
    "props": {}  //  https://element-plus.gitee.io/zh-CN/component/slider.html#attributes
}
```

颜色选择器 color

```json
{
    "type":  "color",
    "field":  "bg_color",
    "label":  "背景色",
    "props": {}  //  https://element-plus.gitee.io/zh-CN/component/color-picker.html#attributes
}
```

穿梭框 transfer

```json
{
    "type":  "transfer",
    "field":  "active_data",
    "label":  "活跃数据",
    "props": {}  
}
```

computed 动态计算

```json
{
  "formItems": [
    {
      "field": "test"
    },
    {
      "field": "test1",
      "computed": [
        {
          "when": "1", // 约束, "1", ["=", "1"] 效果相同,  [">", "1"] 可以定义 非等值的规则
          "set": { // set 为 object, key 为 field
            "test": { // 字段的重写项, 值为 formItems[xxx] 的控件配置
              "value": 1
            }
          }
        },
        {
          "when": "2",
          "set": {
            "test": {
              "value": 2
            }
          }
        }
      ]
    }
  ]
}
```

如上的表单配置, 当 test1 = 1 时, 会自动设置   test 字段的值为 1,   test1=2 时, 会自动设置 test 字段的只为 2,

框架中 when的完整格式 [$field, $operator, $compareValue], 如   ['a', '=', 1],   用于比对的原始数据一般为当前表单, 或当前行的数据,   支持的操作符, =, >, <, >=, <=, !=, in, not_in

动态表单

表单中经常会有联动显示的需求, 比如 当 a 字段值为1是, 显示b字段, 为2时, 显示c字段

`depend` 属性定义了当前字段依赖项, depend.field 字段 满足 depend.value 中定义的值是, 当前字段才会显示

```json
{
    "formItems":[
        {
            "field":"a",
            "label":"A",
            "type":"select",
            "options":[
                {
                    "value":1,
                    "label":"是"
                },
                {
                    "value":2,
                    "label":"否"
                }
            ]
        },
        {
            "field":"b",
            "label":"B",
            "depend":{
                "field":"a",
                "value":1
            }
        },
        {
            "field":"c",
            "label":"C",
            "depend":{
                "field":"a",
                "value":2 // 如果依赖多个值可以为数组 [1,2]
            }
        },
	   {
            "field":"d",
            "label":"d",
			"type": "sub-form",
            "props": {
				"formItems": [
         				{
           				 "field":"b",
           				 "label":"B",
           				 "depend":{
              				  "field":".a", // sub-form 中依赖 上级数据的话可以 使用 .xxx, 会从根节点查询依赖数据的值是否满足依赖
               				  "value":1
            			}
        }, 
				]
			}
        }, 
    ]
}

```
