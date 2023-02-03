function fn() {
  var API_URL = "http://localhost:8098/snapfi"

  var config = {
    apiURL: API_URL,

    uuid: uuid,
    sleep: sleep,
    funcNow: now,

  };

  karate.configure("connectTimeout", 15000);
  karate.configure("readTimeout", 15000);
  karate.configure('retry',{ count:10, interval:2000});

  return config;
}

function uuid() { return java.util.UUID.randomUUID() + '' }
function sleep(milliseconds) { java.lang.Thread.sleep(milliseconds) }
function now() {
  return (new Date()).toJSON();
}