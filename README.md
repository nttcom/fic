# fic

Command Line Interface for Flexible InterConnect

## Installation

```Bash
go get -u github.com/nttcom/fic
```

Make sure that `PATH` includes `$GOPATH/bin`.

```Bash
export PATH=$PATH:$GOPATH/bin
```

## Configuration

There are two ways to set your configurations.

* Export environment variables:

```Bash
export FIC_USERNAME=****
export FIC_PASSWORD=****
export FIC_TENANT_ID=****
```

* Create config file:

```Bash
touch ~/.fic.yaml
echo "username: ****" >> ~/.fic.yaml
echo "password: ****" >> ~/.fic.yaml
echo "tenant_id: ****" >> ~/.fic.yaml
```

## Execution

```Bash
fic areas list
```

## Command reference

See [fic.md](https://github.com/nttcom/fic/blob/master/doc/fic.md).

## Auto completion

### Bash

The completion script for Bash can be generated with the command `fic completion bash`.
It depends on [bash-completion](https://github.com/scop/bash-completion), which means that you have to install this software first.
If you use the completion script on macOS, it requires **bash-completion v2** and **Bash 4.1+**.
[Kubectl docs](https://kubernetes.io/docs/tasks/tools/install-kubectl/#enabling-shell-autocompletion) show how to install bash-completion and bash in details.  

There are two ways to enable auto completion.

* Source the completion script in your `~/.bashrc` file:

```Bash
echo 'source <(fic completion bash)' >>~/.bashrc
``` 

* Add the completion script to the `bash_completion.d` directory:

```Bash
# On Linux
fic completion bash >/etc/bash_completion.d/fic
# On macOS
fic completion bash >/usr/local/etc/bash_completion.d/fic
```

### Zsh

The completion script for Zsh can be generated with the command `fic completion zsh`.

You need to put the generated script somewhere in your *fpath* named *_fic* like this.

```Zsh
fic completion zsh > ${fpath[1]}/_fic
```

## License

Fic is released under the Apache 2.0 license. See [LICENSE](https://github.com/nttcom/fic/blob/master/LICENSE). 