mvn compile exec:java -Dexec.mainClass=org.apache.beam.examples.WordCount \
    -Dexec.args="--runner=SparkRunner --inputFile=sample.txt --output=counts" -Pspark-runner