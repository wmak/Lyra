import matplotlib.pyplot as plt
from matplotlib.mlab import PCA
import numpy as np
import scipy as sp
from sklearn import decomposition
portal_data = np.genfromtxt("portal.wav.mfcc.csv", skip_header = 5, delimiter=",")
running_data = np.genfromtxt("running.wav.mfcc.csv", skip_header = 5, delimiter=",")
boxtop_data = np.genfromtxt("boxtops.wav.mfcc.csv", skip_header = 5, delimiter=",")
print portal_data.shape
def reduce(data):
    data = np.swapaxes(data, 0, 1)
    pca = decomposition.PCA(n_components=1)
    pca.fit(data)
    X = pca.transform(data)
    return X
p = reduce(portal_data)
r = reduce(running_data)
b = reduce(boxtop_data)
print p - r
print p - b
plt.plot(p)
plt.plot(r)
plt.plot(b)
plt.show()
