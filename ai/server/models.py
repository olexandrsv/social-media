from dataclasses import dataclass

@dataclass
class ToneEstimation:
    positive_percentage: float
    negative_percentage: float

@dataclass
class PredictToneResp:
    tone_estimation: ToneEstimation