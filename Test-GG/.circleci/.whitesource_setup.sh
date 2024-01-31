curl -L https://go.dev/dl/${GO_VERSION}.linux-amd64.tar.gz --output go.tar.gz
tar -C /usr/local -xzf go.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> $BASH_ENV
echo "export WS_PROJECTNAME=$CIRCLE_PROJECT_REPONAME" >> $BASH_ENV
