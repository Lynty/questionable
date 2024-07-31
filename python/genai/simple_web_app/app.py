import os

from flask import Flask, request

import google.auth
import vertexai
from vertexai.generative_models import (
    GenerativeModel,
    Content,
    FunctionDeclaration,
    Part,
    Tool,
)

import requests

_, project = google.auth.default()
MODEL_ID = "gemini-1.5-flash"
#MODEL_ID = "gemini-1.5-pro-001"

app = Flask(__name__)

@app.route("/")
def animal_fun_facts():
    vertexai.init(project=project, location="us-central1")
    model = GenerativeModel(
        MODEL_ID,
        system_instruction=[
            "You are a fun kindergarten teacher.",
            "Your mission is to provide information to young children and provide emojis when relevant.",
        ],
    )
    animal = request.args.get("animal", "aardvark") 
    prompt = f"Give me 10 fun facts about {animal}. Return this as html without backticks."
    #prompt = """
    #  User input: I like mangos.
    #  Answer:
    #"""
    response = model.generate_content(prompt)
    return response.text

if __name__ == "__main__":
    app.run(debug=True, host="0.0.0.0", port=int(os.environ.get("PORT", 8080)))
