import string
import numpy as np
import tensorflow as tf
import keras
import os
from domain.dataset.dataset import Dataset
from domain.message.message import Message
from typing import Tuple

class ToneEstimationModel:
    model: keras.models.Sequential
    model_path: string

    def __init__(self, model_path: string = "", train: Dataset = None, valid: Dataset = None, test: Dataset = None):
        self.model_path = model_path
        if (model_path != ""):
            self.model = keras.models.load_model(model_path)
            return
        if (train != None and valid != None and test != None):
            self.create_model(train, valid)
            self.estimate_model(test)
            self.save()
        

    def save(self):
        self.model.save(self.model_path)
    

    def create_model(self, train_dataset: Dataset, valid_dataset: Dataset):
        tf.random.set_seed(42)

        vectorized = tf.keras.layers.TextVectorization(max_tokens=15000)
        vectorized.adapt(train_dataset.map(lambda target, labels: target))

        model = keras.Sequential([
            vectorized,
            keras.layers.Embedding(15000, 128, mask_zero=True),
            keras.layers.GRU(128),
            keras.layers.Dense(2, activation="softmax")
        ])

        model.compile(loss="sparse_categorical_crossentropy", optimizer="adam", metrics=["accuracy"])
        model.fit(train_dataset, validation_data=valid_dataset, epochs=2)

        self.model = model
    
    def estimate_model(self, test: Dataset):
        self.model.summary()

        y = self.model.predict(test)
        print(np.argmax(y))
        loss, accuracy = self.model.evaluate(test)
        print("Accuracy:", accuracy)
    
    def create_dataset(self, text: str, sentiment: int) -> Dataset:
        text = tf.constant(text, dtype=tf.string)
        label = tf.constant(sentiment, dtype=tf.int32)

        raw_ds = tf.data.Dataset.from_tensors((text, label))
        ds = raw_ds.batch(1).prefetch(1)
        return ds

    def estimate_tone(self, messages: list[Message]) -> tuple[float, float]:
        positive_count = 0
        negative_count = 0

        for i, message in enumerate(messages):
            set = self.create_dataset(message.text, 0)
            result = self.model.predict(set)

            negative_value = result[0][0]
            positive_value = result[0][1]
            if negative_value > positive_value:
                negative_count = negative_count+1   
            if positive_value > negative_value:
                positive_count = positive_count+1

        return positive_count/len(messages), negative_count/len(messages)

