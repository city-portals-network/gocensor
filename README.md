# gocensor [![goreport](https://goreportcard.com/badge/github.com/city-portals-network/gocensor)](https://goreportcard.com/report/github.com/city-portals-network/gocensor) 

#Цензор матерных слов
- [x] add parse config
- [] add health response
- [x] add POST ro add dict
- [] add PKGBUILD
- [] add systemd unit
- [] tests
- [] validation input params

/v1/censor/append POST {word:data}
/v1/censor/delete POST {word:data}
/v1/censor/reload GET
/v1/censor/check GET
/v1/status GET


https://github.com/vearutop/php-obscene-censor-rus