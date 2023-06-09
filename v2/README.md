# 路由树总结

- 已经注册的路由，无法被覆盖
- `path`必须以`/`开头，并且结尾不能有`/`，中间也不允许有连续的`/`
- 不能再同一个位置注册不同的参数路由
- 不能在同一个位置同时注册通配符匹配和路径参数匹配
- 同名路径参数，在路由匹配的时候，值会被覆盖

>最后一条可以考虑在注册路由的时候强制`panic`，基本用户群体是程序员，程序员理解程序员。


## 为什么在注册路由的时候使用`panic`
- 可以返回`error`
- 但是从框架的角度来说，用户必须先注册完路由，才能启动服务，那么就可以在服务启动前把一切问题`panic`掉。


## 路由树是线程安全的吗？《药医不死病，佛渡有缘人》
- 路由树使用了`map`，显然不是线程安全的
- 这里没有必要考虑是否线程安全，上面也说了，基本这个是只会在启动前就要做好的事情，不存在并发
场景；至于运行期间动态注册路由，没有必要支持。
- 并发读没有问题【并发读写有问题】
---
要做：封装 HTTPServer，加锁 -> 性能下降

即：约定大于配置，遵从约定即可

## 要点
- 路由树算法核心就是前缀树。前缀的意思就是将 2 个节点的共同前缀抽取出来作为父节点，使用`/`来进行切割
- 路由匹配的优先级：这里设计的是静态匹配 > 路径参数 > 通配符匹配
- 路由查找会回溯吗？这里是不支持的，简单来说：鸡肋
- web 框架是怎么组织路由树的？一个 HTTP 方法一棵路由树，也可以考虑一棵路由树，每个节点标记自己支持的 HTTP 方法。
- 路由查找的性能受什么影响？核心是看树的高度，次要因素是路由树的宽度，从正常开发角度来说，其实两者相差不大，基本不会有很多路由的情况发生，所以基本使用宽度的就可以了
- 路由树不是线程安全的，这是为了性能，约定俗称的都会在服务启动前注册好路由；如果有运行区间动态添加路由的需求，只需要理由装饰器模式，将一个线程不安全的封装为线程安全的路由树
- 具体匹配方式的实现原理：划定优先级，然后根据静态匹配、通配符匹配、路径匹配一个一个匹配过去。