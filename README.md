这是一个自动更新公网IP到cloud flare的小工具。

需要配置如下环境变量：

|  环境变量   | 说明  |
|  ----  | ----  |
| API_TOKEN  | 域名的API token, 需要自己在域名中创建，并且拥有更新DNS record的权限|
| ZONE_ID  | 区域ID，在域名概述页右下角可以获取到|
| RECORD_ID  | 解析记录ID，可以通过打开浏览器调试工具，并更新一条DNS解析获取到 |
| INTERVAL  | 修改IP的频率，单位分钟，建议5或更大的数值|
| NAME  | 二级域名名称|

当然，我们也支持`.env`文件，你需要把你的配置放在`/home`路径下。
`.env`文件如下：
```
API_TOKEN=YOUR_API_TOKEN_HERE
ZONE_ID=DOMAIN_ZONE_ID
RECORD_ID=DNS_RECORD_ID
INTERVAL=10
NAME=SUBDOMAIN
PROXIED=false
```
