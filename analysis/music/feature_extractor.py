from yaafelib import FeaturePlan, Engine, AudioFileProcessor
import mutagen
import numpy
import os
import sys, subprocess, sqlite3
import re
import random

def analyze_file(filename):
    fp = FeaturePlan(sample_rate=44100)
    engine = Engine()
    file_processor = AudioFileProcessor()

    fp.addFeature('mfcc: MFCC blockSize=512 stepSize=256')
    #fp.addFeature('ComplexDomainOnsetDetection: ComplexDomainOnsetDetection blockSize=512 stepSize=256')
    fp.addFeature('loudness: Loudness blockSize=512 stepSize=256')
    data_flow = fp.getDataFlow()
    engine.load(data_flow)
    file_processor.processFile(engine, filename)
    return engine.readAllOutputs()

def setup_cov(matrix):
    width = matrix.shape[1]
    sigma = numpy.array([[0]*width]*width, dtype="object")
    usigma = numpy.array([[0]*width]*width, dtype="object")
    for row in matrix:
        # Calculate \Sigma X_i
        sigma += [row]*width
        # Calculate the \Sigma X_iX_i
        # use [row] because row.T on a 1D makes no sense
        sigma += numpy.transpose([row]) * row
    return sigma, usigma

def convert(filename):
    result = subprocess.check_output(["mpg123",
            "-w",
            filename.replace(".mp3", ".wav"),
            filename],stderr=subprocess.STDOUT)
    metadata = {"filename" : filename.replace(".mp3", ".wav")}
    try:
        regex = re.compile("Title:(.*?)Artist:(.*?)\n.*?\nAlbum:(.*?)\n")
        items = regex.search(result)
        title, artist, album = items.groups()
        metadata["title"] = title.strip()
        metadata["artist"] = artist.strip()
        metadata["album"] = album.strip()
    except Exception as e:
        metadata["title"] = os.path.basename(metadata["filename"]).replace(".wav", "")
        metadata["artist"] = "unknown"
        metadata["album"] = "unknown"
    return metadata

if __name__ == "__main__":
    path = sys.argv[1]
    files = os.walk(path)
    conn = sqlite3.connect("music.db")
    c = conn.cursor()
    c.execute("""
    CREATE TABLE IF NOT EXISTS songs (title varchar(255), artist varchar(255),
    album varchar(255), mfcc varchar(255), msize float, loudness varchar(255),
    lsize float);
    """)
    count = 0
    for dirpath, dirnames, filenames in files:
        for filename in filenames:
            try:
                if ".mp3" in filename:
                    metadata = convert(os.path.join(dirpath, filename))
                    previous = c.execute("""SELECT * FROM songs WHERE title=?""",
                            ([metadata["title"]])).fetchall()
                    if len(previous) == 0:
                        data = analyze_file(metadata["filename"])
                        data["mfcc"].flags.writeable = False
                        data["loudness"].flags.writeable = False
                        total = 0
                        last = 0
                        # precalculate what we need to combine mfcc cov
                        mfcc = data["mfcc"]
                        sigma, usigma = setup_cov(mfcc)
                        mean = numpy.mean(mfcc, 0) # calculate mean
                        nmean = numpy.transpose([mean]*13)
                        cov =  numpy.cov(mfcc.T) # calculate covariance
                        mfccname = str(hash(data["mfcc"].data)).replace("-", "")
                        numpy.save("mfcc/%s-mean" % mfccname, mean)
                        numpy.save("mfcc/%s-cov" % mfccname, cov)
                        numpy.save("mfcc/%s-sigma" % mfccname, sigma)
                        numpy.save("mfcc/%s-usigma" % mfccname, usigma)
                        # precalculate what we need to combine loudness cov
                        loudness= data["loudness"]
                        sigma, usigma = setup_cov(loudness)
                        mean = numpy.mean(loudness, 0) # calculate mean
                        nmean = numpy.transpose([loudness]*loudness.shape[1])
                        cov =  numpy.cov(loudness.T) # calculate covariance
                        name = str(hash(loudness.data)).replace("-", "")
                        numpy.save("loudness/%s-mean" % name, mean)
                        numpy.save("loudness/%s-cov" % name, cov)
                        numpy.save("loudness/%s-sigma" % name, sigma)
                        numpy.save("loudness/%s-usigma" % name, usigma)
                        c.execute("""
                        INSERT INTO songs VALUES (?, ?, ?, ?, ?, ?, ?);
                        """, (metadata["title"], metadata["artist"], metadata["album"],
                            mfccname, mfcc.shape[0], name, loudness.shape[0]))
                        count += 1
                        print count/500.0
                    else:
                        print "You've alraedy analyzed %s" % metadata["title"]
                    print "deleting %s" % metadata["filename"]
                    os.remove(os.path.join(dirpath, metadata["filename"]))
                if count % 50 == 0:
                    conn.commit()
                    conn.close()
                    conn = sqlite3.connect("music.db")
                    c = conn.cursor()
            except Exception as e:
                print e
