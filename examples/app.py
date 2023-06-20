import logging
import argparse

from pydantic import BaseModel, confloat, constr
from ventu import Ventu
import torch
import numpy as np
from transformers import DistilBertTokenizer, DistilBertForSequenceClassification


class Req(BaseModel):
    # the input sentence should be at least 2 characters
    text: constr(min_length=2)

    class Config:
        # examples used for health check and warm-up
        schema_extra = {
            'example': {'text': 'my cat is very cut'},
            'batch_size': 16,
        }


class Resp(BaseModel):
    positive: confloat(ge=0, le=1)
    negative: confloat(ge=0, le=1)


class ModelInference(Ventu):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.tokenizer = DistilBertTokenizer.from_pretrained(
            'distilbert-base-uncased')
        self.model = DistilBertForSequenceClassification.from_pretrained(
            'distilbert-base-uncased-finetuned-sst-2-english')

    def preprocess(self, data: Req):
        tokens = self.tokenizer.encode(data.text, add_special_tokens=True)
        return tokens

    def batch_inference(self, data):
        # batch inference is used in `socket` mode
        data = [torch.tensor(token) for token in data]
        with torch.no_grad():
            result = self.model(torch.nn.utils.rnn.pad_sequence(data, batch_first=True))[0]
        return result.numpy()

    def inference(self, data):
        # inference is used in `http` mode
        with torch.no_grad():
            result = self.model(torch.tensor(data).unsqueeze(0))[0]
        return result.numpy()[0]

    def postprocess(self, data):
        scores = (np.exp(data) / np.exp(data).sum(-1, keepdims=True