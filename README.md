# 5km Run Scrapper

This is a small program for downloading all the publicly available data from the Bulgarian's 5kmrun website: http://5km.5kmrun.bg/

The output would be file `5km-stats.csv` with the following structure:

```
ID,name,is_male,age,place,date,time,position,avg_speed_kph,tempo
4,Георги Станоилов (Junior),1,4,Южен Парк,2013-04-13,27m34s,78,10.88,5m30s
7,Любка  Георгиева -Пелева,0,41,Южен Парк,2013-10-26,36m14s,119,8.28,7m14s
7,Любка  Георгиева -Пелева,0,41,Южен Парк,2013-05-19,46m32s,70,6.45,9m18s
```

where every line is a one particular run for a participant.

## Install

You can grab a binary from the [release page](/releases).

### From Source

As usual, just

```
go get github.com/ironsmile/fivekmscraper
```
## Usage

Just run the binary. From time time you would be asked whether a name is a male or a female. Answer with "f" or "m".

