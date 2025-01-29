#/usr/bin/env python3
import argparse
import asyncio
import os
import sys
import yaml
import socket
import threading
import socketserver


class ThreadedTCPRequestHandler(socketserver.BaseRequestHandler):

    def handle(self):
        data = str(self.request.recv(1024), 'ascii')
        cur_thread = threading.current_thread()
        response = bytes("{}: {}".format(cur_thread.name, data), 'ascii')
        self.request.sendall(response)


class ThreadedTCPServer(socketserver.ThreadingMixIn, socketserver.TCPServer):
    pass


def client(ip, port, message):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        sock.connect((ip, port))
        sock.sendall(bytes(message, 'ascii'))
        response = str(sock.recv(1024), 'ascii')
        print("Received: {}".format(response))


class Listener:
    def __init__(self, args):
        self.args = args


class Connection:
    def __init__(self, args):
        self.args = args


class Muxator:
    def __init__(self, args):
        self.args = args

        async def main(self, args):
            if os.geteuid() == 0:
                print('You cannot run this script as root')
                sys.exit(1)



def run_deployment_tool():
    parser = argparse.ArgumentParser(description='V6 Deployment Tool V2')
    parser.add_argument('-i', '--listen_address',
                        help='IP to serve the proxy on',
                        required=True)
    parser.add_argument('-p', '--listen_port',
                        help='Port to serve the proxy on',
                        required=True)
    parser.add_argument('-n', '--number',
                        help='Number of Tor connections to open',
                        required=True)
    args = parser.parse_args()
    app = Muxator(args)
    asyncio.run(app.main())


if __name__ == "__main__":
    run_deployment_tool()
