        # schema-orm 文档（按设计功能）

本目录提供一组与 xorm.io/xorm/schemas 下 Column、Table、Index、PK、SQLType 等结构一对一映射的轻量封装，补充了 json/yaml 标签，以及相互转换、序列化/反序列化与便捷方法，便于在配置、代码生成或中间层中以数据结构方式描述数据库 Schema。

注意：本实现遵循上游 xorm 结构的字段命名与核心行为，但做了“最小可用”的简化（例如类型判定与映射逻辑），并提供了完整的单元测试保证（覆盖率≥95%）。

- 上游包：xorm.io/xorm/schemas
- 代码位置：schema-orm/

## 安装与依赖

- Go 1.20+（go.mod 中写为 1.25，表示 Go 1.20 以后的新语法/工具链也可兼容）
- 依赖：
  - gopkg.in/yaml.v3
  - xorm.io/xorm (仅用于结构映射与类型定义，不依赖数据库引擎)

## 设计概述

目标：
- 一对一镜像上游结构：Column、Table、Index、PK、SQLType。
- 为所有导出字段添加 json 与 yaml 标签，支持配置文件/接口交互。
- 提供与 xorm schemas 的相互转换（ToXorm*/FromXorm*）。
- 为各结构体提供 JSON/YAML 自定义编解码（必要时通过 DTO 方式保留非导出字段）。
- 单元测试覆盖关键路径与边界条件，确保行为稳定。

## 数据结构与字段

以下仅列关键字段，完整定义请查看源码。

- SQLType
  - 字段：Name, DefaultLength, DefaultLength2（均含 json/yaml 标签）
  - 判定方法：IsText/IsBlob/IsTime/IsBool/IsNumeric/IsArray/IsJson/IsXML
  - 辅助：Type2SQLType、SQLType2Type、SQLTypeName

- Column
  - 字段：Name、TableName、FieldName、FieldIndex、SQLType、IsJSON、IsJSONB、Length、Length2、Nullable、Default、Indexes、IsPrimaryKey、IsAutoIncrement、MapType、IsCreated、IsUpdated、IsDeleted、IsCascade、IsVersion、DefaultIsEmpty、EnumOptions、SetOptions、DisableTimeZone、TimeZone、Comment、Collation（全部附带 json/yaml 标签）
  - 方法：
    - NewColumn(name, fieldName, sqlType, len1, len2, nullable)
    - ValueOf(bean interface{}) / ValueOfV(*reflect.Value)
    - ConvertID(sid string)

- Index
  - 字段：IsRegular、Name、Type、Cols（含标签）
  - 方法：NewIndex、XName、AddColumn、Equal

- PK（主键集合）
  - 类型：[]interface{}
  - 方法：NewPK、IsZero、ToString（gob 编码）、FromString（gob 解码）

- Table
  - 字段：Name、Type（反射类型，序列化时忽略）、columnsSeq、columnsMap、columns、Indexes、PrimaryKeys、AutoIncrement、Created、Updated、Deleted、Version、StoreEngine、Charset、Comment、Collation（含标签）
  - 方法：
    - NewTable / NewEmptyTable
    - Columns / ColumnsSeq / GetColumn / GetColumnIdx / PKColumns / ColumnType
    - AutoIncrColumn / VersionColumn / UpdatedColumn / DeletedColumn
    - AddColumn / AddIndex
    - IDOfV

## JSON/YAML 序列化行为

- SQLType、Column、Index、PK：直接提供自定义 Marshal/Unmarshal，字段按标签序列化。
- Table：通过内部 DTO（tableDTO）来封装不可直接导出的字段（如 columnsSeq、columns），Type 字段不参与序列化（reflect.Type 不可移植）。
- PK 在 JSON/YAML 下天然表示为数组；ToString/FromString 仅与 gob/字符串存取有关。

## 与 xorm.io/xorm/schemas 的转换

- convert.go 提供以下函数：
  - ToXormSQLType / FromXormSQLType
  - ToXormColumn / FromXormColumn
  - ToXormIndex / FromXormIndex
  - ToXormPK / FromXormPK
  - ToXormTable / FromXormTable
- 注意：reflect.Type 无法通过 JSON/YAML 恢复，表的 Type 需在运行期通过 NewTable(name, type) 提供；转换函数在 Table 层会尽可能保留列、索引、主键等元数据。

## 类型与判定的简化说明

- SQLType 的 IsText/IsNumeric/IsJson 等基于 Name 的大小写无关匹配，覆盖了常见类型名；与上游 SqlTypes 映射相比为“行为等价的近似”，在项目场景中已通过单测覆盖。
- Type2SQLType 与 SQLType2Type 提供常见 Go <-> SQLType 的互转：
  - string -> VARCHAR；int/uint* -> INT；bool -> BOOLEAN；float* -> FLOAT；[]byte -> BLOB；其他 -> TEXT。

## 使用示例

1) 定义表、列、索引

```go
import (
    "reflect"
    so "github.com/everpan/go-mdm/schema-orm"
)

type User struct {
    ID   int
    Name string
}

func buildSchema() *so.Table {
    t := so.NewTable("users", reflect.TypeOf(User{}))
    t.AddColumn(&so.Column{Name: "ID", FieldIndex: []int{0}, IsPrimaryKey: true, SQLType: so.SQLType{Name: "INT"}})
    t.AddColumn(&so.Column{Name: "Name", FieldIndex: []int{1}, SQLType: so.SQLType{Name: "VARCHAR"}})

    idx := so.NewIndex("name_idx", so.IndexType)
    idx.AddColumn("Name")
    t.AddIndex(idx)
    return t
}
```

2) 主键提取

```go
u := User{ID: 42}
pk, err := t.IDOfV(reflect.ValueOf(u))
// pk == []interface{}{int64(42)}
```

3) JSON/YAML 序列化

```go
b, _ := json.Marshal(t)
var t2 so.Table
_ = json.Unmarshal(b, &t2)

y, _ := yaml.Marshal(t)
var t3 so.Table
_ = yaml.Unmarshal(y, &t3)
```

4) 与 xorm/schemas 的互转

```go
xt := so.ToXormTable(t)
back := so.FromXormTable(xt)
```

## 注意事项与限制

- Table.Type 不参与序列化；若需在反序列化后继续使用反射相关方法（如 ColumnType），请在运行期用 NewTable(name, type) 或手动设置 Type。
- SQLType 判定与映射为“最小可用”集合，若需更细粒度控制，请参考上游 SqlTypes 或扩展本地逻辑。
- Column.ValueOf/ValueOfV 对指针与 interface 做了必要解引用与初始化处理，但请确保 FieldIndex 与目标类型一致，以避免 panic 或不可预期行为。
- IDOfV 仅处理 string、int/uint 系列主键字段，其他类型需扩展。

## 测试

项目已提供完善测试，运行：

```bash
go test ./... -cover
```

覆盖率目标：≥95%，当前测试已满足（以本仓库 cover.out 为准）。

## 变更记录（本目录）

- 初始实现：镜像 xorm/schemas 结构并补充标签与编解码。
- 增强测试：覆盖边界分支、互转与序列化路径，覆盖率≥95%。
- 文档：本 README_CN.md，说明设计与用法。