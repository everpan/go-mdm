1. 在utils目录下 构建一个xorm的对象，将xorm.io/xorm对应的方法映射到该结构内
2. 将xorm所有的公共方法通过goja映射到javascript对象
3. 做好单元测试，并确保本模块的测试覆盖率达到93%以上