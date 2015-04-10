from yaafelib import FeaturePlan, Engine, AudioFileProcessor
import mutagen
import numpy
import os
import sys, subprocess, sqlite3
import re
import random

def interpret(filename):
    name = filename[filename.rfind("-"):]
    name = name.replace("- ", "")
    name = name.replace(".wav", "")
    name = name.replace("middle", "")
    name = name.replace("higher", "")
    name = name.replace("High Pitch", "")
    name = name.replace("Low Pitch", "")
    return name

def analyze_file(filename):
    name = interpret(filename)
    fp = FeaturePlan(sample_rate=44100)
    engine = Engine()
    file_processor = AudioFileProcessor()

    fp.addFeature('mfcc: MFCC blockSize=2048 stepSize=1024')
    data_flow = fp.getDataFlow()
    engine.load(data_flow)
    file_processor.processFile(engine, filename)
    data = engine.readAllOutputs()
    mfcc = data["mfcc"]
    mean = numpy.mean(mfcc, 0)
    result = []
    for vector in mfcc:
        if sum(numpy.greater(vector,mean)) > 8:
            result.append(vector)
    chord = numpy.mean(numpy.array(result), 0)
    numpy.save("chords/%s" % name, chord)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Not enough arguments please remember to include the file")
        exit()
    else:
        filedir = sys.argv[1]
        for filename in os.listdir(filedir):
            analyze_file(os.path.join(filedir, filename))

