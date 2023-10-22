#CLI Setup
The CLI is used to execute jobs and to tell the agent what to do.

## Install
### Install from Releases
Go to [Releases](https://github.com/XiovV/dokkup/releases) and download the binary for your OS and CPU architecture.

### Install from source:
Clone the repository:
```shell
$ git clone https://github.com/XiovV/dokkup.git
```

Install the binary:
```shell
cd dokkup/cmd/dokkup
go install
```

Run `dokkup` in your terminal to verify it's installed:
```shell
$ dokkup version

Dokkup v0.1.0-beta, build fd13a57
```
