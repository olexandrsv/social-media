from flask import Flask, request, jsonify
from service.service import Service
from domain.message.message import Message
from server.models import PredictToneResp, ToneEstimation

class Server:
    app: Flask
    service: Service

    def __init__(self, service):
        self.app = Flask(__name__)
        self.service = service


        @self.app.route("/data", methods=["POST"])
        def receive_json():
            print("Request received")

            data = request.get_json()
            if data is None:
                return jsonify({"error": "Invalid JSON"}), 400

            messages: list[Message] = list()
            
            for i, messageModel in enumerate(data["messages"]):
                messages.append(Message(messageModel["text"]))
            print("messages", messages)

            positive_percentage, negative_percentage = self.service.estimate_tone(messages)

            print("Percentages: ", positive_percentage, negative_percentage)
            resp = PredictToneResp(
                tone_estimation = ToneEstimation(
                    positive_percentage=positive_percentage,
                    negative_percentage=negative_percentage
                )
            )

            return jsonify(resp)

    def run(self):
        self.app.run(host="0.0.0.0", port=8000, debug=True)
