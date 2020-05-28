## fic firewalls activate

Activate firewall

### Synopsis

Activate firewall

```
fic firewalls activate <id> [flags]
```

### Examples

```
fic firewalls activate F040123456789 --router F022000000335 --user-ip-addresses 192.168.0.0/30,192.168.0.4/30,192.168.0.8/30,192.168.0.12/30
```

### Options

```
  -h, --help                        help for activate
      --router string               (required) Router ID
      --user-ip-addresses strings   Local IP addresses
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.fic.yaml)
      --debug           debug mode
```

### SEE ALSO

* [fic firewalls](fic_firewalls.md)	 - Firewall management

###### Auto generated by spf13/cobra on 28-May-2020