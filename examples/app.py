import logging
import argparse

from pydantic import BaseModel, confloat, constr
from ventu import Ventu
import torch
import numpy as np
from transformers import DistilBertTokenizer, DistilBertForSequenceClassification


class Req(BaseModel):
    # the input sentence 