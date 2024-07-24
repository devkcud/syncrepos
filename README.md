# syncrepos

A simple wrapper for `gh repo sync` to sync multiple repositories simultaneously, since the `gh` command cannot sync more than one repository at once.

## build

To build:

```bash
make build
```

To install:

```bash
make install
```

To do both:

```bash
make
```

The default behaviour of `make` is to build and install (copy into `~/.local/bin` as `sr`)

## usage

```bash
sr devkcud/syncrepos ...
# will add the devkcud/syncrepos repo to the repo list in $HOME/.config/.syncrepos
```

```bash
sr -s
# will sync the repos found in the $HOME/.config/.syncrepos
# To force the sync command, just place the -f flag
#
# Like this: sr -sf or sf -s -f
#
# This won't force anything in the syncrepos script, it's just a way to pass in the --force to gh repo sync command
```

You can modify where the repos file is located using:

```bash
sr -r "/home/pato/.config/.syncrepos"
```

> If you want to create the file if it doesn't exist use `sr -c -r "/home/pato/.config/.syncrepos"`

You can just place all the args together:

```bash
sr -cr "/home/pato/.config/.syncrepos" devkcud/syncrepos -sf
```

For more info view:

```bash
sr -h
```

> You can disable ALL output using the `-q` flag, keep in mind that even errors won't show up
