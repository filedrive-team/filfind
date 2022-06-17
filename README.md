FilFind
==================
[![](https://img.shields.io/github/license/filedrive-team/filfind)](https://github.com/filedrive-team/filfind/blob/main/LICENSE)

## Project Description
FilFind, a storage provider discovery platform!  

Initially, automatically populated data will come from a number of sources:
* Deal success data (storage, retrieval) to come from the dealbot
* Reputation and other data to be sourced from filrep.io
* Deal data to be sourced from filecoin (api.node.glif.io)

## Dependency
Environment
```
go version 1.17
node.js version 14.17.3
npm version 6.14.13
```

## Config
#### Back-end config
Please set database config and smtp config in backend/conf/app.toml.
```
[app]
# Set to false if deployed to a production environment
debug = true
swag = true
filrepApi = "https://api.filrep.io/api"
filecoinApi = "https://api.node.glif.io/rpc/v0"
#filecoinApi = "https://api.chain.love/rpc/v1"
jwtSecret = "7b2274797065223a224853323536222c22707269766174655f6b6579223a226b4238424c6d6765512f4a34714a4c7a6635657562544c67777454594332356f763271372f766e58446e773d227d"
passwordSalt = "6a9741"
publishDate = "2022-06-16T00:00:00Z"
officialWebsite = "https://filfind.info"
officialEmail = "support@filfind.info"

[server]
httpPort = 9095
readTimeout = 60
writeTimeout = 60

[database]
type = "mysql"
dsn = "filfind:filfind@tcp(127.0.0.1)/filfind?charset=utf8&parseTime=True&loc=Local"

[smtp]
[smtp.basic]
host = "smtpout.secureserver.net:25"
user = ""
password = ""
```

#### Front-end config
Please set back-end api with your url in frontend/.env.xxx.

local server: .env.development
```
VUE_APP_BASE_URL=http://localhost:8088
```
product server: .env.production
```
VUE_APP_BASE_URL=YOUR_BACKEND_API_URL
```

## Deployment
### Build and run back-end
#### Build
```bash
cd backend
make
./filfind-backend
```
#### Initializing the system
It may take a few hours, depending on your network.
```bash
./filfind-backend -init
```

#### Run
```bash
./filfind-backend
```

### Build front-end
```bash
cd frontend
yarn install
yarn build
```

### Dist
```bash
cd frontend
copy all files of dist into the root of website
```


## Contribute
PRs are welcome!

## License
MIT

