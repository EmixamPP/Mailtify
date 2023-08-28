# <h1 align="center">mailtify</h1>
A simple server to easily send notifications via a RESTful API and receive them as mails.

# Quick utilization
### Create a token:
`curl -u username:password mailtity.domain.com/new`

-> `{"status": 200 "message":"<token>"}`

### Send a message:
`curl "mailtify.domain.com/msg?token=<token>" -F "title=my title" -F "message=my message" -F "recipients=my@domain.com"`

### Delete a token:
`curl -u username:password -X DELETE "mailtify.domain.com/delete?token=<token>"`

### Get all tokens:
`curl -u username:password "mailtify.domain.com/tokens"`

# Setup
A `config.yml` must be present in the root directory, see [config-example.yml](config-example.yml) for the structure and required fields.

For ssl, and http to https redirection, you should use a reverse proxy (e.g. NGINX Reverse Proxy, Traefik).

# Which features come next?
1. Docker container
2. API doc
3. Retry on network error
