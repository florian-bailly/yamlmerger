# About

SimpleYaml brings you basic YAML parsing, traversing, writing and merging.

With the merger it's possible to have a base config, and one for each environment with its specificities.
Thus, it avoids duplication and improves maintainability.

It doesn't pretend to be a real parser or handle edge cases. However, any contribution is welcomed.


# Motivation

The initial need for this is the lack of Docker Compose to offer extended override to its config files.


# Quick start

```sh
go run yamlmerger.go -i "tests/input1.yml tests/input2.yml" -o "merged.yml" -dpl="args:=,volumes::" -del-tk="nil"
```


# Examples

### Merge

Input 1:
```yml
services:
  php:
    build:
      context: docker/php
      args:
        - FINAL_BASE=dev
        - USER_ID=${USER_ID:-1000}
    networks:
      - web
      - db
    volumes:
      - .:/mnt/app
      - ./storage/logs/php:/mnt/log
    environment:
      XDEBUG_ENABLE: true
      OPCACHE_ENABLE: ${OPCACHE_ENABLE}
```

Input 2:
```yml
services:
  php:
    build:
      context: docker/php_prod
      args:
        - FINAL_BASE=nil
    networks:
      - web:nil
    volumes:
      - .:/mnt/app_prod
    environment:
      XDEBUG_ENABLE: false
```

Ouput:
```yml
services:
  php:
    build:
      context: docker/php_prod
      args:
        - USER_ID=${USER_ID:-1000}
    networks:
      - db
    volumes:
      - .:/mnt/app_prod
      - ./storage/logs/php:/mnt/log
    environment:
      XDEBUG_ENABLE: false
      OPCACHE_ENABLE: ${OPCACHE_ENABLE}
```
