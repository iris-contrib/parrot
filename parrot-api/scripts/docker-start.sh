#!/bin/bash

script_dir="${0%/*}"
cd $script_dir/..

export PARROT_API_ROOT=$PWD
cd $PARROT_API_ROOT/..

if [ "$PARROT_ENV" == "production" ]
then
    echo "Parrot API environment set to 'production'"
    echo "Fetching latest version of parrot..."
    rm -rf parrot && \
        mkdir parrot && \
        cd parrot && \
        git init && \
        git remote add origin https://github.com/iris-contrib/parrot.git && \
        git config core.sparseCheckout true && \
        echo "parrot-api" >> .git/info/sparse-checkout && \
        git pull --depth=1 origin master
else
    echo "Parrot API environment set to 'development'"
    echo "Using local repo..."
fi

cd $PARROT_API_ROOT && \
    chmod +rwx ./scripts/server-start.sh && \
    /bin/bash ./scripts/server-start.sh