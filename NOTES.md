08/16/2024
- If `-t` or `-x` is given for `dir` command, you can continue with only one of them, but should execute all if anyelse is also given
- If `-t`, `-x`, `-s` is given for `file` command, you can continue with only one of them, but should execute all if anyelse is also given
- When using `-x` flag with dir command, I wanna be able to send execute the given command after a file is, you know, encrypted/decrypted. It shouldn't block my encryption stuff, it should spawn a new process or something?
- Ok, so for every file that is decrypted/encrypted by `dir` command, it will be passed to the command provided by the `-x` flag, and the cli will wait till the command execution is over. also, errors in the command execution will be logged by the cli, but won't stop the cli from running