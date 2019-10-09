from pyftpdlib import servers
from pyftpdlib.handlers import FTPHandler
from pyftpdlib.authorizers import DummyAuthorizer
import os

authorizer = DummyAuthorizer()
authorizer.add_anonymous(os.getcwd(), perm="elradfmwMT")

address = ("127.0.0.1", 21)
handler = FTPHandler
handler.authorizer = authorizer
server = servers.FTPServer(address, handler)
server.serve_forever()
