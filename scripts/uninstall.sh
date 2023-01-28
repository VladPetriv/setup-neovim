#!/bin/sh

rm -rf ~/go/bin/setup-nvim

export PATH=${PATH%:/home/go/bin/setup-nvim}

echo "setup-nvim successfully uninstalled"
