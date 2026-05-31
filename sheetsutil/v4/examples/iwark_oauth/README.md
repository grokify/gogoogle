# Using Iwark Spreadsheet

## Instantiation

The [Iwark Spreadsheet `README.md`](https://github.com/Iwark/spreadsheet/blob/41eea14839643a5a737559747bb7318c2eaad600/README.md) provides the following for instantiation. `client`. is an `*http.Client`.`

```go
data, err := ioutil.ReadFile("client_secret.json")
checkError(err)

conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
checkError(err)

client := conf.Client(context.TODO())
service := spreadsheet.NewServiceWithClient(client)
```

## Instantiation using GoAuth.

The same can be accomplished using `goauth`'s `CredentialsGCP` as follows.

```go
creds := goauth.CredentialsGCPReadFile("client_secret.json")

client := creds.NewClient(context.Background()) // returns `*http.Client`
service := spreadsheet.NewServiceWithClient(client)
```
