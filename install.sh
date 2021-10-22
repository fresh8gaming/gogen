#!/bin/bash
case $(uname -sm) in
  'Linux x86_64')
    os='linux_amd64'
    family='linux'
    ;;
  'Darwin x86' | 'Darwin x86_64')
    os='darwin_amd64'
    family='mac'
    ;;
  'Darwin arm64')
    os='darwin_arm64'
    family='mac'
    ;;
  *)
  echo "Sorry, you'll need to install the gogen manually."
  exit 1
    ;;
esac

if [[ -z "${TAG}" ]]; then
  tag=$(basename $(curl -fs -o/dev/null -w %{redirect_url} https://github.com/fresh8gaming/gogen/releases/latest))
else
  tag=${TAG}
fi

filename="gogen_${tag#v}_${os}.tar.gz"

curl -LO https://github.com/fresh8gaming/gogen/releases/download/${tag}/${filename}
tar xzf ${filename}
rm ${filename}

case $family in
  'linux')
    mv ./gogen ~/.local/bin
    chmod +x ~/.local/bin
    ;;
  'mac')
    sudo mv ./gogen /usr/local/bin
    sudo chmod +x /usr/local/bin
    ;;
  *)
  echo "Sorry, you'll need to move the gogen binary manually."
  exit 1
    ;;
esac