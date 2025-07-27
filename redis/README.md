# Redis type

This type manages Redis configs.

## Required variables

The following env variables should be defined when using this type:

* **REDIS_HOST** defines Redis host, it must be a string, its default value is "localhost" 
* **REDIS_PORT** defines Redis port, it must be a valid port number, its default value is 6379
* **REDIS_PASSWORD** defines Redis password, its default value is an empty string
* **REDIS_DATABASE** defines Redis database, its default value is 0 
