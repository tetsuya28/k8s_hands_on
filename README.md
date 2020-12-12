# Kubernetes Hands-on

## Pre-requirement
### API
* api/
```
cp api/.env{.sample,}
cp api/.env.docker{.sample,}
```

### Frontend
* ui/
```
cp ui/.env{.sample,}
```

### Database
* db/
```
cp db/.env{.sample,}
```

## How to run
### All in docker-compose(Recommended)
- Start `docker-compose`
  ```
  $ docker-compose up -d
  ```
### API local
- Install `realize`
- Start `docker-compose` for DB.
  ```
  $ docker-compose up -d db
  ```
- Run API
  ```
  $ make run
  ```
