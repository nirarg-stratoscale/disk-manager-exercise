# How to use auth with swagger

## 1. Update swagger.yaml

Add default security settings to the root of your swagger.yaml file:

```
# Security definitions defines the security of this server.
# We use roles based access, so each API can define the roles it can access.
securityDefinitions:
  roles:
    type: apiKey
    in: header
    name: X-Auth-Roles

# Security section defines the default security permissions for all the APIs.
# This can be override in each API with a security section
security:
  - roles:
    - admin
```

The first section defines that we have a header containing the role, and we want to enforce the
role in some way in our middleware.
The second section defines that the default route permissions are users that hold the `admin` role.
It will be applied to all routes in your server.

If a specific route needs a more specific permissions, for example, `admin` and `_member_` roles are
allowed, you can set the security options in this route:

```diff
paths:
  /pets:
    [...]
    get:
      tags:
        - pet
      summary: List pets
      operationId: PetList
      parameters:
        [...]
      responses:
        [...]
+      security:
+        - roles:
+          - admin
+          - _member_
```

We added a security section, and defined that the allowed roles are admin or _member_

## 2. Generate swagger code

Run the swagger command to generate the swagger code

## 3. Add middleware

In the main.go, where the go-swagger handler is initialized, add the following code:

```diff
import (
+	"github.com/Stratoscale/golib/auth"
)
[...]
func main() {
    [...]
	h, err := restapi.Handler(restapi.Config{
		PetAPI:          pet,
		StoreAPI:        store,
		[...]
		Logger:          a.Log.WithField("pkg", "restapi").Debugf,
+		AuthRoles:       func(token string) (interface{}, error) { return "", nil },
+		InnerMiddleware: auth.Middleware(a.Log.WithField("pkg", "auth")),
	})
	[...]
}
```

This will enforce access to the routes according to the role of the user and the defined
roles on the route. It will also add the user information to the request context, so it
could be used in the business logic.

## 4. Use the user information in your code

In business logic, the user information can be extracted as followes:

```go
import (
	"github.com/Stratoscale/golib/auth"
)
[...]
func (m *manager) PetList(ctx context.Context, params pet.PetListParams) middleware.Responder {
    user := auth.FromContext(ctx)
    // do something with user, if middleware was used, user won't be nil here, however, it's fields
    // can be empty if they were not in the headers.
}
```

## 5. Add tests

The policy enforcement behavior should be tested in `main_test.go`.
See [example](https://github.com/Stratoscale/go-template/blob/9bb00615a88a950aa47caeb29bdce7ec0d8b2274/example-pet-store/main_test.go#L40)

## 6. Use in go clients

You can inject user credentials in go clients using the roundtripper:

```go
user := &auth.User{
    [...]
}
transport := auth.RoundTripper(user, nil)
petStoreClient := petstoreclient.New(petstoreclient.Config{Transport: transport})
// use petStoreClient

httpClient := http.Client{Transport: transport}
// use httpClient
```
