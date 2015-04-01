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

    fp.addFeature('mfcc: MFCC blockSize=256 stepSize=128')
    fp.addFeature('ComplexDomainOnsetDetection: ComplexDomainOnsetDetection blockSize=512 stepSize=256')
    data_flow = fp.getDataFlow()
    engine.load(data_flow)
    file_processor.processFile(engine, filename)
    return engine.readAllOutputs()

def setup_cov(mfcc):
    sigma = numpy.array([[0]*13]*13, dtype="object")
    usigma = numpy.array([[0]*13]*13, dtype="object")
    for row in mfcc:
        # Calculate \Sigma X_i
        usigma += [row]*13
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
    files = [ f for f in os.listdir(path) if os.path.isfile(os.path.join(path,f)) ]
    conn = sqlite3.connect("music.db")
    c = conn.cursor()
    c.execute("""
    CREATE TABLE IF NOT EXISTS songs (title varchar(255), artist varchar(255),
    album varchar(255), mfcc varchar(255), size float, beat float);
    """)
    for filename in files:
        if ".mp3" in filename:
            metadata = convert(os.path.join(path, filename))
            data = analyze_file(metadata["filename"])
            data["mfcc"].flags.writeable = False
            stdev = numpy.std(data["ComplexDomainOnsetDetection"])
            total = 0
            last = 0
            count = 0
            for val in range(len(data["ComplexDomainOnsetDetection"])):
                total += data["ComplexDomainOnsetDetection"][val]
                last = val
                count += 1
            if count:
                beat = float(total)/float(count)
                mfcc = data["mfcc"]
                # precalculate what we need to combine cov
                sigma, usigma = setup_cov(mfcc)
                mean = numpy.mean(mfcc, 0) # calculate mean
                n = mfcc.shape[0]
                nmean = numpy.transpose([mean]*13)
                cov =  numpy.cov(mfcc.T) # calculate covariance
                mfccname = str(hash(data["mfcc"].data)).replace("-", "")
                numpy.save("mfcc/%s-mean" % mfccname, mean)
                numpy.save("mfcc/%s-cov" % mfccname, cov)
                numpy.save("mfcc/%s-sigma" % mfccname, sigma)
                numpy.save("mfcc/%s-usigma" % mfccname, usigma)
                c.execute("""
                INSERT INTO songs VALUES (?, ?, ?, ?, ?, ?);
                """, (metadata["title"], metadata["artist"], metadata["album"],
                    mfccname, mfcc.shape[0], beat))
    conn.commit()
    conn.close()
