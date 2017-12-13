#!/usr/bin/env bash
echo "mode: set" > acc.out
for Dir in $(find ./* -maxdepth 10 -type d | grep -v vendor);
do
        if ls $Dir/*.go &> /dev/null;
        then
            echo "Testing $Dir"
            go test -v -coverprofile=profile.out $Dir
            if [ -f profile.out ]
            then
                cat profile.out | grep -v "mode: set" >> acc.out
            fi
fi
done
goveralls -coverprofile=profile.out -service travis-ci -repotoken $COVERALLS_TOKEN
rm -rf ./profile.out
rm -rf ./acc.out

