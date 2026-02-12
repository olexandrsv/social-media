import os
from wrapt import synchronized
from domain.tone_model.model import ToneEstimationModel
from repository.repository import Repository
from domain.message.message import Message
from typing import Tuple
import tensorflow as tf

class Service:
    model: ToneEstimationModel
    repository: Repository

    def __init__(self, repository: Repository):
        print("Get data")
        model_path = "model1.keras"
        self.repository = repository
        exists = os.path.exists(model_path)
        if exists:
            self.model = ToneEstimationModel(model_path)
            return
        train, valid, test = repository.get_tone_data()
        print("Received data")
        self.model = ToneEstimationModel(model_path, train, valid, test)

    @synchronized   
    def estimate_tone(self, messages: list[Message]) -> tuple[float, float]:
        if (len(messages) == 0):
            return 0, 0
        return self.model.estimate_tone(messages)


       