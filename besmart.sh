#! /bin/bash

# An easy install script for Golang and the program

# Check if GO is installed
if [ ! -d /usr/local/go  ]; then
  echo "Error: Go is not installed."
  dependancyFailure=true
  read -p "Automatically install latest version? [y/n]: " response
  if [[ "$response" =~ ^([yY][eE][sS]|[yY])+$ ]]
  then
    wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
    rm -rf go1.21.6.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    dependancyFailure=false
    echo "Go successfully installed! - Path has been altered temporarily, please add \"/user/local/go/bin\" to ~/.profile"
  else
    echo "Skipping..."
  fi
fi
