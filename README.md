# Web Replay Tool (Replay .har)

<p align="center"><img src="./.docs/demo.gif" width="500" ></p>

## Introduction

This tool lets you replay an exported .har file from the DevTools network tab and replay it.

The tool spins up a local server and replays the .har file on that local server.

I wrote this tool for performance testing the client-side portion of a web app, removing the latency introduced by the server and allowing me to test the client-side performance - but it can be used for website archiving, modifying content for demos and so on. Thought it might be useful for someone.

## Usage

First, record your network activity from your live site and Save the network tab to a `.har` file.

```
web-reply -ssl-key /localhost.key -ssl-cert /localhost.cert ./your.har
```

You will need to generate a self-signed SSL certificate. I recommend using [mkcert](https://github.com/FiloSottile/mkcert) as it's a one line command that does everything for you.

Open Chrome on `https://localhost:3000/path/you/saved`

This will replay the `.har` file exactly, but it will not allow you to use the entire app.

## Patching Sites

You can patch responses, add latency, inject scripts and modify the responses.

Add patches to the `patches/enabled` folder next to the binary to run them.

You can see examples here: [patches/disabled](https://github.com/alshdavid/web-replay/tree/main/commands/web-replay/patches/disabled)https://github.com/alshdavid/web-replay/tree/main/commands/web-replay/patches/disabled

