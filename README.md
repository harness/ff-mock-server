# Feature Flag mock server

Very light mock server which serves simple configurations and target groups. Supports streams and metrics validation.

# How to run
```
docker run -d -p 9090:3000 ff-mock-server:latest
// wait until the server starts
docker exec ${ID} /bin/bash /app/wait-for-it.sh
```
# Authentication keys
You need to specify api key to access all api endpoints
* Server key: `2e182b14-9944-4bd4-9c9f-3e859e2a2954`
* Client key: `2e2ecf62-ce53-4e9e-8006-b4db0386688c`

# Default values
```
* environmentUUID: 265597ad-516c-4575-a16f-b3d17adffc44
* defaultClusterIdentifier: cluster
* project: demo
* environment: dev
* flag: bool-flag
* target group: demo
```

# How to use

When server is started sample app can connect and use custom config Url and event Url

```java
final HarnessConnector hc = new HarnessConnector(SDK_KEY, HarnessConfig.builder().configUrl("http://localhost:9090/api/1.0").build(), null);
final CfClient client = new CfClient(hc);
client.waitForInitialization();
final boolean bResult = client.boolVariation("bool-flag", null, false);

assert bResult == true
```