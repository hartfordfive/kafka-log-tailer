
#!/bin/sh

for build in build/*{darwin,linux}*
do
  tar -czvf "$build.tar.gz" $build && rm -f $build
done

for build in build/*windows*
do
  zip "$build.zip" $build && rm -f $build
done