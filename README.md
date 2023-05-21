# glb

This project is a simple load balancer written in Go. It can be used to distribute traffic between a set of backend servers. The load balancer uses a health check to determine which servers are available and then forwards traffic to the available servers.

The load balancer is implemented using the Go standard library's http.

The load balancer can be configured with a list of backend servers. The list of backend servers can be specified in a configuration file. The load balancer will check the health of the backend servers before forwarding traffic to them. If a backend server is not healthy, the load balancer will not forward traffic to it.

The load balancer can be used to distribute traffic between a set of backend servers. This can help to improve the performance and availability of your applications.

# Getting started

Here are the steps on how to use the load balancer:

1. Install Go.
1. Clone the repository.
1. cd into the directory.
1. Run go build.
1. Write config file at path ~/.glb/config.yaml as below:
1. Run ./glb

```yaml
- path: "/a"
  algorithm: firstActive
  hosts:
      - protocol: http
        hostname: localhost
        port: 8080
        health:
            endpoint: "/health"
            successCode: 200
            method: GET
        minHealthyHits: 2
        minUnhealthyHits: 3
        hitFrequencyInSeconds: 5
      - protocol: http
        hostname: localhost
        port: 8081
        health:
            endpoint: "/health"
            successCode: 200
            method: GET
        minHealthyHits: 2
        minUnhealthyHits: 3
        hitFrequencyInSeconds: 5
```

## Explanation

### `path` (required)

The `path` property specifies the URL path for which the rule applies. Requests with matching paths will be directed to the configured hosts. This property is a string value.

Example:`path: "/a"`

### `algorithm` (optional)

The `algorithm` property determines the load balancing algorithm to be used for distributing traffic among the available hosts. This property is optional, and if not specified, a default algorithm may be used.

Example: `algorithm: firstActive`

Available Algorithms

-   `firstActive`: Selects the first available host based on health checks.

Note: Additional algorithms may be available depending on the load balancer implementation.

### `hosts` (required)

The `hosts` property specifies the list of hosts to which the traffic will be distributed. Each host is defined by a set of properties.

Example:

```yaml
hosts:
    - protocol: http
      hostname: localhost
      port: 8080
      health:
          endpoint: "/health"
          successCode: 200
          method: GET
      minHealthyHits: 2
      minUnhealthyHits: 3
      hitFrequencyInSeconds: 5
    - protocol: http
      hostname: localhost
      port: 8081
      health:
          endpoint: "/health"
          successCode: 200
          method: GET
      minHealthyHits: 2
      minUnhealthyHits: 3
      hitFrequencyInSeconds: 5
```

### Host Properties

#### `protocol` (required)

The `protocol` property specifies the protocol to be used for communication with the host. This property is a string value and can be set to `http` or `https`.

Example: `protocol: http`

#### `hostname` (required)

The `hostname` property specifies the hostname or IP address of the host. This property is a string value.

Example: `hostname: localhost`

#### `port` (required)

The `port` property specifies the port number on which the host is listening for incoming requests. This property is an integer value.

Example: `port: 8080`

#### `health` (required)

The `health` property defines the health check configuration for the host. It determines how the load balancer checks the health of the host.

-   `endpoint` (required): The endpoint property specifies the URL path or endpoint on the host to be used for health checks. This property is a string value.

    Example: `endpoint: "/health"`

-   `successCode` (required): The successCode property specifies the HTTP status code indicating a successful health check response. This property is an integer value.

    Example: `successCode: 200`

-   `method` (required): The method property specifies the HTTP method to be used for health checks. This property is a string value and can be set to GET, POST, PUT, DELETE, or any other valid HTTP method.

    Example: `method: GET`

-   `minHealthyHits` (required): The minHealthyHits property specifies the minimum number of consecutive successful health check responses required for a host to be considered healthy. This property is an integer value.

    Example: `minHealthyHits: 2`

-   `minUnhealthyHits` (required): The minUnhealthyHits property specifies the minimum number of consecutive failed health check responses required for a host to be considered unhealthy. This property is an integer value.

    Example: `minUnhealthyHits: 3`

-   `hitFrequencyInSeconds` (required): The hitFrequencyInSeconds property specifies the frequency at which health checks should be performed. This property is an integer value representing the time interval in seconds.

    Example: `hitFrequencyInSeconds: 5`

## Health check
This load balancer runs at port `8000`. If you need to see what hosts are healthy and what hosts are unhealthy, you can perform a `GET` request to the endpoint `/glb/health` as below:

### Sample request
```shell
curl 'localhost:8000/glb/health'
```

### Sample response
```json
{
    "registry": {
        "52bc80a6-e915-41e5-acb4-455abe83e474": {
            "hostConfig": {
                "uniqueId": "52bc80a6-e915-41e5-acb4-455abe83e474",
                "protocol": "http",
                "hostname": "localhost",
                "port": 8080,
                "health": {
                    "Endpoint": "/health",
                    "SuccessCode": 200,
                    "Method": "GET"
                },
                "minHealthyHits": 2,
                "minUnhealthyHits": 3,
                "hitFrequencyInSeconds": 5
            },
            "lastChecked": "2023-05-21T15:33:48.292662658+05:30",
            "isHealthy": false,
            "healthyHitCount": 0,
            "unhealthyHitCount": 1,
            "lastHitAt": "0001-01-01T00:00:00Z"
        },
        "edbf3cbe-06ba-4105-a0d9-424fbbe27694": {
            "hostConfig": {
                "uniqueId": "edbf3cbe-06ba-4105-a0d9-424fbbe27694",
                "protocol": "http",
                "hostname": "localhost",
                "port": 8081,
                "health": {
                    "Endpoint": "/health",
                    "SuccessCode": 200,
                    "Method": "GET"
                },
                "minHealthyHits": 2,
                "minUnhealthyHits": 3,
                "hitFrequencyInSeconds": 5
            },
            "lastChecked": "2023-05-21T15:33:48.293621018+05:30",
            "isHealthy": false,
            "healthyHitCount": 0,
            "unhealthyHitCount": 1,
            "lastHitAt": "0001-01-01T00:00:00Z"
        }
    }
}
```

## Limitations

This load balancer is a very basic implementation and may lack many features like more algorithm implementations. But it can be developed over time.

