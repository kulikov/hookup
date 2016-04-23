# HookUp

[![Circle CI](https://circleci.com/gh/kulikov/hookup.svg?style=svg&circle-token=6f95ac3d3af195cfb6dc9ff09569b7828cb17133)](https://circleci.com/gh/kulikov/hookup)

Simple http service, listen web hooks form:

 * github

And call all handlers scripts, found at `--handlers` path


```
NAME:
   hookup - Start Webhook Server

USAGE:
   hookup [global options] command [command options] [arguments...]

VERSION:
   1.0.0

GLOBAL OPTIONS:
   --port "9090"
   --handlers "/etc/hookup.d/"	Path to dir with webhook handlers scripts
   --help, -h					show help
   --version, -v				print the version
```
