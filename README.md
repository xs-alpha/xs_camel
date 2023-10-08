# xs_camel

使用ast抽象语法树解析sql建表语句，并将建表语句的字段属性，类型和注释抽离出来，动态绘制GUI界面，使用Fyne GUI框架

# 功能
- 解析sql建表语句
- 开启监听后，监听剪贴板的英文和下划线连接的字符串，自动转驼峰，方便书写开发文档
- 将解析的建表语句绘制出GUI多选框，每个接口可以选择不同的参数
- 转成excel格式，方便 复制到开发文档 中
- 时间戳，base64, sha256, 字数统计， 随机密码生成，二维码解析，二维码生成，json美化，文件md5, 随机选择等小工具集合


# 注意
- 最新分支sql_parse
- v0.3发布版本sql解析转换的是首字母大写版本的驼峰，v0.31版本是首字母小写版本驼峰
- v0.4版本将首字母大写版本和小写版本都抽离出来，按钮控制
- v0.5版本新增大小写转换
- v0.6版本 内存优化& 日志打印可选
- v0.7新写c#二维码增强工具，用go调用（会从第三方网盘直链下载，介意的话可以不使用，吾爱上有哈勃查毒直链 https://www.52pojie.cn/thread-1839818-1-3.html）
  -  由于go的开原二维码解析大多数不成熟，并且存在部分二维码无法解析的问题，特意增加了一个用c#写的二维码解析的小工具 
- v0.8选择困难症福音     

