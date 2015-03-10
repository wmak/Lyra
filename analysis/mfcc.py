import numpy as np
from sklearn.neighbors import NearestNeighbors
import sys

# Read in MFCC csv files generated by Yaafe
files = sys.argv[1:]

points = []
for file in files:
    # read the csv file
    raw_data = np.genfromtxt(file, skip_header = 5, delimiter = ",")
    # Take the average of each column so we have a 1x13 vector
    vector = np.mean(raw_data, 0)
    # Append this to our list of all vectors
    points.append(vector)
# Calculate the nearest neighbors for the entire list
nbrs = NearestNeighbors(n_neighbors=2, algorithm='ball_tree').fit(points)
dist, indices = nbrs.kneighbors(points)

# Print out what we've discovered
for indice in indices:
    print "Song A: %s\nSong B: %s\n" % (
        files[indice[0]][4:-17],
        files[indice[1]][4:-17])
