# aws-secrets-reader

This application is created to be used with [Launch](https://github.com/morfien101/launch).

Currently it is able to collect secrets from AWS Secrets Manager and then do some post processing.

# Help Menu

The help menu should get you going for most of the questions of how to use this. `-h`

```sh
Collects secrets from AWS Secrets manager.
version: development
  -aws-profile string
        AWS Profile to use. Blank by default and omitted.
  -h    Help menu.
  -prepend-with string
        Prepend the returned keys with given string. Upper casing happens after this is applied.
  -region string
        AWS region to use. eu-west-1 by default. (default "eu-west-1")
  -secret string
        The key to use when collecting the secret.
  -upper-case
        Attempt to uppercase all the returned keys
  -v    Shows the version.
```


# Post Processing

## Prepend to keys

Prepend a string at the beginning of all collected secret keys. This is useful if you are looking to collect many keys and they have similar names.
Or if you are looking to add a prefix to the keys to have them automatically ingested.

Example

```sh
# -prepend-with potatoes_
# this
{"badger":"mushroom"}
# becomes
{"potatoes_badger":"mushroom"}

# remember to add separators like _ or - as they are not automagically added.
```

## UPPER CASE keys

Most environment variables are in upper case. However they are not stored in upper case in the secret manager.
Therefore we need a way to quickly uppercase them.

`Upper case happens AFTER the prepend action or any other post process`

Example

```sh
# -upper-case
# this
{"badger":"mushroom"}
# becomes
{"BADGER":"mushroom"}
```

## Upper and prepend together 😱

You can use both the post processors currently available. Just remember that UPPER CASE is always last.

```sh
# -prepend-with potatoes_
# -upper-case
# this
{"badger":"mushroom"}
# becomes
{"POTATOES_BADGER":"mushroom"}
```