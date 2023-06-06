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

    def bat