<h1 align="center">
Ordfanger Telegram Bot
</h1>

This bot was created to collect words. You can deploy it in your own AWS lambda and record words for learning purposes. 

## Install

Required: go, make, AWS account.

```
$ npm i serverless -g
$ make deploy
```

*Note that in the serverless config defined AWS profile name SERVERLESS_USER. You need to create AWS profile with that name or define your own*


## Config
All needed configuration described in serverless.yml.


## Tests
```
$ make test
```

## License

MIT
