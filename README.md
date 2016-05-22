## Summary

A small library that converts CCP xml API responses to json.

## Usage

### As a library

Use either `xmljsonproxy.Transform` (if you have a `io.Reader`) or `xmljsonproxy.TransformString` (if you have a `string`) to transform xml to json.

### As a proxy

Run the binary, and change your endpoint from `https://api.eveonline.com` to `http://localhost:9293`. You can overwrite the port with the `$PORT` environment variable.