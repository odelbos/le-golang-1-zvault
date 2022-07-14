# Disclaimer

This repository is an Golang learning exercise.

**Do not use this code in production.**

_(with all my repositories, the `le-` prefix mean `Learning Exercise`)_

# Synopsis

The goal of this exercise is to write a simple binary used to store/restore files in an encrypted vault.

Encryption will be made with `AES_256_GCM` and password derivation will be made with `PBKDF2`.

It will have 3 commands :

- init _(used to init a vault configuration)_
- put _(used to put a file in the vault)_
- get _(used to get a file given his id)_

## Build the `zvault` binary

```sh
./build.sh
```

## Init command

The `init` command is used to initialize a new vault configuration.

By default the configuration is saved in [USER_HOME]/.config/zvault.json

```sh
./zvault init
```

If you want to use a specific configuration file :

```sh
./zvault -c /path/to/conf/file.json init
```

During the `init` process you will be prompted for :

- The folder where to store the encrypted blocks,
- The folder where to store the encrypted file description,
- A master key password.

## Put command

The `put` command is used to store a file in the encrypted vault.

```sh
./zvault put /absolute/path/to/myfile.txt
```

You will get back the id of the file, like : `a1126d9fc7c2fc240d6c44e267ed2097`

## Get command

The `get` command is used to get back a stored file.

```sh
./zvault get a1126d9fc7c2fc240d6c44e267ed2097
```

The file will be restored in the current directory.

# Example of usage

Create a storage folder structure :

```sh
mkdir storage
mkdir storage/data
mkdir storage/files
```

Create a random file :

```sh
mkdir data
dd if=/dev/random of=./data/file-9mb.bin bs=1 count=9545925
```

Initialize the vault :

```sh
./zvault init
> Data path : storage/data
> Files path : storage/files
> Enter paswsord: *******
> Repeat paswsord: *******
```

Store a file in the vault :

```sh
./zvault put data/file-9mb.bin
> Enter Password: *******
File stored, id: 9deba552fe5c0b04b4e5dbc84cb65324
```

Restore a file from the vault :

```
% ./zvault get 9deba552fe5c0b04b4e5dbc84cb65324
> Enter Password: *******
File restored, name: file-9mb.bin
```

Verify that files are the same :

```
% md5 data/file-9mb.bin
MD5 (data/file-9mb.bin) = 0552c4b808193553cfed8bf562a41d8c

% md5 file-9mb.bin
MD5 (file-9mb.bin) = 0552c4b808193553cfed8bf562a41d8c
```

### If using a specific configuration file

```sh
./zvault -c /path/to/config.json init
./zvault -c /path/to/config.json put data/file-9mb.bin
./zvault -c /path/to/config.json get 9deba552fe5c0b04b4e5dbc84cb65324
```

# Roadmap

- [ ] Better errors handling
- [ ] Use CLI package to manage commands (see: https://github.com/urfave/cli)

# Author

Author : @odelbos