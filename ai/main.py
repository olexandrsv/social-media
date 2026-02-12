import numpy as np
from numpy.random import seed
from repository.repository import Repository
from service.service import Service
from server.server import Server

seed(1)
import matplotlib.pyplot as plt
import pandas as pd
import string
import os
import shutil
import re

import tensorflow as tf

import tensorflow_datasets as tfds
import keras
import keras.models
from keras.models import Sequential
from keras.layers import SimpleRNN, Dense, Embedding
from tensorflow.keras.layers import TextVectorization
from keras.optimizers import Adam
from typing import Tuple

def main():
    repo = Repository()
    service = Service(repo)
    server = Server(service)
    server.run()

main()