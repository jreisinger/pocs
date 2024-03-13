1BRC (The One Billion Rows Challenge) wants you to process a 13GB file of weather station names and temperatures as fast as possible. For each weather station, print out the minimum, mean, and maximum. The file looks like this:

```
$ head measurements.txt 
Kano;29.3
Dili;32.4
Kuala Lumpur;21.2
Dushanbe;12.5
Tegucigalpa;14.2
Irkutsk;21.6
Hargeisa;29.7
Detroit;13.6
Nicosia;-5.8
Chișinău;7.3
```

To generate the file clone [this repo](https://github.com/gunnarmorling/1brc), change to it, install JDK (`sudo apt install openjdk-21-jdk`) and run:

```
$ ./mvnw clean verify
$ ./create_measurements.sh 1000000000
```

To process the file:

```
$ go build
$ time ./1brc [-p] ~/github.com/gunnarmorling/1brc/measurements.txt
```

Based on https://benhoyt.com/writings/go-1brc/.
