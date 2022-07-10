proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 | sort | uniq)

str_args=""
for dir in $proto_dirs; do
  dir=${dir:2}
  str_args="$str_args --go_opt=M$dir=."
done

str_files=""
for file in $proto_dirs; do
  str_files="$str_files ${file:2} "
done


protoc -I=. ${str_args} --go_out=. ${str_files}