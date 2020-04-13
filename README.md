# HMAC-wrapper
Golang HMAC-wrapper to authenticate with KONG HMAC

The latest stable version is 1.0.0, released on April 13, 2020. Latest source code is available from master branch on GitHub.

# How to use the wrapper

1. Import the package in your source code:
```
import (

  Wrapper "github.com/vinando/HMAC-wrapper"

)```

2. Init the class

The Init function need 3 parameters to pass into: client_id string, client_secret string and kong_base_url string.

```wrapper := Wrapper.Init(client_id, client_secret, kong_base_url)```

3. Doing request:

There are 2 method available:

  a. DoGet(endpoint string) (resp interface{}, err error)
  
  b. DoPost(endpoint string, body []byte) (resp interface{}, err error)
