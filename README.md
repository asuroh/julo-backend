# julo-backend 


## Requirement
- install golang
- install redis 
- install postgresql
- install postman
- rabbit MQ


## Documentasi
Step 1 
create database run query : 
```
    CREATE database julo;
```

create table and insert data dummy :
```
 run db in file file\db.sql
```

Step 2
- configuration env

Step 3
```
- Run the server
```

go mod vendor
cd server
go run main.go
```

Step 4 
````
- Run the mqtt 
````
cd amqp_listener_update_balance 
go run main.go
````

Step 4
- Configuration Postman
Postman collection : 
    in file julo-backend.postman_collection.json
Postman Env : 
    in file Local.postman_environment.json
