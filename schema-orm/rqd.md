1. 参考 xorm.io/xorm/schemas 的 column,table,index,pk 等结构的定义 完成一对一的复制，并对每个字段加上json，yml的tag
2. 为每个结构体构建结构体映射转化函数
3. 为每个结构体构建json序列与反序列化定义
4. 为每个结构体构建yml序列与反序列化定义
5. 为每个方法实现测试用例，并保证测试覆盖率达到95%以上
6. 所有的代码，请在 schema-orm 目录下完成