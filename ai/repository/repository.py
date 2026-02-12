import tensorflow as tf
import keras.models
import tensorflow_datasets as tfds

from domain.dataset.dataset import Dataset

class Repository:
    def __init__(self):
        pass

    def get_tone_data(self) -> tuple[Dataset, Dataset, Dataset]:
        raw_train_dataset, raw_valid_dataset, raw_test_dataset = tfds.load(
            name="sentiment140",
            split=["train[:90%]", "train[90%:]", "test"],
            as_supervised=True
        )

        raw_train_dataset = raw_train_dataset.map(self.fix_labels, num_parallel_calls=tf.data.AUTOTUNE)
        raw_valid_dataset = raw_valid_dataset.map(self.fix_labels, num_parallel_calls=tf.data.AUTOTUNE)
        raw_test_dataset = raw_test_dataset.filter(lambda text, label: tf.not_equal(label, 2))
        raw_test_dataset = raw_test_dataset.map(self.fix_labels, num_parallel_calls=tf.data.AUTOTUNE)

        train_dataset = raw_train_dataset.shuffle(20000, seed=42).batch(32).prefetch(1)
        valid_dataset = raw_valid_dataset.batch(32).prefetch(1)
        test_dataset = raw_test_dataset.batch(32).prefetch(1)

        return train_dataset, valid_dataset, test_dataset
    
    def fix_labels(self, text, label):
        new_label = tf.where(label == 4, 1, 0)
        return text, new_label

