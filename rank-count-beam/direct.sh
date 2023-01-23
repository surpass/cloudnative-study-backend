mvn compile exec:java -Dexec.mainClass=cn.easyolap.bigdata.WordCount \
    -Dexec.args="--inputFile=sample.txt --output=counts" -Pdirect-runner