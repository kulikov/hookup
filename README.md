# WebHooker

Simple http service, listen web hooks form:

 * github

And call all handlers scripts, found at `--handlers` path


```
NAME:
   webhooker - Start Webhook Server

USAGE:
   webhooker [global options] command [command options] [arguments...]

VERSION:
   1.0.0

GLOBAL OPTIONS:
   --port "9090"
   --handlers "/etc/webhooker.d/"	Path to dir with webhook handlers scripts
   --help, -h				show help
   --version, -v			print the version
```
