# auth

# Migration
In the context of databases, migration refers to the process of evolving the structure of a database schema over time while preserving the existing data. Database migration typically involves making changes to the schema, such as adding new tables, modifying existing tables, creating indexes, or altering constraints.

utility goose for migration



# Auth Service

Auth Service includes 6 docker containers:
- auth - server processing authentication requests.
- postgres - permanent storage of data.
- migrator - performs migration in database using goose package.
- prometheus - scrapes application and its own metrics.
- grafana - visualizes metrics.
- jaeger - visualizes traces.

## Deploy

Make sure docker network `service-net` is in place for microservices communication. If none exists, then create network:
```
# make docker-net
```

To deploy Auth Service:
```
# make docker-deploy ENV=<environment>
```
*ENV is used then as a config name. Possible ENV values are now `stage` and `prod` as these configs are now in the repository.*

To stop Auth Service:
```
# make docker-stop ENV=<environment>
```

