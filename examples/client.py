
from concurrent import futures
import httpx
import msgpack


URL = 'http://localhost:8080/inference'
HEADER = {'Content-Type': 'application/msgpack'}
packer = msgpack.Packer(
    autoreset=True,