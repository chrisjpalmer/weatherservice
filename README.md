To get started do the following:
```bash
go get github.com/chrisjpalmer/weatherservice
weatherservice
curl "http://localhost:8080/v1/weather?city=sydney"
```

### Tradeoffs / Nice to haves:
* You could generalize the idea of a "weather provider" by making a generic weather provider interface. Then dynamically add these to a the parent `WeatherService` class at bootstrap
* If many developers were to work on this, there would need to be a better abstraction for adding controllers and services. You don't want lots of developers touching the bootstrap of the application - there probably should be some setup function you call which registers a controller or service. Dont know if dependancy injection is possible in go? -- Am I thinking too much like a front-end developer?
* If you were to scale this up, you may use redis to cache the data so every process of the web server would not have its own copy of the data. 
You may also add another microservice whose job is to periodically retrieve data from the providers and store in redis. Then the http service would be a dumb abstraction from the redis cache.