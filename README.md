# bellis(in progress)

[Demo Site](https://bellis.minoic.top)

The distributed cloud monitoring and alarm system adopts a design with clear division of labor between modules. The modules communicate through GRPC and message queue (Redis), and the sampled data is stored through InfluxDB.

![screenshot](https://github.com/bellis-daemon/bellis/assets/55825043/9d7f09bf-5a39-414f-9390-c90a97c2b72c)

## Configure

### ETCD

| key                     | description                                 |
| ----------------------- | ------------------------------------------- |
| `influxdb_token`        | api token for influxdb.                     |
| `telegram_bot_token`    | telegram bot token generated via @BotFather |
| `tencent_smtp_password` | tencent cloud ses service smtp password     |
