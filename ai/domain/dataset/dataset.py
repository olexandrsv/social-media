from typing import Tuple, TypeAlias
import tensorflow as tf

Dataset: TypeAlias = tf.data.Dataset[Tuple[tf.Tensor, tf.Tensor]]
