rm intermediate_*
rm mr-out*
for i in {1..20}
do
   echo "Execution $i"
   go run -race -gcflags="all=-N -l" mrworker.go wc.so
done
