### How to retrieve the JWT token from the Metadata API endpoint

Salad environment provide a metadata Api which offers various services for the containers running on the Salad node. The Api is available on the IP address `169.254.169.254`, port `80`

The `/v1/token` endpoint provides the way for the container to request a JWT token to be used across the Salad Cloud. This endpoint is a simple GET endpoint responding with the json containing the JWT with the following schema:

```
{
  'jwt": "eyJhbGci..."
}
```

Below is an example code how to retrieve the token with Python program using the [requests](https://docs.python-requests.org/en/latest/index.html) library

```
import requests

r = requests.get('http://169.254.169.254:80/v1/token')
token =r.json()['jwt']
print(token)
```
