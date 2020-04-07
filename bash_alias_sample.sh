cd() {
  ddd="$(pwd)/$1"
  if [ -d "$ddd" ]; then
    ntree send WORKDIR "$ddd";
  fi
  builtin cd $1;
}
