#!/bin/sh

rm -rf ./setup-nvim

make build 

mv ./setup-nvim ~/go/bin/setup-nvim

export PATH="$PATH:~/go/bin/setup-nvim"

echo "setup-nvim successfully installed"