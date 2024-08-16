# file-boy

**file-boy** is a no-BS encryption/decryption CLI, made with golang. It allows you to encrypt and decrypt files or entire directories with ease.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
  - [file](#file-command)
  - [dir](#dir-command)
- [Examples](#examples)
- [Flags](#flags)

## Installation

To install file-boy, clone the repository and build it using Go(>=1.20):

```bash
git clone https://github.com/Muhammed-Rajab/file-boy
cd file-boy
go build
```

## Usage

### `file` Command

Encrypts🔒 or decrypts🔓 a specified file using the given passphrase.

```sh
file-boy file [flags]
```

#### Local Flags

- `-f, --from`: The path to the file to encrypt/decrypt from (required).
- `-t, --to`: The path to the directory to encrypt/decrypt to.
- `-m, --mode`: The mode of operation (`encrypt`, `e`, `E`, `decrypt`, `d`, `D`). Default is `encrypt`.
- `-s, --stdout`: Writes the encrypted/decrypted data to stdout.
- `-v, --verbose`: Show detailed output.

#### Description

The file command allows you to encrypt or decrypt a specific file. You can specify the file to operate on with the `--from` flag and optionally provide a destination file with the `--to` flag. If no destination is provided, the output can be directed to stdout with the `--stdout` flag.

### `dir` Command

Encrypts🔒 or decrypts🔓 the specified directory into or from a `.zip` file.

```sh
file-boy dir [flags]
```

#### Local Flags

- `-f, --from`: The path to the directory to encrypt/decrypt from (required).
- `-t, --to`: The path to the directory to encrypt/decrypt to (required).
- `-m, --mode`: The mode of operation (`encrypt`, `e`, `E`, `decrypt`, `d`, `D`). Default is `encrypt`.
- `-v, --verbose`: Show detailed output.
- `-x, --exec`: Executes the given command and sets the `stdin` as the `.zip` bytes. You can use `{1}` in the argument and it will be replaced with the file path.

#### Description

The dir command allows you to encrypt an entire directory into a `.zip` file or decrypt a `.zip` file back into a directory. The `--from` flag specifies the source directory or zip file, while the `--to` flag specifies the destination.

## Examples

#### Encrypt a File

```sh
file-boy file -f secrets.txt -t encrypted/ -m e -v
```

This command encrypts the file `secrets.txt` and saves the encrypted file to `encrypted/` directory. The `-v` flag enables verbose output.

#### Decrypt a File and Output to Stdout

```sh
file-boy file -f secrets.encrypt -m d -s
```

This command decrypts the file `secrets.encrypt` and writes the output to stdout.

#### Encrypt a Directory

```sh
file-boy dir -f /path/to/dir -t /path/to/output -m e -v
```

This command encrypts the directory located at `/path/to/dir` and outputs the encrypted zip file to `/path/to/output`.

#### Decrypt a Directory

```sh
file-boy dir -f /path/to/encrypted/files/directory -t /path/to/store/decrypted.zip -m d -v
```

This command decrypts the `/path/to/encrypted/files/directory` directory and restores the directory to `/path/to/store/decrypted.zip`

#### Decrypt a Directory and pass each decrypted file to another program
```sh
file-boy dir -f /path/to/encrypted/files/directory -m d -v -x "echo {1}"
```

This command will decrypt every encrypted file in the directory and executes the given `-x` flag command, here which is `echo`. `{1}` is replaced by the path of the file which is decrypted. The `stdin` in of `echo` will be set to `zipBuf`, which is the internal representation of the zip file buffer in the cli.
## Flags

### Global Flags

- `-v, --verbose`: Show detailed output.

### `file` Command Flags

- `-f, --from`: The path to the file to encrypt/decrypt from (required).
- `-t, --to`: The path to the directory to encrypt/decrypt to.
- `-m, --mode`: The mode of operation (`encrypt`, `e`, `E`, `decrypt`, `d`, `D`). Default is `encrypt`.
- `-s, --stdout`: Writes the encrypted/decrypted data to stdout.

### `dir` Command Flags

- `-f, --from`: The path to the directory to encrypt/decrypt from (required).
- `-t, --to`: The path to the directory to encrypt/decrypt to (required).
- `-m, --mode`: The mode of operation (`encrypt`, `e`, `E`, `decrypt`, `d`, `D`). Default is `encrypt`.
- `-x, --exec`: Executes the given command and sets the `stdin` as the `.zip` bytes. You can use `{1}` in the argument and it will be replaced with the file path.

## TODO 📝
13 August 2024
- <s>[✅] output encrypted `.zip` file from encrypting a directory to stdout [16 August 2024]</s>
- <s>[✅] output decrypted `.zip` file from decrypting an encrypted directory to stdout [16 August 2024]</s>
- <s>[✅] ability to pass the encrypted directory files with metadata to other programs [16 August 2024]</s>
- <s>[✅] ability to pass the decrypted directory files with metadata to other programs [16 August 2024]</s>

## Contributing

Contributions are always welcome!

If you have a feature request, please open a new pull request with the regarding details. I'll be more than happy to connect with like-minded people😃!
