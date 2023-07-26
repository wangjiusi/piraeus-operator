#!/bin/bash -x

# Get chart info
index=$(curl -sS https://raw.githubusercontent.com/piraeusdatastore/helm-charts/gh-pages/index.yaml )
[ -z $index ] && exit 1

# Download chart
url=$( echo "$index" | yq e .entries.piraeus-ha-controller[0].urls[0] )
cd ./charts/piraeus/charts/ || exit 1
rm -vfr ./ha-controller
curl -sSL "$url" | tar -zxvf -
mv -v piraeus-ha-controller ha-controller

# Update name and version
yq -i '.name = "ha-controller"' ha-controller/Chart.yaml
export version=$( echo "$index" | yq e .entries.piraeus-ha-controller[0].version )
[ -z "$version" ] && exit 1
for i in ../Chart.*; do
    yq e '(.dependencies[] | select(.name == "ha-controller").version) = strenv(version)' -i $i
done
yq -i '.ha-controller.image.tag = "v" + strenv(version)' ../values.yaml
