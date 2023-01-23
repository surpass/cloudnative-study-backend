

CREATE KEYSPACE sport WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 1};


use sport;




INSERT INTO rundata (id,lon,lat,speed,rtime,utime,ctime) values('1','123.123','46.456',1.23,123123123,222323323,2342342342)


mysql:

CREATE DATABASE IF NOT EXISTS sport DEFAULT CHARACTER SET utf8mb4  DEFAULT COLLATE utf8mb4_general_ci;
 


curl -H "Content-Type: application/json" -X POST -d '{"lon": "123.123123", "lat":"46.234234", "speed":12.3, "curtime":1668063277884 }' "https://www.easyolap.cn/api/v1/location"


查询所有数据
https://www.easyolap.cn/api/v1/data/index
https://www.easyolap.cn/api/v1/person


curl -H "Content-Type: application/json" -X POST https://www.easyolap.cn/api/v1/person/add -d '{"userId": 1, "userName":"surpass","gender":1,"email":"surpass_li@aliyun.com","createDate":"2022-03-07T17:11:18+08:00" }'

查找
https://www.easyolap.cn/api/v1/person/id/1

curl -H "Content-Type: application/json" -X GET https://www.easyolap.cn/api/v1/person/id/13


更新
https://www.easyolap.cn/api/v1/person/uid

curl -H "Content-Type: application/json" -X PUT https://www.easyolap.cn/api/v1/person/uid -d '{"userId": 1, "userName":"test" }'



https://www.easyolap.cn/api/v1/person/did?id=1


curl -H "Content-Type: application/json" -X POST "https://www.easyolap.cn/api/v1/location"
 -d '{"lon": "123.123123", "lat":"46.234234", "speed":12.3, "curtime":1668063277884 }' 



https://www.easyolap.cn/api/v1/data/index


curl -H "Content-Type: application/json" -X POST -d '{"lon": "123.123123", "lat":"46.234234", "speed":12.3, "curtime":1668063277884 }' "https://www.easyolap.cn/api/v1/location"


生成行程
curl -H "Content-Type: application/json" -X POST https://www.easyolap.cn/api/v1/location/trip -d '{"uid": 1, "stime":"2022-03-07T17:11:18+08:00" }'



查看详细轨迹：
curl -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJEakFsOWZaVk5feDd4T0doN092WEZlT0NDM0NWOHRfb3p4VmZjb1BCM0RFIn0.eyJleHAiOjE2NzE4ODMyNTAsImlhdCI6MTY3MTg3OTY1MCwianRpIjoiYzViOTg4MmUtMTQxYy00YTFiLWJlYTYtZjQ0ZjIzZWVlYzBlIiwiaXNzIjoiaHR0cHM6Ly93d3cuZWFzeW9sYXAuY24vYXV0aC9yZWFsbXMvZ29sYW5nIiwiYXVkIjoiYWNjb3VudCIsInN1YiI6IjNjN2Q0ZGQyLTA1OTAtNGE1ZC04ZWU2LTIxMTBkYmY1ODU3ZiIsInR5cCI6IkJlYXJlciIsImF6cCI6ImJlZWdvIiwic2Vzc2lvbl9zdGF0ZSI6IjA2ZjNlNzk2LTNiYWYtNDk5ZS1iNWQzLTYwMGZjN2JkN2JjNyIsImFjciI6IjEiLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsiZGVmYXVsdC1yb2xlcy1nb2xhbmciLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgZW1haWwgcHJvZmlsZSBvZmZsaW5lX2FjY2VzcyIsInNpZCI6IjA2ZjNlNzk2LTNiYWYtNDk5ZS1iNWQzLTYwMGZjN2JkN2JjNyIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwicHJlZmVycmVkX3VzZXJuYW1lIjoiMTM0Nzg4NzU2OTUifQ.gR6ZwWVDHfeIcgnaCehUFL6BD-wgkx2zBsMZkLoH7uvokZzMaRwUO0fjoQDiQ8KrbHjFJeCkx6c9Cl-ylhMnzd3BkjwJkhooRlJZnNoSXYY42VFT5VJe9qdn5v8m_xi2mMF9qxbeGI3zbJQkZMlUDJhoFSH69AZCF34j1lW_8YlG1J9IDDn4DrPst2duuSkA7M3o60JCF3SZfrp8YRv1PetuViFd4TQpdIjTATvuMlEXGy17yeA0Dfn7oCEwhRNPyN6TH-77VqSldI-um3ca_UEutFwqYlynC70TukvzOweea8es9TwlBGDKHPoKHY-1DDCXrvGTZKY5ePj6Y6hfgg" -X GET https://www.easyolap.cn/api/v1/location/getLocationByTid -d '{"tid":"8afc894f-b6f1-41d8-8b1d-5741ffbb9d19"}'