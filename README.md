
# sway-icon-to-go

- Renames sway workspaces according to the workspace window names
- Allows usage of [Font Awesome](https://origin.fontawesome.com/icons?d=gallery) icons instead of app names

## Setup
1. To support icons Font Awesome should be available on your system. In case it is not installed - use your favorite package manager to install it.
You can use `fc-list | grep Awesome` or just `sway-icon-to-go awesome` to check the Font Awesome availability on your system

2. **configs** directory contains sample configuration files in yaml format
These files should be placed either under `~/.config/sway` or `~/.config/i3` directory.
In case both locations contain config files `~/.config/sway` takes precedence.

`fa-icons.yaml` sets one-to-one mapping from icon name to UTF-8 code as set by Font Awesome.
`app-icons.yaml` sets one-to-many mapping from icon name to app name (lowercase)
A default `fa-icons.yaml` can be produced by executing `sway-icon-to-go parse > ~/.config/sway/fa-icons.yaml`

3. Just place the executable file anywhere and add this line to your sway config:
`exec sway-icon-to-go`

   **Alternatively**, run as a user-level systemd service:
   ```
   make install-service   # install, enable and start
   make reload-service    # rebuild and restart
   make uninstall-service # remove
   ```

4. Hot reload icons file without restarting the application:
`pkill --signal HUP sway-icon-to-go` 

## Command line parameters
```
  -c         path to the app-icons.yaml config file
  -u         display only unique icons. Default is True
  -l         trim app names to this length. Default is 12
  -d         app delimiter. Default is a pipe character "|"
  -v         show verbose (debug) output. Default is False
```
Sample usage: `./sway-icon-to-go -u=true -d='+'`


Inspired by https://github.com/cboddy/i3-workspace-names-daemon
