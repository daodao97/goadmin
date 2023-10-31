按钮

按钮点击的几种动作:  页面跳转,  打开模态框,  请求后端接口

```json
{
    "type":  "jump",  //  按钮类型, 可选值 jump / modal / api
    "target":  "/user/{id}",  //  目标
    "text":  "编辑",  //  按钮文本
    "props": {
        "type":  "primary",  //  按钮类型 primary / success / warning / danger / info / text
        "icon":  "el-icon-plus"  //  按钮图标
    },
    "extra": {   //  modal / api / form/ table 时的 扩展属性

      "actionHandler": "@topList" // @deleteList @unshiftList @updateList   按钮为 rowButton 或 batchButton 时 几个对当前表的操作

    },  
    //  status, age 的取值(元数据 metaData)一般为列表的行数据
    //  当 当前行 数据满足 when 定义的规则是 才会显示该按钮
    "when": [ ["status",  "=", 1], ["age",  "<", 10] ],
    //  or 单个条件  "when": ["status",  "=", 1]
}
```

actionHandler 按钮对表数据的操作

rowButton  

@topList   置顶当前行数据
@deleteList 删除当前行数据
@unshiftList   在表最上面插入新数据(新数据一般来源于按钮弹窗表单)
@updateList     更新当前行数据, 数据一般来源于当前行的编辑

batchButton

@topList   置顶选中的行
@deleteList   删除选定的行

按钮类型

跳转

```json
{
    "type":  "jump",
    "target":  "/user/{id}",  //  如  /user/123  对应前端的  /user/:id  页面路由
    "text":  "编辑",
}
//  target 也可以时一个 url, 点击后将在新页面中打开该地址
```

打开模态框

```json
{
    "type":  "modal",
    "text":  "编辑",
    "extra": {    //  模态框内部 可渲染 列表/表单 等组件
        "formItems": []  //  同 表单
        "saveApi":  "/xxx/xx"  //  表单保存接口
    }
}
```

```json

{
        "props":{
                "type":"info"
        },
        "extra":{
                "headers":[
                        {
                                "field":"name",
                                "label":"名称"
                        }
                ],
                "listApi":"/user/list"
        },
        "text":"表格",
        "type":"table"
}
```


```json
{
        "props":{
                "type":"info"
        },
        "extra":{
                "formItems":[
                        {
                                "field":"name",
                                "label":"名称"
                        }
                ],
                "saveApi":"/user/create"
        },
        "text":"表单",
        "type":"form"
}
```

请求后端接口

```json
{
    "type":  "api",
    "text":  "下架",
    "extra": {
        "url":  "/goods/offline",
        "method":  "DELETE"  //  http method, 默认为 get
    },

     "confirm": false // 是否需要二次确认, 默认true
}
```

按钮组

```json
{
        "text":"按钮组",
        "type":"group"
        "subButton": []  //  子项 为 单个按钮的 配置
}
```
