#!/usr/bin/env python 

import numpy as np
import cv2
import cv2.cv as cv
import sys

def analyze(filename):
    result = {}
    image = cv2.imread(filename)
    gray_image = cv2.cvtColor(image, cv.CV_RGB2GRAY)
    gray_image = cv2.equalizeHist(gray_image)
    face_cascade = cv2.CascadeClassifier('haarcascade_frontalface_default.xml')
    faces = face_cascade.detectMultiScale(gray_image, 1.1, 2, cv.CV_HAAR_SCALE_IMAGE,
            (20,20))
    result.setdefault("faces", len(faces))
    return result

if __name__ == "__main__":
    print("Starting image analysis")
    print(analyze(sys.argv[1]))
    print("Image analysis complete")

