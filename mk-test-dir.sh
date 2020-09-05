#!/bin/bash

TEST_FOLDER="tmp-test"
TV="tv"
MOVIE="movies"

# $1 (first command line arg) specifies path to write all files to
mkdir -p "$1/$TEST_FOLDER"
if [ $? -ne 0 ] ; then
    echo Failure
    exit 1
fi

DIRS=(
  "$1/$TEST_FOLDER/$MOVIE"
  "$1/$TEST_FOLDER/$TV"
  "$1/$TEST_FOLDER/$TV/Broadchurch.Season.2.Complete.720p.HDTV.x264-SCENE"
  "$1/$TEST_FOLDER/$TV/The Wire Season 1"
  "$1/$TEST_FOLDER/$MOVIE/Captain Marvel (2019)"
)

FILES=(
  "$1/$TEST_FOLDER/$TV/Broadchurch.Season.2.Complete.720p.HDTV.x264-SCENE/Broadchurch.S02E01.720p.HDTV.x264-FTP.mkv"
  "$1/$TEST_FOLDER/$TV/Broadchurch.Season.2.Complete.720p.HDTV.x264-SCENE/Broadchurch.S02E02.720p.HDTV.x264-TLA.mkv"
  "$1/$TEST_FOLDER/$TV/Broadchurch.Season.2.Complete.720p.HDTV.x264-SCENE/Broadchurch.S02E03.720p.HDTV.x264-TLA.mkv"
  "$1/$TEST_FOLDER/$TV/Broadchurch.Season.2.Complete.720p.HDTV.x264-SCENE/Broadchurch.S02E04.REPACK.720p.HDTV.x264-TLA.mkv"
  "$1/$TEST_FOLDER/$TV/Broadchurch.Season.2.Complete.720p.HDTV.x264-SCENE/Broadchurch.S02E05.720p.HDTV.x264-TLA.mkv"
  "$1/$TEST_FOLDER/$TV/Broadchurch.Season.2.Complete.720p.HDTV.x264-SCENE/Broadchurch.S02E06.PROPER.720p.HDTV.x264-SRiZ.mkv"
  "$1/$TEST_FOLDER/$TV/Broadchurch.Season.2.Complete.720p.HDTV.x264-SCENE/Broadchurch.S02E07.720p.HDTV.x264-ORGANiC.mkv"
  "$1/$TEST_FOLDER/$TV/Broadchurch.Season.2.Complete.720p.HDTV.x264-SCENE/Broadchurch.S02E08.720p.HDTV.x264-ORGANiC.mkv"
  "$1/$TEST_FOLDER/$TV/The Wire Season 1/The Wire [1x03] The Buys.mkv"
  "$1/$TEST_FOLDER/$TV/The Wire Season 1/The Wire [1x04] The Old Cases.mkv"
  "$1/$TEST_FOLDER/$TV/The.Wire.S04E02.1080p.5.1Ch.BluRay.ReEnc-DeeJayAhmed.mkv"
  "$1/$TEST_FOLDER/$TV/The.Wire.S04E03.1080p.5.1Ch.BluRay.ReEnc-DeeJayAhmed.mkv"
  "$1/$TEST_FOLDER/$MOVIE/Captain Marvel (2019)/Captain.Marvel.2019.2160p.4K.BluRay.x265.10bit.AAC7.1.1-[YTS.MX].mkv"
  "$1/$TEST_FOLDER/$MOVIE/The.Lion.King.2019.1080p.BluRay.x264-[YTS.LT].mp4"
)

#Make directories
for i in "${DIRS[@]}"
do
   mkdir -p "$i"
   #echo "creating folder: $i"
done

#Make files
for i in "${FILES[@]}"
do
   touch "$i"
   #echo "creating file: $i"
done
