#!/usr/bin/env bash
docker run -p 3306:3306 --name database -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=myaktion mariadb:10.5
