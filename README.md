# mqtt-automation-hat-go
MQTT driven [Pimoroni Automation HAT](https://shop.pimoroni.com/products/automation-hat) driver (Go implementation).

**FAIR WARNING:** This project is a Go learning exercise, expect alpha quality.

The best way to make this program useful at this point is to use it together with [Home Climate Control](https://github.com/home-climate-control/dz) configured with `MqttConnector` enabled. Search [DIY Zoning & Home Climate Control Forum](https://groups.google.com/forum/#!forum/home-climate-control) for "MQTT", you will find configuration instructions.

It is also possible to use any MQTT publisher to control this program, but you'll have to either inspect the code to figure out the JSON payload format, or follow instructions above.

Criticism is welcome - both as [bug reports](https://github.com/climategadgets/mqtt-automation-hat-go/issues) and [suggestions](https://groups.google.com/forum/#!forum/home-climate-control).
