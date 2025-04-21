# WIP: Glimpse - Simple Realtime Server Health view

Glimpse is a simple real time dashboard for homelabs. It shows near real-time metrics such as CPU, Memory etc in a unified dashboard

## But why? Ever heard of prometheus + node exporter?

Because why not! Honestly, it is mostly for the fun of building something.
Also I wanted something simple for my homelab to show real time view of all the servers without much configuration. I could use node_exporter + Grafana to do that, sure, but that is not so fun!

## Project Status

Don't use it. I am working on it. Soon™

## What's the vibe like?

Honestly, I did try to make something via vibe-coding. I am too much of a control freak that it did not go well. Now, I decided the best method that works for me is to use ChatGPT for brainstorming and generate individual functions and then type it up into the IDE.

## Overview

Agent runs on each host, connects to the server over GRPC to send metrics. There is no storage (only in memory short term storage). I am not re-writing prometheus.