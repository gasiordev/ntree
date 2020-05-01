# Below code can be put in .bashrc or .bash_profile
# When directory is changed using 'cd', it's going to be send to running ntree

cd() {
  ddd="$(pwd)/$1"
  if [ -d "$ddd" ]; then
    ntree send WORKDIR "$ddd" > /dev/null 2>&1;
  fi
  builtin cd $1;
}
