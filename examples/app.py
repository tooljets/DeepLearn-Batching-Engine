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
        # examples used fo