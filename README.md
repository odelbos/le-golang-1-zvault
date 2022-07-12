# Disclaimer

This repository is an Golang learning exercise.

**Do not use this code in production.**

_(with all my repositories, the `le-` prefix mean `Learning Exercise`)_

# Sysnopsis

The goal of this exercise is to write a simple binary used to encrypt/decrypt files in a vault.

It will have 3 commands :

- init _(used to init a vault configuration)_
- put _(used to put a file in the vault)_
- get _(used to get a file given his id)_

Encryption will be made with `AES_256_GCM` and password derivation will be made with `PBKDF2`.

## Init command

The `init` command is used to initialize a new vault configuration.

By default the configuration is saved in [USER_HOME]/.config/zvault.json

```
zvault init
```

If you want to specify a particulary configuration file :

```
zvault -c /path/to/conf/file.json init
```

During the `init` process you will be prompted for :

- The folder where to store the encrypted blocks,
- The folder where to store the encrypted file description,
- A master key password.

## Put command

The `put` command is used to store a file in the encrypted vault.

```
zvault put /absolute/path/to/myfile.txt
```

You will get back the id of the file, like : `a1126d9fc7c2fc240d6c44e267ed2097`

## Get command

The `get` command is used to get back a stoed file.

```
zvault get a1126d9fc7c2fc240d6c44e267ed2097
```

The file will be restored in the current directory.