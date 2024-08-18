import json
import random

import torch

from train import NeuralNet, tokenize, bag_of_words, train

from fastapi import FastAPI, Request

device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')

# Load intents
with open('/data/intents.json', 'r') as json_data:
    intents = json.load(json_data)

# Function to load the model
def load_model(file_path: str):
    data = torch.load(file_path, map_location=device, weights_only=True)
    input_size = data["input_size"]
    hidden_size = data["hidden_size"]
    output_size = data["output_size"]
    all_words = data['all_words']
    tags = data['tags']
    model_state = data["model_state"]

    model = NeuralNet(input_size, hidden_size, output_size).to(device)
    model.load_state_dict(model_state)
    model.eval()

    return model, all_words, tags


# Initial load of the model
model, all_words, tags = load_model("/data/data.pth")

app = FastAPI()

@app.post("/")
async def send_message(request: Request):
    req_json = await request.json()

    print(req_json["query"])

    sentence = tokenize(req_json["query"])
    x = bag_of_words(sentence, all_words)
    x = x.reshape(1, x.shape[0])
    x = torch.from_numpy(x).to(device)

    output = model(x)
    _, predicted = torch.max(output, dim=1)

    tag = tags[predicted.item()]

    probs = torch.softmax(output, dim=1)
    prob = probs[0][predicted.item()]
    if prob.item() > 0.75:
        for intent in intents['intents']:
            if tag == intent['tag']:
                return random.choice(intent['responses'])
    else:
        return "I do not understand..."


@app.get("/retrain")
def retrain():
    train()

    # Reload the model after training
    global model, all_words, tags, intents
    model, all_words, tags = load_model("/data/data.pth")

    with open('/data/intents.json', 'r') as json_data:
        intents = json.load(json_data)

    return "Model retrained and reloaded"

@app.post("/new/intent")
async def new_intent(request: Request):
    req_json = await request.json()

    with open('/data/intents.json', 'r') as json_data:
        intents = json.load(json_data)

    intents["intents"].append(req_json)

    f = open('/data/intents.json', 'w')
    f.write(json.dumps(intents))
    f.close()

@app.post("/new/response")
async def new_response(request: Request):
    req_json = await request.json()

    with open('/data/intents.json', 'r') as json_data:
        intents = json.load(json_data)

    for intent in intents['intents']:
        if req_json["tag"] == intent["tag"]:
            intent["responses"].append(req_json["response"])

    f = open("/data/intents.json", "w")
    f.write(json.dumps(intents))
    f.close()

@app.get("/get/intents")
async def get_intents():
    with open('/data/intents.json', 'r') as json_data:
        intents = json.load(json_data)
        return intents
