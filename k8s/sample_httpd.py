#!/usr/bin/env python3
# -*- coding: utf-8 -*-
#
# The original code is written by @tkj
#    https://qiita.com/tkj/items/210a66213667bc038110
#
import os
import http.server as s
from urllib.parse import urlparse
from urllib.parse import parse_qs

class MyHandler(s.BaseHTTPRequestHandler):
    def do_POST(self):
        self.make_data()
    def do_GET(self):
        self.make_data()
    def make_data(self):
        # get URL parameters
        parsed = urlparse(self.path)
        # parse URL parameters
        params = parse_qs(parsed.query)
        # get BODY
        content_len  = int(self.headers.get("content-length"))
        req_body = self.rfile.read(content_len).decode("utf-8")

        # memory alloc
        alloc_size = int(req_body)
        mem_alloc = [ x for x in range(alloc_size) ]

        # Create response body
        body  = "method: " + str(self.command) + "\n"
        body += "params: " + str(params) + "\n"
        body += "body  : " + req_body + "\n"
        body += "data  : " + str(mem_alloc) + "\n"
        self.send_response(200)
        self.send_header('Content-type', 'text/html; charset=utf-8')
        self.send_header('Content-length', len(body.encode()))
        self.end_headers()
        self.wfile.write(body.encode())
host = '0.0.0.0'
port = 8000
httpd = s.HTTPServer((host, port), MyHandler)
print('HTTP server is spawned. Port:%s' % port)
httpd.serve_forever()
