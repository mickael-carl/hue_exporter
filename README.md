# Hue exporter

Prometheus exporter for the Philips Hue IoT system. 

## Features

### Lights

Metrics are available for the following:

 * Reachable
 * On/Off
 * Color
 * Brightness

It supports both colored and white bulbs from Hue or compatible bulbs.

### Sensors

Metrics are available for the following:

 * Reachable
 * On/Off

Additionally the following is also available:
 
 * Temperature sensor:
   * Battery
   * Temperature

## Running the exporter

Currently only supports querying the bridge over HTTP.
The token can be retrieved when registering an application on the bridge.

`./hue_exporter --bridge http://ADDRESS_BRIDGE --token TOKEN`

## Building

Building can be achieved using Bazel as follows:
`bazel build //...`
