# keycloak-spike

This spike creates a keycloak server, creates a client and user, and hooks into it with a simple Go application.

## setup

```
docker run -d -p 8080:8080 --net host -e KEYCLOAK_USER=admin -e KEYCLOAK_PASSWORD=admin --name keycloak jboss/keycloak:latest
```

## configure

1. create a realm
2. turn off ssl requirement for `demo` realm
3. create a client
4. create a user and password
5. configure client redirect URI to `http://localhost:9000/auth/keycloak/callback`
6. configure client access type to be `confidential`
7. pull keycloak.json from installation tab
8. fill in the required fields:

| go code | keycloak.json |
|---------|---------------|
| ClientID | resource |
| ClientSecret | credentials.secret |

