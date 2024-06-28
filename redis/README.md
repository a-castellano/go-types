# redis type

This type manages redis configs.

## Required variables

The following env variables should be definned when using this type:

* **REDIS_HOST** defines redis host, it must be a string, it's default value is "localhost" 
* **REDIS_PORT** defines redis port, it must be a valid port number, it's default value is 6379
* **REDIS_PASSWORD** defines redis password, it's default value is an empty string
* **REDIS_DATABASE** defines redis database, it's default value is 0 
